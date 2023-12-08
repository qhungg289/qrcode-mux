package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	qrcode "github.com/skip2/go-qrcode"
)

type Response struct {
	Message string `json:"message"`
}

func helloWorldHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	res := Response{Message: "Hello, World"}
	b, _ := json.Marshal(res)

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, string(b))
}

func createQRCodeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	data := r.URL.Query().Get("data")
	if data == "" {
		res := Response{Message: "Data is required to create the QR code."}
		b, _ := json.Marshal(res)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, string(b))
		return
	}

	size := r.URL.Query().Get("size")
	qrcodeSize := 256
	if size != "" {
		var err error
		qrcodeSize, err = strconv.Atoi(size)
		if err != nil {
			res := Response{Message: "Size in invalid."}
			b, _ := json.Marshal(res)
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, string(b))
			return
		}
	}

	png, err := qrcode.Encode(data, qrcode.Medium, qrcodeSize)
	if err != nil {
		res := Response{Message: "Failed to create the QR code from the given data."}
		b, _ := json.Marshal(res)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, string(b))
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

	port := 8000
	log.Printf("Listening on port :%d", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), handlers.CompressHandler(r)))
}
