package config

import (
    "bufio"
    "fmt"
    "os"
    "reflect"
    "strconv"
    "strings"
)

// Config represents the configuration struct.
type Config struct {
    MongoDBHost     string `env:"MONGO_HOST"`
    MongoDBPort     int    `env:"MONGO_PORT"`
    MongoDBUser     string `env:"MONGO_USER"`
    MongoDBPassword string `env:"MONGO_PASSWORD"`
    DebugMode       bool   `env:"DEBUG_MODE"`
}

// InitializeConfig initializes the configuration by parsing a .env file or falling back to environment variables.
func InitializeConfig[T any](config *T) error {
    // Attempt to load .env file, but continue if it fails
    err := loadEnvFile("../.env") // Replace with your actual path
    if err != nil {
        fmt.Printf("No .env file loaded (falling back to environment variables): %v\n", err)
    }

    // Populate config from environment variables (set by .env or system)
    return parseConfig(config)
}

// loadEnvFile reads a .env file and sets environment variables using os.Setenv.
func loadEnvFile(filePath string) error {
    file, err := os.Open(filePath)
    if err != nil {
        return err // Return error if file can't be opened (e.g., doesn't exist)
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        line := strings.TrimSpace(scanner.Text())
        if line == "" || strings.HasPrefix(line, "#") {
            continue // Skip empty lines and comments
        }

        parts := strings.SplitN(line, "=", 2)
        if len(parts) != 2 {
            continue // Skip malformed lines
        }

        key := strings.TrimSpace(parts[0])
        value := strings.TrimSpace(parts[1])

        // Set the environment variable if not already set
        if os.Getenv(key) == "" {
            if err := os.Setenv(key, value); err != nil {
                return fmt.Errorf("failed to set env var %s: %v", key, err)
            }
        }
    }

    return scanner.Err()
}

// parseConfig populates a struct of type T using environment variables based on struct field tags.
func parseConfig[T any](config *T) error {
    val := reflect.ValueOf(config)
    if val.Kind() != reflect.Ptr || val.Elem().Kind() != reflect.Struct {
        return fmt.Errorf("config must be a pointer to a struct")
    }

    structVal := val.Elem()
    structType := structVal.Type()

    for i := 0; i < structType.NumField(); i++ {
        field := structType.Field(i)
        fieldVal := structVal.Field(i)

        envKey, ok := field.Tag.Lookup("env")
        if !ok {
            continue
        }

        envVal := os.Getenv(envKey)
        if envVal == "" {
            continue // Skip if no value is set (leave as zero value)
        }

        if !fieldVal.CanSet() {
            return fmt.Errorf("cannot set field %s", field.Name)
        }

        switch fieldVal.Kind() {
        case reflect.String:
            fieldVal.SetString(envVal)
        case reflect.Int, reflect.Int64:
            intVal, err := strconv.Atoi(envVal)
            if err != nil {
                return fmt.Errorf("invalid integer value for %s: %v", envKey, err)
            }
            fieldVal.SetInt(int64(intVal))
        case reflect.Bool:
            boolVal, err := strconv.ParseBool(envVal)
            if err != nil {
                return fmt.Errorf("invalid boolean value for %s: %v", envKey, err)
            }
            fieldVal.SetBool(boolVal)
        case reflect.Float64:
            floatVal, err := strconv.ParseFloat(envVal, 64)
            if err != nil {
                return fmt.Errorf("invalid float value for %s: %v", envKey, err)
            }
            fieldVal.SetFloat(floatVal)
        default:
            return fmt.Errorf("unsupported field type for %s: %v", field.Name, fieldVal.Kind())
        }
    }

    return nil
}

// Example usage
