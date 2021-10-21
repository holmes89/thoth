package handlers

import (
	"context"
	"net/http"
	"net/url"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/gorilla/mux"
)

type customWriter struct {
	events.APIGatewayProxyResponse
}

func (w *customWriter) Header() http.Header {
	if w.MultiValueHeaders == nil {
		w.MultiValueHeaders = make(map[string][]string)
	}
	return w.MultiValueHeaders
}

func (w *customWriter) WriteHeader(statusCode int) {
	w.StatusCode = statusCode
}

func (w *customWriter) Write(content []byte) (int, error) {
	w.Body = string(append([]byte(w.Body), content...))
	return len(content), nil
}

type APIRouter struct {
	*mux.Router
}

func NewAPIRouter() *APIRouter {
	return &APIRouter{
		mux.NewRouter(),
	}
}

func (r *APIRouter) HandleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	u := url.URL{
		Path: request.Path,
	}
	req, err := http.NewRequest(request.HTTPMethod, u.String(), strings.NewReader(request.Body))
	if err != nil {
		panic(err)
	}

	writer := &customWriter{}
	r.ServeHTTP(writer, req)

	return writer.APIGatewayProxyResponse, nil
}
