package ctxrouter_test

import (
	"fmt"
	"log"
	"net/http"

	"github.com/corestoreio/csfw/net/ctxrouter"
)

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Welcome!\n")
}

func Hello(w http.ResponseWriter, r *http.Request) {
	ps := ctxrouter.FromContextParams(r.Context())
	fmt.Fprintf(w, "hello, %s!\n", ps.ByName("name"))
}

func Example() {

	router := ctxrouter.New()
	router.GET("/", Index)
	router.GET("/hello/:name", Hello)

	log.Fatal(http.ListenAndServe(":8080", router))

}
