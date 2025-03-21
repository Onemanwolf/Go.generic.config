package main

import (
    "fmt"
    "github.com/Onemanwolf/go.generic.config/config" // Replace with your module path
)

func main() {
    // Define and initialize a Config instance
    cfg := &config.Config{}

    // Explicitly call InitializeConfig to populate the config
    err := config.InitializeConfig(cfg)
    if err != nil {
        fmt.Printf("Error initializing config: %v\n", err)
        return
    }

    // Print the populated config
    fmt.Printf("MongoDB Host: %s\n", cfg.MongoDBHost)
    fmt.Printf("MongoDB Port: %d\n", cfg.MongoDBPort)
    fmt.Printf("MongoDB User: %s\n", cfg.MongoDBUser)
    fmt.Printf("MongoDB Password: %s\n", cfg.MongoDBPassword)
    fmt.Printf("Debug Mode: %t\n", cfg.DebugMode)

}