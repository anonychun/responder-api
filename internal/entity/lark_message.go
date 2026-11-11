package entity

import "time"

type LarkMessage struct {
	Base

	MessageId string
	ChatId    string
	Content   string
	SenderId  string
	RepliedAt time.Time `bun:"replied_at,nullzero"`
}

func (m *LarkMessage) IsReplied() bool {
	return !m.RepliedAt.IsZero()
}
