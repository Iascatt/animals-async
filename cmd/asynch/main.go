package main

import (
	"lab8/internal/api"
	"log"
)

func main() {
	log.Println("app start")
	api.StartServer()
	log.Println("app terminated")
}
