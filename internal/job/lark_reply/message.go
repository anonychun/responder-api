package lark_reply

import "fmt"

const replyBody = `Thank you for reaching out to the One Segmentation Support Team.

💡 Before we jump in, you may find an immediate solution in our quick reference guides. They cover the most common questions and troubleshooting steps:

- [📚 Frequently Asked Questions (FAQ)](https://gotocompany.sg.larksuite.com/wiki/LQ5qwjOqqi1JDakRuvhly3JygKh?fromScene=spaceOverview)
- [⚙️ Quick Setup & Troubleshooting Guide](https://gotocompany.sg.larksuite.com/wiki/QyHOwFRN6iOZcrkzIgpl3Hyvgky?fromScene=spaceOverview)

If these resources don't resolve your issue, please share any additional details in this thread. Our team will review your request and assist you here as soon as possible.

Thank you for your patience and cooperation!
`

func replyContent(senderId string) string {
	greeting := "Hello"
	if senderId != "" {
		greeting = fmt.Sprintf("Hello <at user_id=%q></at>", senderId)
	}

	return greeting + "! 👋 " + replyBody
}
