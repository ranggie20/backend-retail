package whatsapp

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Client struct {
	url  string
	auth string
	hcl  *http.Client
}

func New(url, auth string, client *http.Client) *Client {
	return &Client{url, auth, client}
}

func (cl *Client) Send(ctx context.Context, phone, msg string) error {
	payload := struct {
		Phone   string
		Message string
	}{Phone: phone, Message: msg}

	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(payload)
	if err != nil {
		log.Fatal(err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, cl.url, &buf)
	if err != nil {
		return fmt.Errorf("error on creating new get request object: %w", err)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprint("Basic ", cl.auth))

	resp, err := cl.hcl.Do(req)
	if err != nil {
		return fmt.Errorf("error on sending request: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		errResp := WhatsaappSendResponse{}
		if err = json.NewDecoder(resp.Body).Decode(&errResp); err != nil {
			return err
		}
		fmt.Printf("%+v \n\n", errResp)
		return fmt.Errorf("try again later")
	}

	return nil
}

type WhatsaappSendResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Results struct {
		MessageID string `json:"message_id"`
		Status    string `json:"status"`
	} `json:"results"`
}
