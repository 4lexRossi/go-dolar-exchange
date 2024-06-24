package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
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

func DolarExchange(ctxReceived context.Context, url string) (*USDBRL, error) {
	req, err := http.NewRequestWithContext(ctxReceived, http.MethodGet, url, nil)
	if err != nil {
		log.Println("Request take too long")
		panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(300*time.Millisecond))
	defer cancel()

	req = req.WithContext(ctx)

	resp, error := http.DefaultClient.Do(req)
	if error != nil {
		log.Println("Response take too long")
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

	file, err := os.Create("exchange-rate.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating file: %v\n", err)
	}
	defer file.Close()
	file.WriteString(fmt.Sprintf("Dólar: %s", e.Usdbrl.Bid))
	return &e, nil
}
