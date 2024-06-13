package expr

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/vm"
	"github.com/ncarlier/za/pkg/events"
)

var exprPlugins = map[string]interface{}{
	"toLower": strings.ToLower,
	"toUpper": strings.ToUpper,
}

// ConditionalExpression is a model for a conditional expression applied on an article
type ConditionalExpression struct {
	expression string
	prog       *vm.Program
}

// NewConditionalExpression creates a new conditional expression
func NewConditionalExpression(expression string, event events.Event) (*ConditionalExpression, error) {
	var prog *vm.Program
	if strings.TrimSpace(expression) != "" {
		args := buildExprArgs(event)

		options := []expr.Option{
			expr.Env(args),
			expr.AsBool(),
		}
		var err error
		prog, err = expr.Compile(expression, options...)
		if err != nil {
			return nil, fmt.Errorf("invalid conditional expression: %s", err.Error())
		}
	}
	return &ConditionalExpression{
		expression: expression,
		prog:       prog,
	}, nil
}

// Match test ifthe article match the conditional expression
func (c *ConditionalExpression) Match(event events.Event) bool {
	if c.prog == nil {
		return true
	}
	args := buildExprArgs(event)
	output, err := expr.Run(c.prog, args)
	if err != nil {
		slog.Error("unable to build expression arguments", "expr", c.expression)
		return false
	}
	return output.(bool)
}

// String returns the string expression
func (c *ConditionalExpression) String() string {
	return c.expression
}

func buildExprArgs(event events.Event) map[string]interface{} {
	env := event.ToMap()
	for k, v := range exprPlugins {
		env[k] = v
	}
	return env
}
