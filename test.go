package main

import (
	"fmt"
	"github.com/coocood/jas"
	"net/http"
)

type Url struct {
	//url string
}

func (x *Url) Get(ctx *jas.Context) { // `GET /v1/hello`
	ctx.Data = ctx.id //response: `{"data":"hello world","error":null}`
}

func main() {
	router := jas.NewRouter(new(Url))
	router.BasePath = "/v1/"
	fmt.Println(router.HandledPaths(true))
	//output: `GET /v1/hello`
	http.Handle(router.BasePath, router)
	http.ListenAndServe(":8080", nil)
}
