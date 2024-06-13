package main

import (
	"encoding/json"
	"net/http"

	client "github.com/4lexRossi/go-dolar-exchange/Client"
)

func main() {
	http.HandleFunc("/", DolarExchangeHandler)
	http.ListenAndServe(":8080", nil)
}

func DolarExchangeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	currency := r.URL.Query().Get("currency")
	if currency == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	exchange, error := client.DolarExchange(currency)
	if error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(exchange.Usdbrl.Bid)
}
