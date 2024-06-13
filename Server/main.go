package main

import (
	"encoding/json"
	"net/http"

	client "github.com/4lexRossi/go-dolar-exchange/Client"
)

func main() {
	http.HandleFunc("/exchange-rate", DolarExchangeHandler)
	http.ListenAndServe(":8080", nil)
}

func DolarExchangeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/exchange-rate" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	exchange, error := client.DolarExchange("https://economia.awesomeapi.com.br/json/last/USD-BRL")
	if error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(exchange.Usdbrl.Bid)
}
