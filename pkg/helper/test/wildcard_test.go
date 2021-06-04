package test

import (
	"testing"

	"github.com/ncarlier/za/pkg/helper"
	"github.com/stretchr/testify/assert"
)

func TestWildcardPatterns(t *testing.T) {
	assert.True(t, helper.Match("https://example.com", "https://example.com"))
	assert.True(t, helper.Match("https://example.com", "https://example.com/foo"))
	assert.False(t, helper.Match("https://*.example.com", "https://example.com"))
	assert.True(t, helper.Match("https://*.example.com", "https://www.example.com"))
	assert.True(t, helper.Match("https://*example.com", "https://www.example.com"))
	assert.True(t, helper.Match("https://example.com/*", "https://example.com/"))
	assert.True(t, helper.Match("https://example.com/*", "https://example.com/foo"))
	// assert.True(t, helper.Match("http?://*.example.com", "http://foo.example.com"))
	assert.True(t, helper.Match("http?://*.*.example.com", "https://foo.bar.example.com"))
}
