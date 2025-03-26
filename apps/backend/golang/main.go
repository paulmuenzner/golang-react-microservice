package main

import (
	"log"
	"os"
	"os/exec"
)

func main() {
	log.Println("Starting Microservices...")

	// Starte api_gateway
	go runCommand("modules/user_service", "air")
	go runCommand("modules/order_service", "air")
	go runCommand("api_gateway", "air")

	select {} // Halte die Hauptfunktion am Laufen
}

func runCommand(path, command string) {
	cmd := exec.Command(command)
	cmd.Dir = path
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		log.Printf("Error running %s in %s: %v", command, path, err)
	}
}
