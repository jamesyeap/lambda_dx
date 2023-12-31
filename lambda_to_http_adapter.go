package lambda_dx

import (
	"encoding/base64"
	"github.com/aws/aws-lambda-go/events"
	"io"
	"net/http"
)

type LambdaHandlerFunction func(events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)
type HttpHandlerFunction func(w http.ResponseWriter, r *http.Request)

type LambdaToHttpAdapter struct {
	handler *LambdaHandlerFunction
}

func NewLambdaToHttpAdapter(lambdaHandler LambdaHandlerFunction) HttpHandlerFunction {
	wrappedHandler := LambdaToHttpAdapter{handler: &lambdaHandler}

	return wrappedHandler.ServeHttp
}

func (l *LambdaToHttpAdapter) ServeHttp(w http.ResponseWriter, r *http.Request) {
	lambdaRequest, err := convertHTTPRequestToAPIGatewayProxyRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError) // 500 Internal Server Error
		return
	}

	lambdaResponse, lambdaError := (*l.handler)(lambdaRequest)
	if lambdaError != nil {
		// if status code is unset, set it to 500
		statusCode := lambdaResponse.StatusCode
		if statusCode == 0 {
			statusCode = 500
		}

		// if error message is unset, set it to a default error message
		errorMessage := lambdaError.Error()
		if errorMessage == "" {
			errorMessage = "Operation failed due to an unspecified error."
		}

		http.Error(w, errorMessage, statusCode)
		return
	}

	err = writeAPIGatewayProxyResponseToHTTPResponse(w, lambdaResponse)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError) // 500 Internal Server Error
		return
	}
}

func convertHTTPRequestToAPIGatewayProxyRequest(req *http.Request) (events.APIGatewayProxyRequest, error) {
	// prepare query string parameters
	queryStringParameters := make(map[string]string)
	for name, values := range req.URL.Query() {
		if len(values) > 0 {
			queryStringParameters[name] = values[0]
		}
	}

	// prepare the body (if any)
	var body string
	if req.Body != nil {
		bytes, err := io.ReadAll(req.Body)
		if err != nil {
			return events.APIGatewayProxyRequest{}, err
		}
		body = string(bytes)
	}

	// prepare headers
	headers := make(map[string]string)
	for name, values := range req.Header {
		if len(values) > 0 {
			headers[name] = values[0]
		}
	}

	// create the APIGatewayProxyRequest
	lambdaRequest := events.APIGatewayProxyRequest{
		HTTPMethod:            req.Method,
		Headers:               headers,
		QueryStringParameters: queryStringParameters,
		Body:                  body,
		IsBase64Encoded:       true, // set to true since body is base64 encoded
	}

	return lambdaRequest, nil
}

func writeAPIGatewayProxyResponseToHTTPResponse(w http.ResponseWriter, lambdaResponse events.APIGatewayProxyResponse) error {
	// set the status code
	w.WriteHeader(lambdaResponse.StatusCode)

	// set the headers
	for key, value := range lambdaResponse.Headers {
		w.Header().Set(key, value)
	}

	// check if the body is base64 encoded
	if lambdaResponse.IsBase64Encoded {
		// decode the body if it's base64 encoded
		decodedBody, err := base64.StdEncoding.DecodeString(lambdaResponse.Body)
		if err != nil {
			// handle error (for example, you might want to write an internal server error response)
			http.Error(w, "Failed to decode base64 body", http.StatusInternalServerError)
			return err
		}
		_, err = w.Write(decodedBody)
		if err != nil {
			return err
		}
	} else {
		// write the body as is
		_, err := w.Write([]byte(lambdaResponse.Body))
		if err != nil {
			return err
		}
	}

	// return no errors if everything is successful
	return nil
}
