package lambda_dx

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type BoxHttpServer struct {
	lambdaSpecList []*LambdaSpec
}

func NewBoxHttpServer() *BoxHttpServer {
	return &BoxHttpServer{}
}

func (s *BoxHttpServer) AddLambdaFunction(path string, lambdaHandler LambdaHandlerFunction) {
	newLambdaSpec := newLambdaSpec(path, lambdaHandler)

	s.lambdaSpecList = append(s.lambdaSpecList, newLambdaSpec)
}

func (s *BoxHttpServer) Start() {
	// start a new router instance
	r := mux.NewRouter()

	// add all lambda routes to the router
	for _, lambdaSpec := range s.lambdaSpecList {
		r.HandleFunc(lambdaSpec.Path, lambdaSpec.Handler)
	}

	// start the router
	log.Print("Listening...")
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatal(err)
	}
}

type LambdaSpec struct {
	Path    string
	Handler HttpHandlerFunction
}

func newLambdaSpec(path string, lambdaHandler LambdaHandlerFunction) *LambdaSpec {
	httpHandler := NewLambdaToHttpAdapter(lambdaHandler)

	return &LambdaSpec{Path: path, Handler: httpHandler}
}

/*
	r := mux.NewRouter()

	// put all handlers here
	r.HandleFunc("/query", common.NewLambdaToHttpAdapter(handler.QueryHandler))

	// start the server
	log.Print("Listening...")
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Hello")
	r.HandleFunc("/query", common.NewLambdaToHttpAdapter(handler.QueryHandler))

*/
