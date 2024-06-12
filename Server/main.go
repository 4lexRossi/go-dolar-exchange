package main

import (
	"encoding/json"
	"io"
	"net/http"
)

type USDBRL struct {
	Usdbrl struct {
		Code       string `json:"code"`
		Codein     string `json:"codein"`
		Name       string `json:"name"`
		High       string `json:"high"`
		Low        string `json:"low"`
		VarBid     string `json:"varBid"`
		PctChange  string `json:"pctChange"`
		Bid        string `json:"bid"`
		Ask        string `json:"ask"`
		Timestamp  string `json:"timestamp"`
		CreateDate string `json:"create_date"`
	} `json:"USDBRL"`
}

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
	exchange, error := DolarExchange(currency)
	if error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(exchange.Usdbrl.Bid)
}

func DolarExchange(currency string) (*USDBRL, error) {
	resp, error := http.Get("https://economia.awesomeapi.com.br/json/last/" + currency)
	if error != nil {
		return nil, error
	}
	defer resp.Body.Close()

	body, error := io.ReadAll(resp.Body)
	if error != nil {
		return nil, error
	}
	var e USDBRL
	error = json.Unmarshal(body, &e)
	if error != nil {
		return nil, error
	}

	return &e, nil
}
