package openai

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type BatchMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type BatchRequestBody struct {
	Model     string    `json:"model"`
	Messages  []Message `json:"messages"`
	MaxTokens int       `json:"max_tokens"`
}

type BatchRequest struct {
	CustomID string           `json:"custom_id"`
	Method   string           `json:"method"`
	URL      string           `json:"url"`
	Body     BatchRequestBody `json:"body"`
}

type CreateBatchRequestFileRequest struct {
	InputFileID      string `json:"input_file_id"`
	EndPoint         string `json:"end_point"`
	CompletionWindow int    `json:"completion_window"`
}

type RequestCounts struct {
	Total     int `json:"total"`
	Completed int `json:"completed"`
	Failed    int `json:"failed"`
}

type Metadata struct {
	CustomerID       string `json:"customer_id"`
	BatchDescription string `json:"batch_description"`
}

type Batch struct {
	ID               string        `json:"id"`
	Object           string        `json:"object"`
	Endpoint         string        `json:"endpoint"`
	Errors           interface{}   `json:"errors"` // Use interface{} to handle null or different types
	InputFileID      string        `json:"input_file_id"`
	CompletionWindow string        `json:"completion_window"`
	Status           string        `json:"status"`
	OutputFileID     string        `json:"output_file_id"`
	ErrorFileID      string        `json:"error_file_id"`
	CreatedAt        int64         `json:"created_at"`
	InProgressAt     int64         `json:"in_progress_at"`
	ExpiresAt        int64         `json:"expires_at"`
	FinalizingAt     int64         `json:"finalizing_at"`
	CompletedAt      int64         `json:"completed_at"`
	FailedAt         *int64        `json:"failed_at"`     // Use pointers to handle null values
	ExpiredAt        *int64        `json:"expired_at"`    // Use pointers to handle null values
	CancellingAt     *int64        `json:"cancelling_at"` // Use pointers to handle null values
	CancelledAt      *int64        `json:"cancelled_at"`  // Use pointers to handle null values
	RequestCounts    RequestCounts `json:"request_counts"`
	Metadata         Metadata      `json:"metadata"`

	httpHeader
}

// CreateBatchRequestFile .
func (c *Client) CreateBatchRequestFile(
	ctx context.Context,
	request []BatchRequest) (response File, err error) {

	// create a json file with the request
	fileName := fmt.Sprintf("batch_request_%s.json", request[0].CustomID)
	file, err := os.Create(fileName)
	if err != nil {
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(request)
	if err != nil {
		return
	}

	// create a file request
	fileRequest := FileRequest{
		Purpose:  "batch",
		FilePath: fileName,
	}

	// upload the file
	response, err = c.CreateFile(ctx, fileRequest)
	if err != nil {
		return
	}

	return
}

// CreateBatch creates a new batch.
func (c *Client) CreateBatch(ctx context.Context, request CreateBatchRequestFileRequest) (batch Batch, err error) {
	req, err := c.newRequest(ctx, http.MethodPost, c.fullURL("/v1/batches"), withBody(request),
		withBetaAssistantVersion(c.config.AssistantVersion))
	if err != nil {
		return
	}

	err = c.sendRequest(req, &batch)

	return
}
