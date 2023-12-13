package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	qrcode "github.com/skip2/go-qrcode"
)

type response struct {
	Message string `json:"message"`
}

func helloWorldHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response{Message: "Hello, World"})
}

func createQRCodeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	data := r.URL.Query().Get("data")
	if data == "" {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response{Message: "Data is required to create the QR code."})
		return
	}

	size := r.URL.Query().Get("size")
	qrcodeSize := 256
	if size != "" {
		var err error
		qrcodeSize, err = strconv.Atoi(size)
		const maxSize = 2048
		if qrcodeSize > maxSize {
			qrcodeSize = maxSize
		}
		if err != nil {
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(response{Message: "Size in invalid."})
			return
		}
	}

	png, err := qrcode.Encode(data, qrcode.Medium, qrcodeSize)
	if err != nil {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response{Message: "Failed to create the QR code from the given data."})
		return
	}

	w.Header().Add("Content-Type", "image/png")
	w.Write(png)
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/", helloWorldHandler).Methods(http.MethodGet, http.MethodOptions)
	r.Handle("/qrcode", handlers.LoggingHandler(os.Stdout, http.HandlerFunc(createQRCodeHandler))).Methods(http.MethodGet, http.MethodOptions)
	r.Use(mux.CORSMethodMiddleware(r))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Listening on port :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, handlers.CompressHandler(r)))
}
