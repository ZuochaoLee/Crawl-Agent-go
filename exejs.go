package main

import (
	"github.com/mcuadros/go-candyjs"
)

// func main() {
// 	ctx := candyjs.NewContext()
// 	a := ctx.EvalString("(1+2)")
// }

func Eval(js string) {
	ctx := candyjs.NewContext()
	ctx.EvalString(js)
}
