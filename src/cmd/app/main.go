package main

import "os"

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type 'Bearer ' followed by your token
func main() {
	app := NewApp()

	if err := app.Run(); err != nil {
		app.Logger().Error("Server stopped with error", "error", err)
		os.Exit(1)
	}
}
