package lark

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTextMessageBody(t *testing.T) {
	t.Run("escapes quotes and line breaks into valid JSON", func(t *testing.T) {
		text := "Hello <at user_id=\"ou_1\"></at>! 👋\n\n- [FAQ](https://example.com)\nThank you!"

		body, err := textMessageBody(text)

		require.NoError(t, err)
		require.True(t, json.Valid([]byte(body)), "body is not valid JSON: %s", body)

		var message struct {
			Text string `json:"text"`
		}
		require.NoError(t, json.Unmarshal([]byte(body), &message))
		assert.Equal(t, text, message.Text)
	})
}
