package handlers

import (
	"app/definitions"
	"fmt"
	"net/http"
	"time"
)

func WhatsAppHandler(w http.ResponseWriter, r *http.Request) {
	// Create a new custom writer.
	cw := definitions.CustomResponseWriter{}.New(w)

	// Loop to keep sending data
	for {
		select {
		case signal := <-definitions.WhatAppChan:
			// Stream the received signal payload to the client
			fmt.Printf("data: %v\n\n", signal)
			cw.StreamPayload(signal)

		case <-time.After(30 * time.Second):
			// Send a heartbeat every 30 seconds to keep the connection alive
			cw.Stream("heartbeat")

		case <-r.Context().Done():
			// Handle client disconnection or request cancellation
			fmt.Println("Client disconnected")
			return
		}
	}
}
