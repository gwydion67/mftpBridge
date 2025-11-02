package main

import (
	// "mftpBridge/src/backend"
	whatsapphandler "mftpBridge/src/whatsappHandler"
	"os"
	"os/signal"
	"syscall"
	"github.com/joho/godotenv"
	"go.mau.fi/whatsmeow/types"
)

func main() {
	// backend.StartBackend()
	var client whatsapphandler.MyClient
	
	godotenv.Load()

	groupJid := os.Getenv("MFTP_COMMUNITY_JID")
	if groupJid == "" {
		println("!ERROR group jid not found, please set the environment variables")
		return
	}

	whatsapphandler.Connect(&client)

	jid, err := types.ParseJID(groupJid)
	if err != nil {
		print("failed to parse jid ")
	} else {
		client.SendMessage("hello", jid)
	}

	// Listen to Ctrl+C (you can also do something else that prevents the program from exiting)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	client.Logout()
}
