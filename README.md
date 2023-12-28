# Lambda DX
A simple wrapper that lets you run AWS Lambda functions locally through local HTTP endpoints.

## Example Usage
```go
package main

import (
	"github.com/jamesyeap/lambda_dx"
	lambdahandler1 "lambda_function_1/handler"
	lambdahandler2 "lambda_function_2/handler"
)

func main() {
	server := lambda_dx.NewBoxHttpServer()

	/* add all the lambda functions for testing here */
	server.AddLambdaFunction([]string{"GET"}, "/route_to_lambda_function_1", lambdahandler1.HandleRequest)
	server.AddLambdaFunction([]string{"POST"}, "/route_to_lambda_function_2", lambdahandler2.HandleRequest)
	// ... ... ... ... ... ... ... ... ... ... ... ... ... ... ...
	// ... ... add more lambda functions for testing here ... ....
	// ... ... ... ... ... ... ... ... ... ... ... ... ... ... ...

	server.Start(8080) // start the server on localhost:8080
}
```
- [link to example](https://github.com/jamesyeap/lambda_dx_example)