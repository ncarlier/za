package test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ncarlier/za/pkg/conditional"
	"github.com/ncarlier/za/pkg/events"
)

func TesInvalidExpressionSyntax(t *testing.T) {
	conf := &conditional.Config{
		Condition: "###",
	}
	_, err := conditional.NewConditionalExpression(conf)
	assert.NotNil(t, err, "expression should not be valid")
}

func TestMatchingExpression(t *testing.T) {
	conf := &conditional.Config{
		Condition: "browser == \"foo\"",
	}
	condition, err := conditional.NewConditionalExpression(conf)
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
