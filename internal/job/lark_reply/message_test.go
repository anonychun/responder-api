package lark_reply

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReplyContent(t *testing.T) {
	t.Run("greets the sender with a mention and includes the guide", func(t *testing.T) {
		content := replyContent("ou_1")

		assert.Contains(t, content, `Hello <at user_id="ou_1"></at>! 👋 `)
		assert.Contains(t, content, "One Segmentation Support Team")
		assert.Contains(t, content, "https://gotocompany.sg.larksuite.com/wiki/LQ5qwjOqqi1JDakRuvhly3JygKh")
	})

	t.Run("greets without a mention when the sender is unknown", func(t *testing.T) {
		content := replyContent("")

		assert.Contains(t, content, "Hello! 👋 ")
		assert.NotContains(t, content, "<at ")
	})
}
