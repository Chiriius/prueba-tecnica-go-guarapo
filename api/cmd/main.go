package main

import (
	"os"
	"prueba_tecnica_go_guarapo/api/server"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func main() {
	_ = godotenv.Load()

	logger := logrus.New()
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	srv := server.NewServer(logger)
	srv.Start(":" + port)
}
