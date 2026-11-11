package message

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMessageReceiveRequest_IsRootMessage(t *testing.T) {
	t.Run("reports a message without a thread reference as a root message", func(t *testing.T) {
		req := MessageReceiveRequest{MessageId: "om_1"}

		assert.True(t, req.IsRootMessage())
	})

	t.Run("reports a reply as not a root message", func(t *testing.T) {
		req := MessageReceiveRequest{MessageId: "om_2", RootId: "om_1"}

		assert.False(t, req.IsRootMessage())
	})

	t.Run("reports a reply that only carries a parent as not a root message", func(t *testing.T) {
		req := MessageReceiveRequest{MessageId: "om_2", ParentId: "om_1"}

		assert.False(t, req.IsRootMessage())
	})
}

func TestMessageReceiveRequest_MentionsSupportTeam(t *testing.T) {
	t.Run("matches a mention of the support team handle", func(t *testing.T) {
		req := MessageReceiveRequest{
			Content:  `{"text":"@_user_1 please help"}`,
			Mentions: []string{"someone", "one-segmentation-support"},
		}

		assert.True(t, req.MentionsSupportTeam())
	})

	t.Run("matches the handle typed as plain text", func(t *testing.T) {
		req := MessageReceiveRequest{Content: `{"text":"one-segmentation-support please help"}`}

		assert.True(t, req.MentionsSupportTeam())
	})

	t.Run("does not match other mentions", func(t *testing.T) {
		req := MessageReceiveRequest{
			Content:  `{"text":"@_user_1 please help"}`,
			Mentions: []string{"someone"},
		}

		assert.False(t, req.MentionsSupportTeam())
	})
}
