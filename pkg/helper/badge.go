package helper

import (
	"fmt"
	"strings"
	"sync"

	"github.com/narqo/go-badge"
)

var badgeCache = map[string][]byte{}
var lock = sync.Mutex{}

// GetBadge transform badge code into SVG badge
func GetBadge(code string) ([]byte, error) {
	if data, ok := badgeCache[code]; ok {
		return data, nil
	}

	parts := strings.Split(code, "|")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid badge code: %s", code)
	}
	data, err := badge.RenderBytes(parts[0], parts[1], badge.Color(parts[2]))
	if err != nil {
		return nil, err
	}
	lock.Lock()
	defer lock.Unlock()
	badgeCache[code] = data
	return data, nil
}
