package main

import (
	"flag"
	"fmt"
	"log"
	"time"
)

// NOTE: becareful of global variables. they should be read only. or they should be singletons
const (
	PATH = "/sys/class/power_supply/BAT0/" // Battery information path
)

// I would make a type alias of these. It can be easier to find and seperate concerns
type threshold = int

const (
	CRITICAL threshold = iota // Levels for the battery
	VERY_LOW
	LOW
	SAFE
)

var THRESHOLD = map[threshold]int{ // Battery Percentage thresholds that define the levels
	CRITICAL: 5,
	VERY_LOW: 10,
	LOW:      20,
}

func main() {
	intervalString := flag.String("interval", "1m", "")
	flag.Parse()

	interval, err := time.ParseDuration(*intervalString)
	if err != nil {
		log.Fatalf("parse interval: %w", err)
	}
	// Declare our battery struct
	b := battery{
		Level: SAFE,
	}

	// Main loop
	for {

		// Get battery status from file
		isCharing, charge, err := b.getValues()
		if err != nil {
			n := notification{
				"-t":   "ERROR",
				"-c":   err.Error(),
				"--fs": "12",
			}
			e := n.notify()
			if e != nil {
				log.Fatal(err)
			}
		}

		// updateLevel checks to see if the battery level drops below or rises above a defined threshold
		if b.updateLevel(charge) == -1 && !isCharing { // Battery level dropped below a threshold
			n := notification{
				"-t":    "LOW BATTERY",
				"-c":    fmt.Sprintf("Your battery is currently at %d%%. Please provide power.", charge),
				"-s":    "50",
				"-d":    "5000",
				"--fs":  "30",
				"--fg":  "red",
				"--pos": "bottom_right",
			}
			err := n.notify()
			if err != nil {
				log.Fatal(err)
			}
		}

		time.Sleep(interval)
	}
}
