package test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ncarlier/za/pkg/events"
	"github.com/ncarlier/za/pkg/expr"
)

func TesInvalidExpressionSyntax(t *testing.T) {
	ev := &events.PageView{}
	_, err := expr.NewConditionalExpression("###", ev)
	assert.NotNil(t, err, "expression should not be valid")
}

func TestMatchingExpression(t *testing.T) {
	ev := &events.PageView{}
	condition, err := expr.NewConditionalExpression("browser == \"foo\"", ev)
	assert.Nil(t, err, "expression should be valid")
	assert.NotNil(t, condition)

	event := &events.PageView{
		BaseEvent: events.BaseEvent{
			Browser: "foo",
		},
	}

	ok := condition.Match(event)
	assert.True(t, ok, "event should match")
}
