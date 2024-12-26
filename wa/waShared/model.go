package waShared

import (
	"app/definitions"
	"encoding/json"
	"fmt"
)

func PrepareModel(Sender, Name, Time, MessageID, MessageType, MessageText, MessageCaption, Uri string) (string, error) {
	model := definitions.Model{
		Sender:         Sender,
		Name:           Name,
		Time:           Time,
		MessageID:      MessageID,
		MessageType:    MessageType,
		MessageText:    MessageText,
		MessageCaption: MessageCaption,
		Uri:            Uri,
	}

	data, err := json.Marshal(model)
	if err != nil {
		fmt.Println(err)
		return "failed to JSON", err
	}
	//	fmt.Printf("Model: %#v", model)
	return string(data), nil

}
