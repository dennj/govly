package lib

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

// VATRequest represents the incoming VAT request payload
type VATRequest struct {
	Country                  string `json:"country" binding:"required"`
	Name                     string `json:"name" binding:"required"`
	RegNum                   string `json:"regnum" binding:"required"`
	StartDate                string `json:"startdate" binding:"required"`
	EndDate                  string `json:"enddate" binding:"required"`
	Sales                    string `json:"sales" binding:"required"`
	Purchases                string `json:"purchs" binding:"required"`
	GoodsTo                  string `json:"goodsto" binding:"required"`
	GoodsFrom                string `json:"goodsfrom" binding:"required"`
	ServicesTo               string `json:"servicesto" binding:"required"`
	ServicesFrom             string `json:"servicesfrom" binding:"required"`
	PostponedAccounting      string `json:"postponedAccounting"`
	UnusualExpenditure       string `json:"unusualExpenditure"`
	UnusualExpenditureAmount string `json:"unusualExpenditureAmt"`
	UnusualExpenditureDetail string `json:"unusualExpenditureDtl"`
	Token                    string `json:"token" binding:"required"` // Token field added
}

// SendHTTPRequest is a utility function for sending HTTP requests
func SendHTTPRequest(method, url, contentType, body, token string) (*Response, error) {
	req, err := http.NewRequest(method, url, bytes.NewBufferString(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	return &Response{
		Status: resp.Status,
		Body:   string(responseBody),
	}, nil
}

// Response represents the HTTP response structure
type Response struct {
	Status string
	Body   string
}
