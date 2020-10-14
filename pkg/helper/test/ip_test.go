package test

import (
	"net/http"
	"testing"

	"github.com/ncarlier/za/pkg/helper"

	"github.com/stretchr/testify/assert"
)

func TestParseIPClient(t *testing.T) {
	req := &http.Request{}
	req.Header = make(http.Header)
	req.Header.Set("x-forwarded-for", "90.80.70.60, 127.0.0.1")
	ip := helper.ParseClientIP(req)
	assert.Equal(t, "90.80.70.60", ip)
}
