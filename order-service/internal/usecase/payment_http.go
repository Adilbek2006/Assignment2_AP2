package usecase

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)

type HTTPPaymentClient struct {
	Client *http.Client
	URL    string
}

func (pc *HTTPPaymentClient) CreatePayment(orderID string, amount int64) (string, error) {
	body, err := json.Marshal(map[string]interface{}{
		"order_id": orderID,
		"amount":   amount,
	})
	if err != nil {
		return "", err
	}

	resp, err := pc.Client.Post(pc.URL+"/payments", "application/json", bytes.NewBuffer(body))
	if err != nil {
		return "", errors.New("payment service unavailable")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", errors.New("payment failed")
	}

	var result struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	return result.Status, nil
}
