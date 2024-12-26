package waHandler

import (
	"app/definitions"
	"app/wa/waShared"
	"context"
	"fmt"
	"log"

	"go.mau.fi/whatsmeow/proto/waE2E"
	"google.golang.org/protobuf/proto"
)

func Send(sender string) {
	fmt.Println(sender)
	// Check if the client is initialized
	if definitions.Client == nil {
		log.Println("websocket.Client is not initialized")
		return
	}

	// Prepare the message
	m := &waE2E.Message{
		Conversation: proto.String("*1234* is your verification code for Plastbau Platform."),
	}

	// content, err := os.ReadFile("./logo.png")
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// resp, err := websocket.Client.Upload(context.Background(), content, whatsmeow.MediaImage)
	// if err != nil {
	// 	fmt.Println(err)
	// }

	// msg := &waE2E.ImageMessage{
	// 	Caption:  proto.String("Do you see any text other than the logo?"),
	// 	Mimetype: proto.String("image/png"), // replace this with the actual mime type
	// 	// you can also optionally add other fields like ContextInfo and JpegThumbnail here

	// 	URL:           &resp.URL,
	// 	DirectPath:    &resp.DirectPath,
	// 	MediaKey:      resp.MediaKey,
	// 	FileEncSHA256: resp.FileEncSHA256,
	// 	FileSHA256:    resp.FileSHA256,
	// 	FileLength:    &resp.FileLength,
	// 	ViewOnce:      proto.Bool(true),
	// }

	// m := &waE2E.Message{ImageMessage: msg}
	// // Parse the JID
	jid, ok := waShared.ParseJID(sender)
	if !ok {
		log.Printf("Failed to parse JID: %s", sender)
		return // Safely exit if JID parsing fails
	}

	fmt.Printf("sent to: %s", jid.String())

	// Log.Infof("Parsed JID: %v", jid)

	// // Attempt to send the message
	resp, err := definitions.Client.SendMessage(context.Background(), jid, m)
	if err != nil {
		fmt.Println(err)
		// Log.Errorf("Error sending message to JID %v", err)
		return
	}
	fmt.Println(resp.Timestamp) // time.Now()

	fmt.Printf("Message sent to %s\n", jid.String())

	data, _ := waShared.PrepareModel(sender, "self", resp.Timestamp.Local().Format("Mon 02-Jan-2006 15:04"),
		jid.String(), "text", fmt.Sprintf("%v", m), "", "")
	fmt.Println(data)

	definitions.InputChan <- definitions.Payload{
		Event:   "message", // default: source.onmessage = function (event) {}
		Message: data,
	}

}
