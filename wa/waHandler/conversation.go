package waHandler

import (
	"app/definitions"
	"app/wa/waShared"
	"fmt"

	"go.mau.fi/whatsmeow/types/events"
)

func Conversation(evt *events.Message) {
	fmt.Println(evt.Message)
	sender := evt.Info.Chat.User
	pushName := evt.Info.PushName

	msgReceived := evt.Message.GetConversation()
	data, _ := waShared.PrepareModel(sender, pushName, evt.Info.Timestamp.Local().Format("Mon 02-Jan-2006 15:04"),
		evt.Info.ID, "text", fmt.Sprintf("%v", msgReceived), "", "")
	fmt.Println(data)

	definitions.InputChan <- definitions.Payload{
		Event:   "message", // default: source.onmessage = function (event) {}
		Message: data,
	}
}
