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

type DepartmentACLClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewDepartmentACLClient(baseURL string) *DepartmentACLClient {
	return &DepartmentACLClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

type DepartmentACLResponse struct {
	Code    int           `json:"code"`
	Message string        `json:"message"`
	Data    DepartmentData `json:"data"`
}

type DepartmentData struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

func (c *DepartmentACLClient) GetDepartment(ctx context.Context, departmentID int64) (*DepartmentInfo, error) {
	url := fmt.Sprintf("%s/api/v1/departments/%d", c.baseURL, departmentID)

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
		return nil, errors.New("department not found")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var aclResp DepartmentACLResponse
	if err := json.Unmarshal(body, &aclResp); err != nil {
		return nil, err
	}

	if aclResp.Code != 0 {
		return nil, errors.New(aclResp.Message)
	}

	return &DepartmentInfo{
		ID:   aclResp.Data.ID,
		Name: aclResp.Data.Name,
		Type: aclResp.Data.Type,
	}, nil
}

func (c *DepartmentACLClient) ValidateDepartment(ctx context.Context, departmentID int64) error {
	_, err := c.GetDepartment(ctx, departmentID)
	return err
}