package ctxrouter_test

import (
	"fmt"
	"github.com/corestoreio/csfw/net/ctxrouter"
	"golang.org/x/net/context"
	"log"
	"net/http"
)

func Index(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	fmt.Fprint(w, "Welcome!\n")
	return nil
}

func Hello(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	ps := ctxrouter.ParamsFromContext(ctx)
	fmt.Fprintf(w, "hello, %s!\n", ps.ByName("name"))
	return nil
}

func Example() {

	router := ctxrouter.New()
	router.GET("/", Index)
	router.GET("/hello/:name", Hello)

	log.Fatal(http.ListenAndServe(":8080", router))

}
