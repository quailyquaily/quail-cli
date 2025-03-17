package main

import (
	"context"

	"github.com/quailyquaily/quail-cli/cmd"
	"github.com/quailyquaily/quail-cli/cmd/common"
)

var (
	Version = "0.0.1"
)

func main() {
	ctx := context.Background()
	ctx = context.WithValue(ctx, common.CTX_VERSION{}, Version)
	cmd.ExecuteContext(ctx)
}
