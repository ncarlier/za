package outputs

import "github.com/ncarlier/trackr/pkg/model"

// Creator function for an output
type Creator func() model.Output

// Outputs registry
var Outputs = map[string]Creator{}

// Add output to the registry
func Add(name string, creator Creator) {
	Outputs[name] = creator
}
