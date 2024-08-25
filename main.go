package main

import (
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"net/http"
)

func WebSocketUpgradeHeader(req *http.Request) (string, error) {
	key := req.Header.Get("Sec-WebSocket-Key")
	if key == "" {
		return "", fmt.Errorf("Sec-WebSocket-Key is missing...")
	}

	acceptedKey := ComputeAcceptKey(key)
	return acceptedKey, nil

}

func ComputeAcceptKey(key string) string {
	const webSocketGUID = "258EAFA5-E914-47DA-95CA-C5AB0DC85B11"
	hash := sha1.New()
	hash.Write([]byte(key + webSocketGUID))
	return base64.StdEncoding.EncodeToString(hash.Sum(nil))
}

func HandleConnection(w http.ResponseWriter, r *http.Request) {

	acceptedKey, err := WebSocketUpgradeHeader(r)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	w.Header().Set("Upgrade", "websocket")
	w.Header().Set("Connection", "Upgrade")
	w.Header().Set("Sec-WebSocket-Accept", acceptedKey)

	// establish websocket connection
	conn, _, err := w.(http.Hijacker).Hijack()

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	defer conn.Close()

	fmt.Println("Websocket connection established")

}

func main() {

	http.HandleFunc("/ws", HandleConnection)

	serverAddress := "localhost:3000"

	fmt.Printf("Websocket server started on %s\n", serverAddress)

	err := http.ListenAndServe(serverAddress, nil)

	if err != nil {
		fmt.Println("Error starting server", err)
	}
}
