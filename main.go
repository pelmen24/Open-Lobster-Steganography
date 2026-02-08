package main

import (
	"encoding/json"
	"fmt"
	"lobster-server/lobster"
	"log"
	"net/http"
	"strings"
)

type EncodeRequest struct {
	Text    string `json:"text"`
	Version int    `json:"version"`
	Key     string `json:"key,omitempty"`
}

type DecodeRequest struct {
	Encoded string `json:"encoded"`
	Version int    `json:"version"`
	Key     string `json:"key,omitempty"`
}

type Response struct {
	Success        bool   `json:"success"`
	Encoded        string `json:"encoded,omitempty"`
	Result         string `json:"result,omitempty"`
	Text           string `json:"text,omitempty"`
	LobsterCount   int    `json:"lobster_count"`
	OriginalLength int    `json:"original_length"`
	EncodedLength  int    `json:"encoded_length"`
	DecodedLength  int    `json:"decoded_length"`
	VersionName    string `json:"version_name"`
	Error          string `json:"error,omitempty"`
}

func enableCORS(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	(*w).Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

func handleEncode(w http.ResponseWriter, r *http.Request) {
	enableCORS(&w)
	if r.Method == "OPTIONS" {
		return
	}

	var req EncodeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var result string
	if req.Version <= 3 {
		result = lobster.EncodeTernary([]byte(req.Text))
	} else {
		result = lobster.EncodeBinary([]byte(req.Text), req.Key)
	}

	count := (strings.Count(result, "🦞")) / 2

	json.NewEncoder(w).Encode(Response{
		Success:        true,
		Encoded:        result,
		LobsterCount:   count,
		OriginalLength: len([]byte(req.Text)),
		EncodedLength:  len([]byte(result)),
		VersionName:    fmt.Sprintf("V%d", req.Version),
	})
}

func handleDecode(w http.ResponseWriter, r *http.Request) {
	enableCORS(&w)
	if r.Method == "OPTIONS" {
		return
	}

	var req DecodeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	text, err := lobster.Decode(req.Encoded, req.Version, req.Key)
	if err != nil {
		json.NewEncoder(w).Encode(Response{Success: false, Error: err.Error()})
		return
	}

	json.NewEncoder(w).Encode(Response{
		Success:       true,
		Result:        text,
		Text:          text,
		DecodedLength: len([]byte(text)),
		LobsterCount:  strings.Count(req.Encoded, "🦞"),
	})
}

func main() {
	http.HandleFunc("/api/encode", handleEncode)
	http.HandleFunc("/api/decode", handleDecode)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "🦞 OLD/EOLD API Server (Go Version) is running!")
	})

	fmt.Println("Starting server on :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
