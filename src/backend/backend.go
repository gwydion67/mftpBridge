package backend

import (
	// "fmt"
	"fmt"
	whatsapphandler "mftpBridge/src/whatsappHandler"
	"net/http"
	"os"

	"go.mau.fi/whatsmeow/types"
)

// func hello(w http.ResponseWriter, req *http.Request) {
//
// 	fmt.Fprintf(w, "hello\n")
// }
//
// func headers(w http.ResponseWriter, req *http.Request) {
//
// 	for name, headers := range req.Header {
// 		for _, h := range headers {
// 			fmt.Fprintf(w, "%v: %v\n", name, h)
// 		}
// 	}
// }

type Handler struct {
	client *whatsapphandler.MyClient
}

func createHandler(client *whatsapphandler.MyClient) *Handler {
	return &Handler{client: client}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {

	groupJid := os.Getenv("MFTP_COMMUNITY_JID")
	jid, err := types.ParseJID(groupJid)
	if err != nil {
		print("failed to parse jid ")
	} else {
		h.client.SendMessage("hello", jid)
	}
}

func StartBackend(client *whatsapphandler.MyClient) {

	// http.HandleFunc("/hello", hello)
	// http.HandleFunc("GET /headers", headers)

	router := http.NewServeMux()
	handler := createHandler(client)
	router.Handle("/send", handler)

	server := http.Server{
		Addr:    ":8090",
		Handler: router,
	}

	fmt.Println("Server listening on port 8090")
	server.ListenAndServe()
}
