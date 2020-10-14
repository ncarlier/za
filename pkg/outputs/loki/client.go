package loki

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/ncarlier/za/pkg/outputs/loki/logproto"

	"github.com/gogo/protobuf/proto"
	"github.com/golang/snappy"
)

const maxBodyResponseSize = 1024

const maxBackoffRetry = 3

// Config for Loki client
type Config struct {
	URL     string
	Timeout time.Duration
}

// Client for Loki
type Client struct {
	cfg Config
}

// NewClient create new Loki client
func NewClient(cfg Config) *Client {
	return &Client{
		cfg: cfg,
	}
}

// Send streams to Loki instance
func (c *Client) Send(streams []*logproto.Stream) (err error) {
	payload, err := c.encode(streams)
	if err != nil {
		return
	}

	ctx := context.Background()
	b := &Backoff{
		Min:    100 * time.Millisecond,
		Max:    5 * time.Second,
		Factor: 2,
		Jitter: false,
	}
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		retry := false
		for {
			retry, err = c.send(ctx, payload)
			if err == nil || !retry || b.Attempt() >= maxBackoffRetry {
				return
			}
			time.Sleep(b.Duration())
		}
	}()
	wg.Wait()
	return
}

func (c *Client) send(ctx context.Context, buf []byte) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, c.cfg.Timeout)
	defer cancel()
	req, err := http.NewRequest("POST", c.cfg.URL, bytes.NewReader(buf))
	if err != nil {
		return false, err
	}
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-protobuf")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode/100 != 2 {
		scanner := bufio.NewScanner(io.LimitReader(resp.Body, maxBodyResponseSize))
		body := ""
		if scanner.Scan() {
			body = scanner.Text()
		}
		return true, fmt.Errorf("bad status response: %s (%d) - %s", resp.Status, resp.StatusCode, body)
	}
	return false, nil
}

func (c *Client) encode(streams []*logproto.Stream) ([]byte, error) {
	req := logproto.PushRequest{
		Streams: streams,
	}
	buf, err := proto.Marshal(&req)
	if err != nil {
		return nil, err
	}
	buf = snappy.Encode(nil, buf)
	return buf, nil
}
