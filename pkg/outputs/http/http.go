package http

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/ncarlier/za/pkg/config"
	"github.com/ncarlier/za/pkg/events"
	"github.com/ncarlier/za/pkg/helper"
	"github.com/ncarlier/za/pkg/logger"
	"github.com/ncarlier/za/pkg/outputs"
	"github.com/ncarlier/za/pkg/serializers"
)

const (
	defaultClientTimeout = 5 * time.Second
	defaultMethod        = http.MethodPost
	defaultUserAgent     = "Mozilla/5.0 (compatible; ZeroAnalytics/1.0; +https://github.com/ncarlier/za)"
)

// HTTP output
type HTTP struct {
	URL      string            `toml:"url"`
	Timeout  config.Duration   `toml:"timeout"`
	Method   string            `toml:"POST"`
	Username string            `toml:"username"`
	Password string            `toml:"password"`
	Headers  map[string]string `toml:"headers"`
	Gzip     bool              `toml:"gzip"`

	client     *http.Client
	serializer serializers.Serializer
}

var sampleConfig = `
  ## URL is the address to send events to
  url = "http://127.0.0.1:8080/"
  ## Timeout for HTTP message
  # timeout = "5s"
  ## HTTP method, one of: "POST" or "PUT"
  # method = "POST"
  ## HTTP Basic Auth credentials
  # username = "username"
  # password = "pa$$word"
  ## Compress body request using GZIP
  # gzip = true
  ## Data format to output.
  data_format = "json"
`

func (h *HTTP) createClient(ctx context.Context) (*http.Client, error) {
	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
		},
		Timeout: h.Timeout.Duration,
	}

	return client, nil
}

// SetSerializer set data serializer
func (h *HTTP) SetSerializer(serializer serializers.Serializer) {
	h.serializer = serializer
}

// Connect activate the output writer
func (h *HTTP) Connect() error {
	if h.URL == "" {
		return fmt.Errorf("invalid URL: %s", h.URL)
	}
	if h.Method == "" {
		h.Method = defaultMethod
	}
	h.Method = strings.ToUpper(h.Method)
	if h.Method != http.MethodPost && h.Method != http.MethodPut {
		return fmt.Errorf("invalid method [%s] %s", h.URL, h.Method)
	}

	if h.Timeout.Duration == 0 {
		h.Timeout.Duration = defaultClientTimeout
	}

	ctx := context.Background()
	client, err := h.createClient(ctx)
	if err != nil {
		return err
	}

	h.client = client

	logger.Debug.Printf("using HTTP output: %s\n", h.URL)
	return nil
}

// Close the output writer
func (h *HTTP) Close() error {
	return nil
}

// SampleConfig returns sample configuration
func (h *HTTP) SampleConfig() string {
	return sampleConfig
}

// Description returns description
func (h *HTTP) Description() string {
	return "Send page view to HTTP endpoint"
}

// SendEvent send event to the Output
func (h *HTTP) SendEvent(event events.Event) error {
	b, err := h.serializer.Serialize(event)
	if err != nil {
		return fmt.Errorf("unable to serialize page view: %v", err)
	}

	if err := h.send(b); err != nil {
		return err
	}

	return nil
}

func (h *HTTP) send(reqBody []byte) error {
	var reqBodyBuffer io.Reader = bytes.NewBuffer(reqBody)

	var err error
	if h.Gzip {
		rc, err := helper.CompressWithGzip(reqBodyBuffer)
		if err != nil {
			return err
		}
		defer rc.Close()
		reqBodyBuffer = rc
	}

	req, err := http.NewRequest(h.Method, h.URL, reqBodyBuffer)
	if err != nil {
		return err
	}

	if h.Username != "" || h.Password != "" {
		req.SetBasicAuth(h.Username, h.Password)
	}

	req.Header.Set("User-Agent", defaultUserAgent)
	req.Header.Set("Content-Type", h.serializer.ContentType())
	if h.Gzip {
		req.Header.Set("Content-Encoding", "gzip")
	}
	for k, v := range h.Headers {
		if strings.ToLower(k) == "host" {
			req.Host = v
		}
		req.Header.Set(k, v)
	}

	resp, err := h.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	_, err = ioutil.ReadAll(resp.Body)

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("bad status code (%d) when writing to [%s]", resp.StatusCode, h.URL)
	}

	return nil
}

func init() {
	outputs.Add("http", func() outputs.Output {
		return &HTTP{
			Timeout: config.Duration{Duration: defaultClientTimeout},
			Method:  defaultMethod,
			Gzip:    true,
		}
	})
}
