package main

import (
	"github.com/joho/godotenv"
	// "go.mau.fi/whatsmeow/types"
	"mftpBridge/src/backend"

	"mftpBridge/src/whatsappHandler"

	"os"
	"os/signal"
	"syscall"
)

func main() {

	var client whatsapphandler.MyClient

	godotenv.Load()

	groupJid := os.Getenv("MFTP_COMMUNITY_JID")
	if groupJid == "" {
		println("!ERROR group jid not found, please set the environment variables")
		println("To get the jid, inspect element on whatsapp web select the group and search for @g.us in the html")
		return
	}

	whatsapphandler.Connect(&client)

	// jid, err := types.ParseJID(groupJid)
	// if err != nil {
	// 	print("failed to parse jid ")
	// } else {
	// 	client.SendMessage("hello", jid)
	// }

	// after inital connection to whatsapp start the backend server and pass the client
	// use this client to manage the start/stop of the whatsapp server

	if client.WAClient == nil {
		println("Whatsapp Client creation failed")
		return
	}

	backend.StartBackend(&client)

	// Listen to Ctrl+C (you can also do something else that prevents the program from exiting)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	client.Logout()
}
