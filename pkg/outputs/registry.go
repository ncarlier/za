package outputs

// Creator function for an output
type Creator func() Output

// Outputs registry
var Outputs = map[string]Creator{}

// Add output to the registry
func Add(name string, creator Creator) {
	Outputs[name] = creator
}
