package main

import (
	"assign/config"
	"assign/core"
)

func main() {
	config.LoadConfig()
	core.Assign()
}
