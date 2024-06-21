package handlers

import (
	"context"
	"encoding/json"
	"log"
	"math/big"
	"net/http"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
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

func CUSDBalanceHandler(w http.ResponseWriter, r *http.Request) {
    address := r.URL.Query().Get("address")
    if address == "" {
        http.Error(w, "Address not provided", http.StatusBadRequest)
        return
    }

    balance := getCUSDBalance(r.Context(), address)
    w.Header().Set("Connection", "close")
    json.NewEncoder(w).Encode(BalanceResponse{Balance: balance})
}


func TransferCUSDHandler(w http.ResponseWriter, r *http.Request) {
    var req TransferRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    // In a real application, you would perform the transfer here
    // For now, we'll just return a success message
    w.Header().Set("Connection", "close")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{"status": "Transaction initiated"})
}


func getCUSDBalance(ctx context.Context, address string) string {
    client, err := ethclient.DialContext(ctx, "https://alfajores-forno.celo-testnet.org")
    if err != nil {
        log.Printf("Failed to connect to the Celo network: %v", err)
        return "0"
    }
    defer client.Close()

    contractAddress := common.HexToAddress(cUSDAddress)
    parsedABI, err := abi.JSON(strings.NewReader(`[{"constant":true,"inputs":[{"name":"_owner","type":"address"}],"name":"balanceOf","outputs":[{"name":"","type":"uint256"}],"payable":false,"type":"function"}]`))
    if err != nil {
        log.Printf("Failed to parse ABI: %v", err)
        return "0"
    }

    callData, err := parsedABI.Pack("balanceOf", common.HexToAddress(address))
    if err != nil {
        log.Printf("Failed to pack call data: %v", err)
        return "0"
    }

    msg := ethereum.CallMsg{
        To:   &contractAddress,
        Data: callData,
    }

    result, err := client.CallContract(ctx, msg, nil)
    if err != nil {
        log.Printf("Failed to call contract: %v", err)
        return "0"
    }

    balance := new(big.Int)
    balance.SetBytes(result)
    balanceInDecimals := new(big.Float).Quo(new(big.Float).SetInt(balance), big.NewFloat(float64(1e18)))

    return balanceInDecimals.Text('f', 6)
}