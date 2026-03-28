package acl

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

type PatientACLClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewPatientACLClient(baseURL string) *PatientACLClient {
	return &PatientACLClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

type PatientACLResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    PatientData `json:"data"`
}

type PatientData struct {
	ID     int64  `json:"id"`
	Name   string `json:"name"`
	Phone  string `json:"phone"`
	Gender int    `json:"gender"`
	Age    int    `json:"age"`
}

func (c *PatientACLClient) GetPatient(ctx context.Context, patientID int64) (*PatientInfo, error) {
	url := fmt.Sprintf("%s/api/v1/patients/%d", c.baseURL, patientID)

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
		return nil, errors.New("patient not found")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var aclResp PatientACLResponse
	if err := json.Unmarshal(body, &aclResp); err != nil {
		return nil, err
	}

	if aclResp.Code != 0 {
		return nil, errors.New(aclResp.Message)
	}

	return &PatientInfo{
		ID:     aclResp.Data.ID,
		Name:   aclResp.Data.Name,
		Phone:  aclResp.Data.Phone,
		Gender: aclResp.Data.Gender,
		Age:    aclResp.Data.Age,
	}, nil
}

func (c *PatientACLClient) ValidatePatient(ctx context.Context, patientID int64) error {
	_, err := c.GetPatient(ctx, patientID)
	return err
}