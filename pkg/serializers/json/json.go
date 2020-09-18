package json

import (
	"encoding/json"

	"github.com/ncarlier/trackr/pkg/events"
)

type serializer struct {
}

// NewSerializer create new JSON serializer
func NewSerializer() (*serializer, error) {
	s := &serializer{}
	return s, nil
}

func (s *serializer) Serialize(event events.Event) ([]byte, error) {
	serialized, err := json.Marshal(event)
	if err != nil {
		return []byte{}, err
	}
	serialized = append(serialized, '\n')

	return serialized, nil
}
