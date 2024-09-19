package handler

import (
	"net/http"

	"github.com/go-chi/render"
)

// SuccessResponse represents a success response
type SuccessResponse struct {
	HTTPStatusCode int         `json:"status,omitempty"`  // http response status code
	Message        string      `json:"message,omitempty"` // user-level success message
	Data           interface{} `json:"data,omitempty"`    // optional data to include in the success response
}

func (s *SuccessResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, s.HTTPStatusCode)
	return nil
}

// SuccessOK is a success response for HTTP status code 200
var SuccessOK = &SuccessResponse{HTTPStatusCode: http.StatusOK, Message: "OK"}

// SuccessCreated is a success response for HTTP status code 201
var SuccessCreated = &SuccessResponse{HTTPStatusCode: http.StatusCreated, Message: "Created"}

// SuccessNoContent is a success response for HTTP status code 204
var SuccessNoContent = &SuccessResponse{HTTPStatusCode: http.StatusNoContent, Message: "No Content"}
