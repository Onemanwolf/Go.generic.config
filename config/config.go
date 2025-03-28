package config

import (
    "fmt"
    "os"
    "reflect"
    "strconv"

    "github.com/joho/godotenv"
)

// Config represents the configuration struct.
type Config struct {
    MongoDBHost     string `env:"MONGO_HOST"`
    MongoDBPort     int    `env:"MONGO_PORT"`
    MongoDBUser     string `env:"MONGO_USER"`
    MongoDBPassword string `env:"MONGO_PASSWORD"`
    DebugMode       bool   `env:"DEBUG_MODE"`
}

// InitializeConfig initializes the configuration by first attempting to load a .env file,
// and falling back to system environment variables if the file is not available.
func InitializeConfig[T any](config *T) error {
    // Attempt to load .env file from a default location (optional)
    err := godotenv.Load("../.env")
    if err != nil {
        // If .env file is not found or cannot be loaded, proceed with system env vars
        fmt.Printf("Warning: No .env file loaded, falling back to system environment variables: %v\n", err)
    }

    // Parse config from environment variables (system vars override .env if both exist)
    return parseConfig(config)
}

// parseConfig is a generic helper function that populates a struct of type T
// using environment variables based on struct field tags.
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
            continue
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

