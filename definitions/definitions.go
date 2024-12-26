package definitions

import (
	"net/http"

	"go.mau.fi/whatsmeow"
	waLog "go.mau.fi/whatsmeow/util/log"
)

type Payload struct {
	Event   string
	Message string
	Context <-chan struct{} // Cancellation context.
}

type CustomResponseWriter struct {
	http.ResponseWriter
	flusher http.Flusher
}

type Model struct {
	Sender, Name, Time, MessageID, MessageType, MessageText, MessageCaption, Uri string
}

var (
	InputChan   = make(chan Payload, 10) // Buffered channel.
	WhatAppChan = make(chan Payload, 10) // Buffered channel.
	OutputChan2 = make(chan Payload)
	OutputChan3 = make(chan Payload)

	Client *whatsmeow.Client
	Log    waLog.Logger
)

// 	signal := channels.Signal{
// 		ID:      1,
// 		Payload: "Hello, World! from backgroundTask",
// 	}

// 	channels.InputChan <- signal
// }
