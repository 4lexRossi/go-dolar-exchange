package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	client "github.com/4lexRossi/go-dolar-exchange/Client"
)

type RateExchange struct {
	ID  int `gorm:"primaryKey"`
	Bid string
	gorm.Model
}

func main() {
	http.HandleFunc("/exchange-rate", DolarExchangeHandler)
	http.ListenAndServe(":8080", nil)
}

func DolarExchangeHandler(w http.ResponseWriter, r *http.Request) {
	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&RateExchange{})

	ctx := r.Context()
	if r.URL.Path != "/exchange-rate" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	select {
	case <-time.After(200 * time.Millisecond):
		log.Println("Request successfully")
		exchange, error := client.DolarExchange("https://economia.awesomeapi.com.br/json/last/USD-BRL")
		if error != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		json.NewEncoder(w).Encode(exchange.Usdbrl.Bid)

		db.Create(&RateExchange{
			Bid: exchange.Usdbrl.Bid,
		})
	case <-ctx.Done():
		log.Println("Resquest Failed")
	}

}
