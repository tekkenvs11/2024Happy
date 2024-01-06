package main

import (
	_ "go-medium/internal/packed"

	_ "go-medium/internal/logic"

	"github.com/gogf/gf/v2/os/gctx"

	"go-medium/internal/cmd"
)

func main() {
	cmd.Main.Run(gctx.GetInitCtx())
}
