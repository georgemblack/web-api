package main

import (
	"log/slog"
	"os"

	"github.com/georgemblack/web-api/pkg/api"
)

func main() {
	err := api.Run()
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}
