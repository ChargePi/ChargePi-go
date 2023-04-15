package main

import (
	"github.com/xBlaz3kx/ChargePi-go/internal/api/http"
)

func main() {

	// Launch UI at http://localhost:4269/
	// The UI should be integrated for portability.
	ui := http.NewUi()
	ui.Serve("0.0.0.0:8080")
}
