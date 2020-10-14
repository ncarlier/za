package test

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ncarlier/za/pkg/geoip"
)

func TestLookupCountry(t *testing.T) {
	t.Skip()
	db, err := geoip.New("../../../var/dbip-city-lite.mmdb")
	assert.Nil(t, err, "unable to load Geo IP database")
	defer db.Close()
	ip := net.ParseIP("8.8.8.8")
	country, err := db.LookupCountry(ip)
	assert.Nil(t, err, "unable to get country code")
	assert.Equal(t, "US", country)
}
