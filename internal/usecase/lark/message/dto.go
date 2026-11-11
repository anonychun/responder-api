package message

import "strings"

const supportTeamHandle = "one-segmentation-support"

type MessageReceiveRequest struct {
	MessageId string
	ChatId    string
	RootId    string
	ParentId  string
	SenderId  string
	Content   string
	Mentions  []string
}

func (r MessageReceiveRequest) IsRootMessage() bool {
	return r.RootId == "" && r.ParentId == ""
}

func (r MessageReceiveRequest) MentionsSupportTeam() bool {
	for _, mention := range r.Mentions {
		if strings.EqualFold(mention, supportTeamHandle) {
			return true
		}
	}

	return strings.Contains(strings.ToLower(r.Content), supportTeamHandle)
}
