package conditional

import (
	"maps"

	"github.com/ncarlier/za/pkg/events"
)

var ARGS map[string]interface{}

func init() {
	ARGS = (&events.PageView{}).ToMap()
	maps.Copy(ARGS, (&events.Exception{}).ToMap())
	maps.Copy(ARGS, (&events.SimpleEvent{}).ToMap())
}
