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

type DoctorACLClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewDoctorACLClient(baseURL string) *DoctorACLClient {
	return &DoctorACLClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

type DoctorACLResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    DoctorData  `json:"data"`
}

type DoctorData struct {
	ID           int64  `json:"id"`
	Name         string `json:"name"`
	Title        string `json:"title"`
	DepartmentID int64  `json:"department_id"`
}

func (c *DoctorACLClient) GetDoctor(ctx context.Context, doctorID int64) (*DoctorInfo, error) {
	url := fmt.Sprintf("%s/api/v1/doctors/%d", c.baseURL, doctorID)

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
		return nil, errors.New("doctor not found")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var aclResp DoctorACLResponse
	if err := json.Unmarshal(body, &aclResp); err != nil {
		return nil, err
	}

	if aclResp.Code != 0 {
		return nil, errors.New(aclResp.Message)
	}

	return &DoctorInfo{
		ID:           aclResp.Data.ID,
		Name:         aclResp.Data.Name,
		Title:        aclResp.Data.Title,
		DepartmentID: aclResp.Data.DepartmentID,
	}, nil
}

func (c *DoctorACLClient) ValidateDoctor(ctx context.Context, doctorID int64) error {
	_, err := c.GetDoctor(ctx, doctorID)
	return err
}