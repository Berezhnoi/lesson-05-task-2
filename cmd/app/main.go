package main

import (
	"fmt"
	"os"

	"error-handling-lab/internal/audit"
	"error-handling-lab/internal/config"
	"error-handling-lab/internal/handler"
	"error-handling-lab/internal/repository"
	"error-handling-lab/internal/service"
)

func main() {
	if err := config.Load("config.json"); err != nil {
		exitWithError(err)
	}

	userRepository := repository.New()
	userService := service.New(userRepository)
	appHandler := handler.New(userService, audit.New("audit.log"), os.Stdout)
	if err := appHandler.Handle(os.Args[1:]); err != nil {
		exitWithError(err)
	}
}

func exitWithError(err error) {
	if _, writeErr := fmt.Fprintln(os.Stderr, "error:", err); writeErr != nil {
		os.Exit(1)
	}
	os.Exit(1)
}
