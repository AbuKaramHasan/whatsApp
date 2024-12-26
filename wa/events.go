package wa

import (
	"app/definitions"
	"app/wa/waHandler"
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"
	//	"knowlege-navigator/msgHandler"
	//	"knowlege-navigator/websocket"
	//
	// Replace with your actual module name
)

func eventHandler(evt interface{}) {
	switch v := evt.(type) {
	case *events.Message:
		fmt.Println("Received a message:", v.Message.GetConversation())
		waHandler.HandelMessage(v)
	case *events.Disconnected:
		fmt.Println("Disconnected")
	}
}

func Run() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigChan
		fmt.Println("Shutting down...")
		definitions.InputChan <- definitions.Payload{
			Event:   "notification",
			Message: "Server is shutting down",
		}

		time.Sleep(1 * time.Second)
		os.Exit(0)
	}()

	dbLog := waLog.Stdout("Database", "DEBUG", true)
	container, err := sqlstore.New("sqlite3", "file:examplestore.db?_foreign_keys=on", dbLog)
	if err != nil {
		panic(err)
	}

	deviceStore, err := container.GetFirstDevice()
	if err != nil {
		panic(err)
	}

	clientLog := waLog.Stdout("Client", "DEBUG", true)
	client := whatsmeow.NewClient(deviceStore, clientLog)
	client.AddEventHandler(eventHandler)

	definitions.Client = client

	for {
		if client.Store.ID == nil {
			// New login
			qrChan, _ := client.GetQRChannel(context.Background())
			err = client.Connect()
			if err != nil {
				fmt.Println("Failed to connect:", err)
				time.Sleep(5 * time.Second)
				continue
			}

			fmt.Println("Waiting for QR code scan...")
			for evt := range qrChan {
				if evt.Event == "code" {
					definitions.InputChan <- definitions.Payload{
						Event:   "qrCode",
						Message: evt.Code,
					}
					//	qrterminal.GenerateHalfBlock(evt.Code, qrterminal.L, os.Stdout)
				} else if evt.Event == "timeout" {
					fmt.Println("QR code timeout, retrying...")
					client.Disconnect()
					break
				} else {
					fmt.Println("Login event:", evt.Event)
				}
			}
		} else {
			// Reconnect
			err = client.Connect()
			if err != nil {
				fmt.Println("Failed to reconnect:", err)
				time.Sleep(5 * time.Second)
				continue
			}

			fmt.Println("Client connected and logged in")
			if definitions.Client != nil {
				waHandler.Send("966569291028")
			} else {
				fmt.Println("websocket.Client is nil; cannot send disposable message")
			}
			break
		}
	}

	<-sigChan
	client.Disconnect()
	fmt.Println("Client disconnected")
}
