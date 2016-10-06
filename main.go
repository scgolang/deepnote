package main

import (
	"log"
	"math/rand"
	"time"
)

func main() {
	// Seed the random number generator.
	rand.Seed(time.Now().Unix())

	// Get the configuration from the CLI.
	config, err := NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	// Create a new deepnote app.
	app, err := NewApp(config)
	if err != nil {
		log.Fatal(err)
	}

	// Run the app.
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}

// rrand generates a float in the interval [min, max)
func rrand(min, max float64) float64 {
	return (rand.Float64() * (max - min)) + min
}
