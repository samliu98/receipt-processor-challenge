package main

import (
	"ReceiptApi/pkg/server"
	"net/http"
)

func main() {
	r := server.SetupRouter()
	http.ListenAndServe(":8080", r)
}
