package conditional

import (
	"fmt"
	"log/slog"
	"maps"
	"strings"

	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/vm"
	"github.com/ncarlier/za/pkg/events"
)

var exprPlugins = map[string]interface{}{
	"toLower": strings.ToLower,
	"toUpper": strings.ToUpper,
}

// internalExpression is a model for a conditional expression applied on an article
type intenalExpression struct {
	expression string
	prog       *vm.Program
}

// newConditionalExpression creates a new conditional expression
func newConditionalExpression(expression string, input map[string]interface{}) (Expression, error) {
	var prog *vm.Program
	if strings.TrimSpace(expression) != "" {
		args := make(map[string]interface{})
		maps.Copy(args, input)
		maps.Copy(args, exprPlugins)

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
	return &intenalExpression{
		expression: expression,
		prog:       prog,
	}, nil
}

// Match test if the event match the conditional expression
func (c *intenalExpression) Match(event events.Event) bool {
	if c.prog == nil {
		return true
	}
	output, err := expr.Run(c.prog, event.ToMap())
	if err != nil {
		slog.Error("unable to build expression arguments", "expr", c.expression)
		return false
	}
	return output.(bool)
}

// String returns the string expression
func (c *intenalExpression) String() string {
	return c.expression
}
