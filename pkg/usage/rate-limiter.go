package usage

import (
	"time"

	"github.com/sethvargo/go-limiter"
	"github.com/sethvargo/go-limiter/memorystore"
	"github.com/sethvargo/go-limiter/noopstore"
)

// RateLimiter is an interface use to apply rate limiting
type RateLimiter limiter.Store

// RateLimiterConfig for rate-limiter configuration section
type RateLimitingConfig struct {
	// Provider of the rate limiting store
	Provider string `toml:"provider"`
	// Tokens allowed per interval
	Tokens int `toml:"tokens"`
	// Interval until tokens reset
	Interval string `toml:"interval"`
}

// NewRateLimiter create new rate limiter
func NewRateLimiter(conf RateLimitingConfig) (RateLimiter, error) {
	switch conf.Provider {
	case "memory":
		interval, err := time.ParseDuration(conf.Interval)
		if err != nil {
			return nil, err
		}
		store, err := memorystore.New(&memorystore.Config{
			Tokens:   uint64(conf.Tokens),
			Interval: interval,
		})
		if err != nil {
			return nil, err
		}
		return store, nil
	default:
		store, err := noopstore.New()
		if err != nil {
			return nil, err
		}
		return store, nil
	}
}
