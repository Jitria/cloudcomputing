package main

import (
	"assign/config"
)

func main() {
	config.LoadConfig()
	Assign()
}
