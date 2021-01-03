package http

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/ncarlier/za/pkg/events"
	"github.com/ncarlier/za/pkg/logger"
	"github.com/ncarlier/za/pkg/serializers/json"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func getEvent() events.Event {
	return &events.SimpleEvent{
		BaseEvent: events.BaseEvent{
			TrackingID:  "test",
			ClientIP:    "127.0.0.1",
			CountryCode: "fr",
			UserAgent:   "none",
			Tags: map[string]string{
				"foo": "bar",
			},
			Timestamp: time.Now(),
		},
		Payload: map[string]interface{}{
			"value": 42.0,
		},
	}
}

func TestInvalidURL(t *testing.T) {
	plugin := &HTTP{
		URL: "",
	}

	err := plugin.Connect()
	assert.Error(t, err)
}

func TestMethod(t *testing.T) {
	ts := httptest.NewServer(http.NotFoundHandler())
	defer ts.Close()

	u, err := url.Parse(fmt.Sprintf("http://%s", ts.Listener.Addr().String()))
	assert.Nil(t, err)

	tests := []struct {
		name           string
		plugin         *HTTP
		expectedMethod string
		connectError   bool
	}{
		{
			name: "default method is POST",
			plugin: &HTTP{
				URL:    u.String(),
				Method: defaultMethod,
			},
			expectedMethod: http.MethodPost,
		},
		{
			name: "put is okay",
			plugin: &HTTP{
				URL:    u.String(),
				Method: http.MethodPut,
			},
			expectedMethod: http.MethodPut,
		},
		{
			name: "get is invalid",
			plugin: &HTTP{
				URL:    u.String(),
				Method: http.MethodGet,
			},
			connectError: true,
		},
		{
			name: "method is case insensitive",
			plugin: &HTTP{
				URL:    u.String(),
				Method: "poST",
			},
			expectedMethod: http.MethodPost,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts.Config.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				require.Equal(t, tt.expectedMethod, r.Method)
				w.WriteHeader(http.StatusOK)
			})

			serializer, _ := json.NewSerializer()
			tt.plugin.SetSerializer(serializer)
			err = tt.plugin.Connect()
			if tt.connectError {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			err = tt.plugin.SendEvent(getEvent())
			require.NoError(t, err)
		})
	}
}

func init() {
	logger.Init("debug")
}
