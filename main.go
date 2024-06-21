package main

import (
	"context"
	"encoding/json"
	"html/template"
	"log"
	"math/big"
	"net/http"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

var (
	templates = template.Must(template.ParseGlob("templates/*.html"))
)

const cUSDAddress = "0x874069Fa1Eb16D44d622F2e0Ca25eeA172369bC1"
const cUSDDecimals = 18

type BalanceResponse struct {
	Balance string `json:"balance"`
}

type TransferRequest struct {
	To     string `json:"to"`
	Amount string `json:"amount"`
}

func main() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/wallet-status", walletStatusHandler)
	http.HandleFunc("/cusd-balance", cUSDBalanceHandler)
	http.HandleFunc("/transfer-cusd", transferCUSDHandler)

	log.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "index.html", nil)
}

func walletStatusHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "wallet_status.html", nil)
}

func cUSDBalanceHandler(w http.ResponseWriter, r *http.Request) {
	address := r.URL.Query().Get("address")
	if address == "" {
		http.Error(w, "Address not provided", http.StatusBadRequest)
		return
	}

	balance := getCUSDBalance(address)
	json.NewEncoder(w).Encode(BalanceResponse{Balance: balance})
}

func transferCUSDHandler(w http.ResponseWriter, r *http.Request) {
	var req TransferRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// In a real application, you would perform the transfer here
	// For now, we'll just return a success message
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "Transaction initiated"})
}

func getCUSDBalance(address string) string {
	client, err := ethclient.DialContext(context.Background(), "https://alfajores-forno.celo-testnet.org")
	if err != nil {
		log.Printf("Failed to connect to the Celo network: %v", err)
		return "0"
	}
	defer client.Close()

	contractAddress := common.HexToAddress(cUSDAddress)
	data, err := abi.JSON(strings.NewReader(`[{"constant":true,"inputs":[{"name":"_owner","type":"address"}],"name":"balanceOf","outputs":[{"name":"","type":"uint256"}],"payable":false,"type":"function"}]`))
	if err != nil {
		log.Printf("Failed to parse ABI: %v", err)
		return "0"
	}

	callData, err := data.Pack("balanceOf", common.HexToAddress(address))
	if err != nil {
		log.Printf("Failed to pack call data: %v", err)
		return "0"
	}

	msg := ethereum.CallMsg{
		To:   &contractAddress,
		Data: callData,
	}

	result, err := client.CallContract(context.Background(), msg, nil)
	if err != nil {
		log.Printf("Failed to call contract: %v", err)
		return "0"
	}

	balance := new(big.Int)
	balance.SetBytes(result)
	balanceInDecimals := new(big.Float).Quo(new(big.Float).SetInt(balance), big.NewFloat(float64(1e18)))

	return balanceInDecimals.Text('f', 6)
}