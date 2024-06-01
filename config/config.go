package config

import (
    "bufio"
    "fmt"
    "os"
    "strings"
    "time"
)

type Config struct {
    Port            string
    DataFile        string
    WindowDuration  time.Duration
    PersistInterval time.Duration
}

// LoadEnvFile loads environment variables from a given file
func LoadEnvFile(filename string) error {
    file, err := os.Open(filename)
    if err != nil {
        return fmt.Errorf("could not open .env file: %w", err)
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        line := scanner.Text()

        // Ignore comments and empty lines
        if strings.HasPrefix(line, "#") || strings.TrimSpace(line) == "" {
            continue
        }

        // Split key and value
        parts := strings.SplitN(line, "=", 2)
        if len(parts) != 2 {
            return fmt.Errorf("invalid line in .env file: %s", line)
        }
        key := strings.TrimSpace(parts[0])
        value := strings.TrimSpace(parts[1])

        // Set the environment variable
        fmt.Printf(" set key %s to value %s\n", key, value)
        if err := os.Setenv(key, value); err != nil {
            return fmt.Errorf("could not set environment variable: %w", err)
        }
    }

    if err := scanner.Err(); err != nil {
        return fmt.Errorf("error reading .env file: %w", err)
    }

    return nil
}

func LoadConfig() (*Config, error) {
    // Load .env
    err := LoadEnvFile(".env")
    if err != nil {
        return nil, fmt.Errorf("error loading .env file: %w", err)
    }

    // Retrieve variables
    port := os.Getenv("PORT")
    dataFile := os.Getenv("DATA_FILE")
    windowDuration, _ := time.ParseDuration(os.Getenv("WINDOW_DURATION"))
    persistInterval, _ := time.ParseDuration(os.Getenv("PERSIST_INTERVAL"))

    cfg := &Config{
        Port:            port,
        DataFile:        dataFile,
        WindowDuration:  windowDuration,
        PersistInterval: persistInterval,
    }

    return cfg, nil
}
