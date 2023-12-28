package lambda_dx

import (
	"fmt"
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

func (s *BoxHttpServer) AddLambdaFunction(httpVerbs []string, path string, lambdaHandler LambdaHandlerFunction) {
	newLambdaSpec := newLambdaSpec(httpVerbs, path, lambdaHandler)

	s.lambdaSpecList = append(s.lambdaSpecList, newLambdaSpec)
}

func (s *BoxHttpServer) Start(port int) {
	// start a new router instance
	r := mux.NewRouter()

	// add all lambda routes to the router
	for _, lambdaSpec := range s.lambdaSpecList {
		r.HandleFunc(lambdaSpec.Path, lambdaSpec.Handler).Methods(lambdaSpec.HttpVerbs...)
	}

	listenAddr := fmt.Sprintf(":%d", port)

	// start the router
	log.Printf("Listening on %s\n", listenAddr)
	err := http.ListenAndServe(listenAddr, r)

	if err != nil {
		log.Fatal(err)
	}
}

type LambdaSpec struct {
	HttpVerbs []string
	Path      string
	Handler   HttpHandlerFunction
}

func newLambdaSpec(httpVerbs []string, path string, lambdaHandler LambdaHandlerFunction) *LambdaSpec {
	httpHandler := NewLambdaToHttpAdapter(lambdaHandler)

	return &LambdaSpec{HttpVerbs: httpVerbs, Path: path, Handler: httpHandler}
}
