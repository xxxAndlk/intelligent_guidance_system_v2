package acl

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

type PaymentACLClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewPaymentACLClient(baseURL string) *PaymentACLClient {
	return &PaymentACLClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

type PaymentACLResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    PaymentData `json:"data"`
}

type PaymentData struct {
	RecordID int64   `json:"record_id"`
	Amount   float64 `json:"amount"`
	Method   string  `json:"method"`
	Status   string  `json:"status"`
	PayTime  string  `json:"pay_time"`
}

type ProcessPaymentRequest struct {
	RecordID int64   `json:"record_id"`
	Amount   float64 `json:"amount"`
	Method   string  `json:"method"`
}

func (c *PaymentACLClient) ProcessPayment(ctx context.Context, recordID int64, amount float64, method string) error {
	url := fmt.Sprintf("%s/api/v1/payments/process", c.baseURL)

	reqBody := ProcessPaymentRequest{
		RecordID: recordID,
		Amount:   amount,
		Method:   method,
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(bodyBytes))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var aclResp PaymentACLResponse
	if err := json.Unmarshal(respBody, &aclResp); err != nil {
		return err
	}

	if aclResp.Code != 0 {
		return errors.New(aclResp.Message)
	}

	return nil
}

func (c *PaymentACLClient) GetPaymentStatus(ctx context.Context, recordID int64) (*PaymentStatusInfo, error) {
	url := fmt.Sprintf("%s/api/v1/payments/record/%d", c.baseURL, recordID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, errors.New("payment not found")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var aclResp PaymentACLResponse
	if err := json.Unmarshal(body, &aclResp); err != nil {
		return nil, err
	}

	if aclResp.Code != 0 {
		return nil, errors.New(aclResp.Message)
	}

	return &PaymentStatusInfo{
		RecordID: aclResp.Data.RecordID,
		Amount:   aclResp.Data.Amount,
		Method:   aclResp.Data.Method,
		Status:   aclResp.Data.Status,
		PayTime:  aclResp.Data.PayTime,
	}, nil
}

