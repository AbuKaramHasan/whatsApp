package main

import (
	"app/definitions"
	"app/handlers"
	"app/wa"
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go Broadcast()
	go wa.Run()

	server := &http.Server{Addr: ":8080", Handler: http.DefaultServeMux}
	http.HandleFunc("/wa", handlers.WhatsAppHandler)

	// Serve static files from the "templates/static" directory
	fs := http.FileServer(http.Dir("templates/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", handlers.Index)

	go func() {
		fmt.Println("Starting server on :8080...")
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			fmt.Printf("Server error: %v\n", err)
			stop()
		}
	}()

	<-ctx.Done() // Wait for shutdown signal.
	fmt.Println("Shutting down server...")
	server.Shutdown(context.Background())
}

// Broadcast sends signals to different channels based on the signal ID.
func Broadcast() {
	for signal := range definitions.InputChan {
		select {
		case <-signal.Context:
			fmt.Println("Request has been canceled.")
			continue
		default:
			fmt.Println("Signal received from InputChan.", signal.Message)
			switch signal.Event {
			case "qrCode":
				definitions.WhatAppChan <- signal
			case "message":
				definitions.WhatAppChan <- signal
			default:
				definitions.OutputChan3 <- signal
			}
		}
	}
}
