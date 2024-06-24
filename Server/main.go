package main

import (
	"context"
	"encoding/json"
	"fmt"
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
	db, err := gorm.Open(sqlite.Open("exchange.db"), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to the database: %v", err)
	}
	// db.AutoMigrate(&RateExchange{})

	if r.URL.Path != "/exchange-rate" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	defer log.Println("Resquest Finish")

	ctx, cancel := context.WithTimeout(context.Background(), (1 * time.Nanosecond))
	defer cancel()

	exchange, err := client.DolarExchange(ctx, "https://economia.awesomeapi.com.br/json/last/USD-BRL")
	log.Println("Resquest Started")
	if err != nil {
		panic(err)
	}

	ctxDB, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	data := exchange.Usdbrl.Bid

	err = saveData(ctxDB, db, data)
	if err != nil {
		log.Printf("failed to save data: %v", err)
	} else {
		log.Println("data saved successfully")
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(exchange.Usdbrl.Bid)

}

func saveData(ctx context.Context, db *gorm.DB, data string) error {
	record := RateExchange{
		Bid: data,
	}

	if err := db.WithContext(ctx).Create(&record).Error; err != nil {
		return fmt.Errorf("failed to save data: %w", err)
	}

	return nil
}
