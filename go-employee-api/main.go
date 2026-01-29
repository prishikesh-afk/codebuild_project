
package main

import (
	"context"
	"encoding/base64"
	"io"
	"net/http"
	"net/url"
	"strings"

	"go-employee-api/routers"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// Lambda ResponseWriter
type lambdaResponseWriter struct {
	headers http.Header
	body    []byte
	status  int
}

func newLambdaResponseWriter() *lambdaResponseWriter {
	return &lambdaResponseWriter{
		headers: make(http.Header),
		status:  http.StatusOK,
	}
}

func (l *lambdaResponseWriter) Header() http.Header {
	return l.headers
}

func (l *lambdaResponseWriter) Write(b []byte) (int, error) {
	l.body = append(l.body, b...)
	return len(b), nil
}

func (l *lambdaResponseWriter) WriteHeader(statusCode int) {
	l.status = statusCode
}

// Build query string
func buildQuery(params map[string]string) string {
	values := url.Values{}
	for k, v := range params {
		values.Add(k, v)
	}
	return values.Encode()
}

// Lambda handler
func handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	defer func() {
		if r := recover(); r != nil {
		}
	}()

	// Construct full URL
	urlStr := "http://localhost" + req.Path
	if len(req.QueryStringParameters) > 0 {
		urlStr += "?" + buildQuery(req.QueryStringParameters)
	}

	// Create HTTP request
	r, err := http.NewRequest(req.HTTPMethod, urlStr, nil)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Internal server error",
		}, nil
	}

	// Copy headers
	for k, v := range req.Headers {
		r.Header.Set(k, v)
	}

	// Copy body safely
	if req.Body != "" {
		var bodyBytes []byte
		if req.IsBase64Encoded {
			bodyBytes, err = base64.StdEncoding.DecodeString(req.Body)
			if err != nil {
				return events.APIGatewayProxyResponse{
					StatusCode: 400,
					Body:       "Invalid Base64 body",
				}, nil
			}
		} else {
			bodyBytes = []byte(req.Body)
		}
		r.Body = io.NopCloser(strings.NewReader(string(bodyBytes)))
	} else {
		r.Body = http.NoBody
	}

	// Response writer
	w := newLambdaResponseWriter()

	// Call router
	routers.MainRouter(w, r)

	// Convert headers
	respHeaders := map[string]string{}
	for k, v := range w.headers {
		if len(v) > 0 {
			respHeaders[k] = v[0]
		}
	}

	// Ensure body is string
	bodyStr := string(w.body)
	if bodyStr == "" {
		bodyStr = ""
	}

	return events.APIGatewayProxyResponse{
		StatusCode: w.status,
		Headers:    respHeaders,
		Body:       bodyStr,
	}, nil
}

func main() {
	lambda.Start(handler)
}
