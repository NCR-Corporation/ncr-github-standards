package main

import (
	"log"
	"time"
)

func main() {
	log.Println("Starting Github-Standard routine")
	start := time.Now()

	// Gather QLIDs, send emails
	GetUsersAndFilter()
	elapsed := time.Since(start)
	log.Printf("\nProgram execution finished in %s\n", elapsed)
}
