package runtime

import (
	"github.com/srmagura/luma/shared"
)

func Execute(ast shared.Node) {
	shared.PrintAST(ast)

}
