package api

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/ncarlier/za/pkg/config"
	"github.com/stretchr/testify/assert"
)

var conf *config.Config

func newTestRequest(values url.Values) (*http.Request, error) {
	req, err := http.NewRequest("GET", "/collect", http.NoBody)
	if err != nil {
		return nil, err
	}
	req.URL.RawQuery = values.Encode()
	req.Header.Set("Referer", "http://localhost:8080")
	return req, nil
}

func newSimpleEvent() url.Values {
	payload := url.Values{}
	payload.Set("tid", "test")
	payload.Set("t", "event")
	return payload
}

func TestCollectHandlerWithDNT(t *testing.T) {
	rr := httptest.NewRecorder()
	handler := collectHandler(nil, conf)

	req, err := newTestRequest(newSimpleEvent())
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("DNT", "1")

	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "N", rr.Header().Get("Tk"))
	assert.Equal(t, "image/gif", rr.Header().Get("Content-Type"))
}

func TestCollectHandlerWithPrefetch(t *testing.T) {
	rr := httptest.NewRecorder()
	handler := collectHandler(nil, conf)

	req, err := newTestRequest(newSimpleEvent())
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("X-Moz", "prefetch")

	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusNoContent, rr.Code)
}

func TestCollectHandlerWithBot(t *testing.T) {
	rr := httptest.NewRecorder()
	handler := collectHandler(nil, conf)

	req, err := newTestRequest(newSimpleEvent())
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)")

	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusNoContent, rr.Code)
}

func TestCollectHandlerWithInvalidReferer(t *testing.T) {
	rr := httptest.NewRecorder()
	handler := collectHandler(nil, conf)

	req, err := newTestRequest(newSimpleEvent())
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Referer", "http://example.com")

	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestCollectHandlerWithBadge(t *testing.T) {
	rr := httptest.NewRecorder()
	handler := collectHandler(nil, conf)

	event := newSimpleEvent()
	event.Set("t", "badge")
	req, err := newTestRequest(event)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Referer", "http://example.com")

	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "P", rr.Header().Get("Tk"))
	assert.Equal(t, "image/svg+xml", rr.Header().Get("Content-Type"))
}

func TestCollectHandler(t *testing.T) {
	tt := []struct {
		name  string
		query string
		want  string
		code  int
	}{
		{
			name:  "simple pageview event",
			query: "t=pageview&tid=test",
			want:  "P",
			code:  http.StatusOK,
		},
		{
			name:  "simple exception event",
			query: "t=exception&tid=test&exl=1&exc=1",
			want:  "P",
			code:  http.StatusOK,
		},
		{
			name:  "bad exception event",
			query: "t=exception&tid=test&exl=1",
			code:  http.StatusBadRequest,
		},
		{
			name:  "bad type event",
			query: "t=foo&tid=test",
			code:  http.StatusBadRequest,
		},
		{
			name:  "custom event with valid payload",
			query: "t=event&tid=test&d=eyJmb28iOiJiYXIifQ==",
			want:  "P",
			code:  http.StatusOK,
		},
		{
			name:  "custom event with invalid payload",
			query: "t=event&tid=test&d=xxx",
			code:  http.StatusBadRequest,
		},
	}

	handler := collectHandler(nil, conf)
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			payload, err := url.ParseQuery(tc.query)
			if err != nil {
				t.Fatal(err)
			}
			req, err := newTestRequest(payload)
			if err != nil {
				t.Fatal(err)
			}

			handler.ServeHTTP(rr, req)
			assert.Equal(t, tc.code, rr.Code)
			if tc.code == http.StatusOK {
				assert.Equal(t, tc.want, rr.Header().Get("Tk"))
				assert.Equal(t, "image/gif", rr.Header().Get("Content-Type"))
			}
		})
	}
}

func init() {
	conf = config.NewConfig()
	conf.Trackers = append(conf.Trackers, config.TrackerConfig{
		Origin:     "http://localhost:8080",
		TrackingID: "test",
		Badge:      "ZerØ|analytics|#00a5da",
	})
}
