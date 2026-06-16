package api

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

// HelloInput is the input for the hello endpoint.
type HelloInput struct{}

// HelloOutput is the response for the hello endpoint.
type HelloOutput struct {
	Body struct {
		Message string `json:"message" example:"Hello from the backend!"`
	}
}

func registerHello(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID: "get-hello",
		Method:      http.MethodGet,
		Path:        "/api/hello",
		Summary:     "Hello world",
		Description: "Returns a hello-world message confirming the backend is reachable.",
	}, func(ctx context.Context, input *HelloInput) (*HelloOutput, error) {
		resp := &HelloOutput{}
		resp.Body.Message = "Hello from the backend!"
		return resp, nil
	})
}
