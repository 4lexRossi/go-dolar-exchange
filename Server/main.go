package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

const (
	timeoutAPICall   = 200 * time.Millisecond
	timeoutDBPersist = 10 * time.Millisecond
	sqliteDBPath     = "./cotacao.db"
)

type ExchangeRate struct {
	Bid string `json:"bid"`
}

func persistExchangeRate(ctx context.Context, bid float64) error {
	db, err := sqlx.Open("sqlite3", sqliteDBPath)
	if err != nil {
		return err
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(ctx, timeoutDBPersist)
	defer cancel()

	query := "INSERT INTO cotacoes (valor, timestamp) VALUES ($1, $2)"
	_, err = db.ExecContext(ctx, query, bid, time.Now().Unix())
	return err
}

func cotacaoHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ctx, cancel := context.WithTimeout(ctx, timeoutAPICall)
	defer cancel()

	// Consumir API externa para obter a cotação do dólar
	resp, err := http.Get("https://economia.awesomeapi.com.br/json/last/USD-BRL")
	if err != nil {
		http.Error(w, fmt.Sprintf("Erro ao chamar API de cotação: %v", err), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		http.Error(w, fmt.Sprintf("Erro na resposta da API: %s", resp.Status), http.StatusInternalServerError)
		return
	}

	var exchangeRate map[string]ExchangeRate
	if err := json.NewDecoder(resp.Body).Decode(&exchangeRate); err != nil {
		http.Error(w, fmt.Sprintf("Erro ao decodificar resposta JSON: %v", err), http.StatusInternalServerError)
		return
	}

	if rate, ok := exchangeRate["USDBRL"]; ok {
		// Parse the bid string to float64
		bid, err := strconv.ParseFloat(rate.Bid, 64)
		if err != nil {
			http.Error(w, fmt.Sprintf("Erro ao converter bid para float64: %v", err), http.StatusInternalServerError)
			return
		}

		// Persistir no banco de dados
		if err := persistExchangeRate(ctx, bid); err != nil {
			log.Printf("Erro ao persistir cotação no banco de dados: %v", err)
			http.Error(w, "Erro ao persistir cotação no banco de dados", http.StatusInternalServerError)
			return
		}

		// Retornar apenas o valor do bid para o cliente
		response := map[string]float64{"bid": bid}
		json.NewEncoder(w).Encode(response)

	} else {
		http.Error(w, "Não foi possível obter a cotação do dólar", http.StatusInternalServerError)
		return
	}
}

func main() {
	// Criar banco de dados SQLite se não existir
	createDBIfNotExist()

	// Criar o servidor HTTP
	mux := http.NewServeMux()
	mux.HandleFunc("/cotacao", cotacaoHandler)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	log.Println("Servidor iniciado na porta :8080")
	log.Fatal(srv.ListenAndServe())
}

func createDBIfNotExist() {
	db, err := sqlx.Open("sqlite3", sqliteDBPath)
	if err != nil {
		log.Fatalf("Erro ao abrir banco de dados: %v", err)
	}
	defer db.Close()

	// Criar a tabela cotacoes se não existir
	query := `
			CREATE TABLE IF NOT EXISTS cotacoes (
					id INTEGER PRIMARY KEY AUTOINCREMENT,
					valor REAL NOT NULL,
					timestamp INTEGER NOT NULL
			)
	`
	if _, err := db.Exec(query); err != nil {
		log.Fatalf("Erro ao criar tabela no banco de dados: %v", err)
	}
}
