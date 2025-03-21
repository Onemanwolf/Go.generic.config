### Documentation: Generic Configuration Implementation in Go

This code provides a generic implementation for loading configuration values into a Go struct. It supports reading from a .env file and falling back to environment variables. The implementation uses reflection to map environment variables to struct fields based on struct tags.

---

### Key Components

#### 1. **Config Struct**
The Config struct is an example configuration structure. Each field is annotated with a `env` tag that specifies the corresponding environment variable name.

```go
type Config struct {
    MongoDBHost     string `env:"MONGO_HOST"`
    MongoDBPort     int    `env:"MONGO_PORT"`
    MongoDBUser     string `env:"MONGO_USER"`
    MongoDBPassword string `env:"MONGO_PASSWORD"`
    DebugMode       bool   `env:"DEBUG_MODE"`
}
```

- **Field Tags**: The `env` tag specifies the environment variable name to map to the field.
- **Supported Types**: `string`, `int`, `bool`, and `float64`.

---

#### 2. **`InitializeConfig` Function**
This function initializes the configuration by:
1. Attempting to load a .env file.
2. Populating the provided struct with values from environment variables.

```go
func InitializeConfig[T any](config *T) error {
    err := loadEnvFile("../.env") // Replace with your actual path
    if err != nil {
        fmt.Printf("No .env file loaded (falling back to environment variables): %v\n", err)
    }
    return parseConfig(config)
}
```

- **Generic Parameter**: The function is generic (`T any`), allowing it to work with any struct type.
- **Error Handling**: If the .env file cannot be loaded, it falls back to system environment variables.

---

#### 3. **`loadEnvFile` Function**
This function reads a .env file and sets environment variables using `os.Setenv`.

```go
func loadEnvFile(filePath string) error {
    file, err := os.Open(filePath)
    if err != nil {
        return err
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        line := strings.TrimSpace(scanner.Text())
        if line == "" || strings.HasPrefix(line, "#") {
            continue
        }

        parts := strings.SplitN(line, "=", 2)
        if len(parts) != 2 {
            continue
        }

        key := strings.TrimSpace(parts[0])
        value := strings.TrimSpace(parts[1])

        if os.Getenv(key) == "" {
            if err := os.Setenv(key, value); err != nil {
                return fmt.Errorf("failed to set env var %s: %v", key, err)
            }
        }
    }

    return scanner.Err()
}
```

- **File Parsing**: Reads key-value pairs from the .env file.
- **Environment Variables**: Sets environment variables only if they are not already set.

---

#### 4. **`parseConfig` Function**
This function populates a struct of type `T` using environment variables based on the `env` tags.

```go
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
```

- **Reflection**: Uses the `reflect` package to iterate over struct fields and set their values.
- **Type Handling**: Supports `string`, `int`, `bool`, and `float64` types.
- **Error Handling**: Returns errors for unsupported types or invalid values.

---

### Example Usage

```go
package main

import (
    "fmt"
    "config"
)

func main() {
    var cfg config.Config
    err := config.InitializeConfig(&cfg)
    if err != nil {
        fmt.Printf("Failed to initialize config: %v\n", err)
        return
    }

    fmt.Printf("Loaded Config: %+v\n", cfg)
}
```

- **Initialization**: Pass a pointer to the Config struct to `InitializeConfig`.
- **Output**: The struct is populated with values from the .env file or environment variables.

---

### Advantages
1. **Generic Implementation**: Works with any struct type.
2. **Fallback Mechanism**: Uses .env file or system environment variables.
3. **Type Safety**: Validates and converts environment variable values to the appropriate types.

### Limitations
1. **Unsupported Types**: Only supports basic types (`string`, `int`, `bool`, `float64`).
2. **Error Handling**: Requires careful handling of malformed .env files or invalid environment variable values.

### Additional Example: Using `github.com/joho/godotenv` for Configuration Management in the `config` folder

The updated implementation uses the `github.com/joho/godotenv` package to load environment variables from a .env file. This simplifies the process of managing environment variables and eliminates the need for custom .env file parsing logic.

---

### Example Usage with `godotenv`

Below is an example of how to use the `InitializeConfig` function with the `godotenv` dependency to load configuration values.

#### .env File Example
Create a .env file in your project directory with the following content:

```
MONGO_HOST=localhost
MONGO_PORT=27017
MONGO_USER=admin
MONGO_PASSWORD=secret
DEBUG_MODE=true
```

