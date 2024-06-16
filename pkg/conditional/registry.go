package conditional

import "github.com/ncarlier/za/pkg/events"

// Output is an interface for output plugins that are conditioned to an expression.
type Output interface {
	// SetCondition sets the condition expression for the interface.
	SetCondition(condition Expression)
}

// Serializer is an interface defining functions that a serializer plugin must satisfy.
type Expression interface {
	// Match test if the event match the conditional expression
	Match(event events.Event) bool
	// String returns condition as text
	String() string
}

// Config is a struct that define output condition expression.
type Config struct {
	// Condition to produce metric
	// See syntax here: https://expr-lang.org/docs/language-definition
	Condition string `toml:"condition"`
}

// NewConditionalExpression creat a conditional Output interface based on the given config.
func NewConditionalExpression(config *Config) (Expression, error) {
	return newConditionalExpression(config.Condition, ARGS)
}
