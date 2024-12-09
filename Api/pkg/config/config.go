package config

import (
	"fmt"
	"log"

	"github.com/go-playground/validator"
	"github.com/knadh/koanf/parsers/json"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

// APIConfig represents the API-related configuration.
type ApiConfig struct {
	Port          int      `koanf:"port" validate:"required,min=1,max=65535"`                        
	Host          string   `koanf:"host" validate:"required,hostname"`                               
	DatabaseName  string   `koanf:"database_name" validate:"required"`                         
	DatabaseType  string   `koanf:"database_type" validate:"required"` 
	MetricsLookup []string `koanf:"metrics_lookup" validate:"required,min=1"`          
	SchemaPath    string   `koanf:"schema_file" validate:"required"`                            
	
}

// AppConfig represents the general application configuration.
type ClientConfig struct {
	HostURL              string `koanf:"host_url" validate:"required,url"`
	ClientVersion        string `koanf:"client_version" validate:"required"`
	FrequencyInNanoseconds int   `koanf:"frequency_in_nanoseconds" validate:"required,min=1"`
}


// ConfigurationManager handles loading configurations.
type ConfigurationManager struct {
    koanf *koanf.Koanf
}

// NewConfigurationManager initializes a new Koanf instance.
func NewConfigurationManager() *ConfigurationManager {
    return &ConfigurationManager{
        koanf: koanf.New("."),
    }
}

// LoadConfig loads a configuration file and unmarshals it into the provided struct.
func (cm *ConfigurationManager) LoadConfig(filePath string, configType interface{}) error {
    // Load configuration from a file.
    if err := cm.koanf.Load(file.Provider(filePath), json.Parser()); err != nil {
        return err
    }

    // Unmarshal into the specified struct.
    if err := cm.koanf.Unmarshal("", configType); err != nil {
        return err
    }

	if err := validateConfig(configType); err != nil {
		log.Fatalf("configuration validation failed: %v", err)
	}

    return nil
}


func validateConfig(config interface{}) error {
	validate := validator.New()
	err := validate.Struct(config)
	if err != nil {
		// Return detailed validation errors
		return fmt.Errorf("validation failed: %v", err)
	}
	return nil
}