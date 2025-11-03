package backend

import (
	// "fmt"
	"encoding/json"
	"fmt"
	"io"
	whatsapphandler "mftpBridge/src/whatsappHandler"
	"net/http"
	"os"
	"strings"

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
	jid    string
}

type TextData struct {
	Message string `json:"message"`
}

func createHandler(client *whatsapphandler.MyClient) *Handler {
	return &Handler{
		client: client,
		jid:    os.Getenv("MFTP_COMMUNITY_JID"),
	}
}

func (h *Handler) handleSendText(w http.ResponseWriter, req *http.Request) {

	body, err := io.ReadAll(req.Body)

	if err != nil {
		print("Req has no body")
		return
	}

	defer req.Body.Close()

	var receivedData TextData
	err = json.Unmarshal(body, &receivedData)
	if err != nil {
		http.Error(w, "Error unmarshalling JSON", http.StatusBadRequest)
		return
	}

	// Process the received data (e.g., print it, save to a database).
	fmt.Printf("Received POST data: Message='%s'\n", receivedData.Message)
	message := receivedData.Message
	jid, err := types.ParseJID(h.jid)
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

func (h *Handler) handleSendDoc(w http.ResponseWriter, req *http.Request) {

	contentType := req.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "multipart/form-data") {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		response := map[string]string{"status": "Failed", "Error": fmt.Sprintf("Invalid data, required multipart form got, %s", contentType)}
		json.NewEncoder(w).Encode(response)
		return
	}

	req.ParseMultipartForm(10 << 20)

	file, handler, err := req.FormFile("attachment")
	message := req.FormValue("message")
	if err != nil {
		http.Error(w, "Error retrieving the file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	fmt.Fprintf(w, "Uploaded File: %s\n", handler.Filename)
	fmt.Fprintf(w, "File Size: %d\n", handler.Size)
	fmt.Fprintf(w, "MIME Header: %v\n", handler.Header)

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Invalid file", http.StatusBadRequest)
		return
	}

	jid, err := types.ParseJID(h.jid)
	if err != nil {
		print("failed to parse jid ")
	} else {

		h.client.SendDocument(jid, fileBytes, handler.Filename, message)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := map[string]string{"status": "success", "received_message": message}
	json.NewEncoder(w).Encode(response)
}

func StartBackend(client *whatsapphandler.MyClient) {

	// http.HandleFunc("/hello", hello)
	// http.HandleFunc("GET /headers", headers)

	router := http.NewServeMux()
	handler := createHandler(client)
	router.HandleFunc("/sendtext", handler.handleSendText)
	router.HandleFunc("/senddoc", handler.handleSendDoc)

	server := http.Server{
		Addr:    ":8090",
		Handler: router,
	}

	fmt.Println("Server listening on port 8090")
	server.ListenAndServe()
}
