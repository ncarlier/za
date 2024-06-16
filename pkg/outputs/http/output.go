package http

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/ncarlier/za/pkg/conditional"
	"github.com/ncarlier/za/pkg/config"
	"github.com/ncarlier/za/pkg/events"
	"github.com/ncarlier/za/pkg/helper"
	"github.com/ncarlier/za/pkg/outputs"
	"github.com/ncarlier/za/pkg/serializers"
)

const (
	defaultClientTimeout = 5 * time.Second
	defaultMethod        = http.MethodPost
	defaultUserAgent     = "Mozilla/5.0 (compatible; ZeroAnalytics/1.0; +https://github.com/ncarlier/za)"
)

// Output for HTTP
type Output struct {
	URL      string            `toml:"url"`
	Timeout  config.Duration   `toml:"timeout"`
	Method   string            `toml:"POST"`
	Username string            `toml:"username"`
	Password string            `toml:"password"`
	Headers  map[string]string `toml:"headers"`
	Gzip     bool              `toml:"gzip"`

	client     *http.Client
	serializer serializers.Serializer
	condition  conditional.Expression
}

func (h *Output) createClient(ctx context.Context) (*http.Client, error) {
	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
		},
		Timeout: h.Timeout.Duration,
	}

	return client, nil
}

// SetSerializer set data serializer
func (o *Output) SetSerializer(serializer serializers.Serializer) {
	o.serializer = serializer
}

// SetCondition set condition expression
func (o *Output) SetCondition(condition conditional.Expression) {
	o.condition = condition
}

// Connect activate the output writer
func (o *Output) Connect() error {
	if o.URL == "" {
		return fmt.Errorf("invalid URL: %s", o.URL)
	}
	if o.Method == "" {
		o.Method = defaultMethod
	}
	o.Method = strings.ToUpper(o.Method)
	if o.Method != http.MethodPost && o.Method != http.MethodPut {
		return fmt.Errorf("invalid method [%s] %s", o.URL, o.Method)
	}

	if o.Timeout.Duration == 0 {
		o.Timeout.Duration = defaultClientTimeout
	}

	ctx := context.Background()
	client, err := o.createClient(ctx)
	if err != nil {
		return err
	}

	o.client = client

	slog.Debug("using HTTP output", "uri", o.URL)

	return nil
}

// Close the output writer
func (o *Output) Close() error {
	return nil
}

// Description returns description
func (o *Output) Description() string {
	return "Send page view to HTTP endpoint"
}

// SendEvent send event to the Output
func (o *Output) SendEvent(event events.Event) error {
	if !o.condition.Match(event) {
		return nil
	}
	b, err := o.serializer.Serialize(event)
	if err != nil {
		return fmt.Errorf("unable to serialize page view: %v", err)
	}

	if err := o.send(b); err != nil {
		return err
	}

	return nil
}

func (o *Output) send(reqBody []byte) error {
	var reqBodyBuffer io.Reader = bytes.NewBuffer(reqBody)

	var err error
	if o.Gzip {
		rc, err := helper.CompressWithGzip(reqBodyBuffer)
		if err != nil {
			return err
		}
		defer rc.Close()
		reqBodyBuffer = rc
	}

	req, err := http.NewRequest(o.Method, o.URL, reqBodyBuffer)
	if err != nil {
		return err
	}

	if o.Username != "" || o.Password != "" {
		req.SetBasicAuth(o.Username, o.Password)
	}

	req.Header.Set("User-Agent", defaultUserAgent)
	req.Header.Set("Content-Type", o.serializer.ContentType())
	if o.Gzip {
		req.Header.Set("Content-Encoding", "gzip")
	}
	for k, v := range o.Headers {
		if strings.EqualFold(k, "host") {
			req.Host = v
		}
		req.Header.Set(k, v)
	}

	resp, err := o.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	_, err = io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("bad status code (%d) when writing to [%s]", resp.StatusCode, o.URL)
	}

	return nil
}

func init() {
	outputs.Add("http", func() outputs.Output {
		return &Output{
			Timeout: config.Duration{Duration: defaultClientTimeout},
			Method:  defaultMethod,
			Gzip:    true,
		}
	})
}
