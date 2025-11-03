package whatsapphandler

import (
	"context"
	"fmt"
	"mftpBridge/src/utils"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"github.com/mdp/qrterminal/v3"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"
)

type MyClient struct {
	WAClient       *whatsmeow.Client
	eventHandlerId uint32
}

func (mycli *MyClient) register() {
	mycli.eventHandlerId = mycli.WAClient.AddEventHandler(mycli.myEventHandler)
}

func (mycli *MyClient) myEventHandler(evt interface{}) {
	switch v := evt.(type) {
	case *events.Message:
		fmt.Println("Received a message!", v.Message.GetConversation())
		fmt.Println("From: ", v.Info.Sender.User)
		// file, err := os.Create("./jid.txt")
		// if err != nil {
		// 	fmt.Println("Error creating file")
		// }
		//
		// defer func() {
		// 	if closeErr := file.Close(); closeErr != nil {
		// 		fmt.Printf("Error closing file: %v\n", closeErr)
		// 	}
		// }()
		//
		// _, err = file.WriteString(v.Info.Chat.String())
		// if err != nil {
		// 	fmt.Printf("Error writing to file: %v\n", err)
		// 	return
		// }
		//
		// fmt.Printf("Successfully wrote to %s\n", "./jid.txt")
		// mycli.WAClient.SendMessage(context.Background(), v.Info.Chat, utils.TextToWaMessage("hi"))
	}
}

func (mycli *MyClient) Logout() {
	mycli.WAClient.Disconnect()
}

func (mycli *MyClient) SendMessage(msg string, jid types.JID) {
	mycli.WAClient.SendMessage(context.Background(), jid, utils.TextToWaMessage(msg))
}

func Connect(client *MyClient) {

	dbLog := waLog.Stdout("Database", "ERROR", true)
	ctx := context.Background()
	println("Conneting to db")

	dataDB := os.Getenv("DB_DIRECTORY")
	if dataDB == "" {
		dataDB = "."
		fmt.Println("Warning: DB_DIRECTORY environment variable is not set. Using current directoy")
	}

	container, err := sqlstore.New(ctx, "sqlite3", dataDB + "/whatsapp.db?_foreign_keys=on", dbLog)
	if err != nil {
		panic(err)
	}
	// If you want multiple sessions, remember their JIDs and use .GetDevice(jid) or .GetAllDevices() instead.
	deviceStore, err := container.GetFirstDevice(ctx)
	if err != nil {
		panic(err)
	}

	clientLog := waLog.Stdout("Client", "ERROR", true)

	client.WAClient = whatsmeow.NewClient(deviceStore, clientLog)
	client.WAClient.EnableAutoReconnect = true
	client.register()

	if client.WAClient.Store.ID == nil {
		// No ID stored, new login
		qrChan, _ := client.WAClient.GetQRChannel(context.Background())
		err = client.WAClient.Connect()
		if err != nil {
			panic(err)
		}
		for evt := range qrChan {
			if evt.Event == "code" {
				// Render the QR code here
				qrterminal.GenerateHalfBlock(evt.Code, qrterminal.L, os.Stdout)
				// or just manually `echo 2@... | qrencode -t ansiutf8` in a terminal
				fmt.Println("QR code:", evt.Code)
			} else {
				fmt.Println("Login event:", evt.Event)
			}
		}
	} else {
		// Already logged in, just connect
		err = client.WAClient.Connect()
		if err != nil {
			panic(err)
		}
	}

}
