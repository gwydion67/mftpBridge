package backend

import (
	// "fmt"
	"encoding/json"
	"fmt"
	"io"
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

type Data struct {
	Message string `json:"message"`
}

func createHandler(client *whatsapphandler.MyClient) *Handler {
	return &Handler{client: client}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {

	groupJid := os.Getenv("MFTP_COMMUNITY_JID")

	body, err := io.ReadAll(req.Body)

	if err != nil {
		print("Req has no body")
		return
	}

	defer req.Body.Close()

	var receivedData Data
	err = json.Unmarshal(body, &receivedData)
	if err != nil {
		http.Error(w, "Error unmarshalling JSON", http.StatusBadRequest)
		return
	}

	// Process the received data (e.g., print it, save to a database).
	fmt.Printf("Received POST data: Message='%s'\n", receivedData.Message)
	message := receivedData.Message
	jid, err := types.ParseJID(groupJid)
	if err != nil {
		print("failed to parse jid ")
	} else {

		h.client.SendMessage(message, jid)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := map[string]string{"status": "success", "received_message": receivedData.Message}
	json.NewEncoder(w).Encode(response)
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