#### Main Application Code
```go
package main

import (
    "fmt"
    "config"
)

func main() {
    // Define a configuration struct
    var cfg config.Config

    // Initialize the configuration
    err := config.InitializeConfig(&cfg)
    if err != nil {
        fmt.Printf("Failed to initialize config: %v\n", err)
        return
    }

    // Print the loaded configuration
    fmt.Printf("Loaded Config: %+v\n", cfg)
}
```

---

### Explanation of the Example

1. **`.env` File**:
   - The .env file contains key-value pairs for configuration.
   - These values are loaded into environment variables using the `godotenv.Load` function.

2. **`InitializeConfig` Function**:
   - Attempts to load the .env file using `godotenv.Load("../.env")`.
   - If the .env file is not found or cannot be loaded, it falls back to system environment variables.
   - Calls `parseConfig` to populate the Config struct with values from the environment.

3. **Output**:
   - If the .env file is successfully loaded, the Config struct is populated with the values from the file.
   - If the .env file is missing, the program uses system environment variables instead.

---

### Advantages of Using `godotenv`

1. **Simplified Parsing**:
   - The `godotenv` package handles parsing .env files, reducing the need for custom logic.

2. **Fallback Mechanism**:
   - If the .env file is not found, the program seamlessly falls back to system environment variables.

3. **Compatibility**:
   - Works well with existing tools and workflows that rely on .env files.

---

### Example Output

Assuming the .env file contains the values shown above, running the program will produce the following output:

```
Loaded Config: {MongoDBHost:localhost MongoDBPort:27017 MongoDBUser:admin MongoDBPassword:secret DebugMode:true}
```

If the .env file is missing, and the corresponding environment variables are not set, the fields in the Config struct will retain their zero values.

---

This example demonstrates how to integrate the `github.com/joho/godotenv` package into your configuration management workflow, making it easier to manage environment variables in Go applications.



### Documentation: Setting Environment Variables in Bash

The provided commands.sh script is used to set environment variables for a Go application. These variables can be accessed by the application during runtime to configure its behavior.

---

### Script Overview

The script uses the `export` command to define and set environment variables in the current shell session. These variables are typically used to configure application settings such as database connections, authentication credentials, and debug modes.

#### Example Script: commands.sh
```bash
## Set environment variables
export MONGO_HOST="mongodb.atlas.example.com"
export MONGO_PORT="27017"
export MONGO_USER="admin"
export MONGO_PASSWORD="supersecretpassword"
export DEBUG_MODE="true"
```

---

### Explanation of Each Variable

1. **`MONGO_HOST`**:
   - Specifies the hostname or IP address of the MongoDB server.
   - Example: `mongodb.atlas.example.com`.

2. **`MONGO_PORT`**:
   - Specifies the port number on which the MongoDB server is running.
   - Example: `27017`.

3. **`MONGO_USER`**:
   - Specifies the username for authenticating with the MongoDB server.
   - Example: `admin`.

4. **`MONGO_PASSWORD`**:
   - Specifies the password for authenticating with the MongoDB server.
   - Example: `supersecretpassword`.

5. **`DEBUG_MODE`**:
   - Enables or disables debug mode in the application.
   - Example: `true` (enabled) or `false` (disabled).

---

### How to Use the Script

1. **Run the Script**:
   Execute the script in your terminal to set the environment variables in the current shell session:
   ```bash
   source ./scripts/commands.sh
   ```

   - The `source` command ensures that the variables are set in the current shell session, not a subshell.

2. **Verify the Variables**:
   After running the script, you can verify that the variables are set using the `echo` command:
   ```bash
   echo $MONGO_HOST
   echo $MONGO_PORT
   echo $MONGO_USER
   echo $MONGO_PASSWORD
   echo $DEBUG_MODE
   ```

3. **Run the Application**:
   Once the environment variables are set, you can run your Go application, and it will use these variables for configuration.

---

### Example Workflow

1. Run the script to set the environment variables:
   ```bash
   source ./scripts/commands.sh
   ```

2. Verify the variables:
   ```bash
   echo $MONGO_HOST
   # Output: mongodb.atlas.example.com
   ```

3. Run the Go application:
   ```bash
   go run main.go
   ```

---

### Notes

- **Security**: Avoid hardcoding sensitive information (e.g., passwords) in scripts. Use secure methods like .env files or secret management tools in production.
- **Persistence**: Environment variables set using `export` are only available in the current shell session. To make them persistent, add the `export` commands to your shell's configuration file (e.g., `.bashrc` or `.zshrc`).

This script is a simple and effective way to configure your application using environment variables.
