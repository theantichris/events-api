package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	args := Args{
		conn: os.Getenv("DB"),
		port: os.Getenv("PORT"),
	}

	if err := Run(args); err != nil {
		log.Println(err)
	}
}
