package json

import (
	"encoding/json"

	"github.com/ncarlier/trackr/pkg/model"
)

type serializer struct {
}

// NewSerializer create new JSON serializer
func NewSerializer() (*serializer, error) {
	s := &serializer{}
	return s, nil
}

func (s *serializer) Serialize(pageview model.PageView) ([]byte, error) {
	serialized, err := json.Marshal(pageview)
	if err != nil {
		return []byte{}, err
	}
	serialized = append(serialized, '\n')

	return serialized, nil
}
