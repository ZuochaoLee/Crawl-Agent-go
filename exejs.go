package crawlagent

import (
	"github.com/mcuadros/go-candyjs"
)

func Eval(js string) {
	ctx := candyjs.NewContext()
	ctx.EvalString(js)
}
