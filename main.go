package main

import (
	"skillmaster/cmd"
	"skillmaster/pkg/config"
)

func main() {
	// Initialize global config if it doesn't exist
	config.InitializeConfig()
	
	// Execute root command
	cmd.Execute()
}
