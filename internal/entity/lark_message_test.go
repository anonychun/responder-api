package entity

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLarkMessage_IsReplied(t *testing.T) {
	t.Run("reports a message without a reply timestamp as not replied", func(t *testing.T) {
		larkMessage := &LarkMessage{}

		assert.False(t, larkMessage.IsReplied())
	})

	t.Run("reports a message with a reply timestamp as replied", func(t *testing.T) {
		larkMessage := &LarkMessage{RepliedAt: time.Now()}

		assert.True(t, larkMessage.IsReplied())
	})
}
