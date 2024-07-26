package main

import (
	"flag"
	"log"
)

// special predefined function
func init() {
	//setting prefix for all logs
	log.SetPrefix("Blockchain: ")
}

func main() {
	port := flag.Uint("port", 5000, "TCP Port Number for Blockchain Server")
	flag.Parse()
	app := NewBlockchainServer(uint16(*port))
	app.Run()
}
