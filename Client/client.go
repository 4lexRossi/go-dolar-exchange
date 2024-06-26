package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	serverURL  = "http://localhost:8080/cotacao"
	timeout    = 300 * time.Millisecond
	outputFile = "cotacao.txt"
)

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// Fazer requisição para o servidor
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, serverURL, nil)
	if err != nil {
		log.Fatalf("Erro ao criar requisição HTTP: %v", err)
	}

	// Enviar a requisição e obter a resposta
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Erro ao realizar requisição HTTP: %v", err)
	}
	defer resp.Body.Close()

	// Ler a resposta
	var exchangeRate map[string]float64
	if err := json.NewDecoder(resp.Body).Decode(&exchangeRate); err != nil {
		log.Fatalf("Erro ao decodificar resposta JSON: %v", err)
	}

	// Salvar a cotação em um arquivo
	bid, ok := exchangeRate["bid"]
	if !ok {
		log.Fatal("Campo 'bid' não encontrado na resposta")
	}

	content := fmt.Sprintf("Dólar: %.2f", bid)

	file, err := os.Create(outputFile)
	if err != nil {
		log.Fatalf("Erro ao criar arquivo: %v", err)
	}
	defer file.Close()

	if _, err := file.WriteString(content); err != nil {
		log.Fatalf("Erro ao escrever no arquivo: %v", err)
	}

	log.Printf("Cotação do dólar salva com sucesso no arquivo %s", outputFile)
}
