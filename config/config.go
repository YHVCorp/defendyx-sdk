package config

import (
	"fmt"

	"github.com/YHVCorp/defendyx-sdk/helpers"
)

const (
	DetectionTypeFileHash DetectionType = "ioc_file_hash"

	DBTypeSQLite   DBType = "sqlite"
	DBTypePostgres DBType = "postgres"
	DBTypeCSV      DBType = "csv"
)

type DetectionType string

type Config struct {
	Detections []DetectionConfig `yaml:"detections"`
}

type DBType string

type DetectionConfig struct {
	Name          string         `yaml:"name"`
	Database      DatabaseConfig `yaml:"database"`
	Fields        FieldsConfig   `yaml:"fields"`
	DetectionType string         `yaml:"detection_type"`
}

type DatabaseConfig struct {
	Type   string `yaml:"type"` // sqlite, postgres, csv
	Source string `yaml:"source"`
	Table  string `yaml:"table,omitempty"`
}

type FieldsConfig struct {
	IOC         string   `yaml:"ioc"`
	Informative []string `yaml:"informative"`
}

func LoadConfig(filePath string) (*Config, error) {
	config := &Config{}
	err := helpers.ReadYAML(filePath, config)
	if err != nil {
		return nil, err
	}

	if len(config.Detections) == 0 {
		return nil, fmt.Errorf("no detections found in config file")
	}

	for _, detection := range config.Detections {
		if detection.Name == "" {
			return nil, fmt.Errorf("detection name is required")
		}
		if detection.Database.Type != string(DBTypeSQLite) && detection.Database.Type != string(DBTypePostgres) && detection.Database.Type != string(DBTypeCSV) {
			return nil, fmt.Errorf("invalid database type: %s", detection.Database.Type)
		}
		if detection.Database.Source == "" {
			return nil, fmt.Errorf("database source is required for detection: %s", detection.Name)
		}
		if (detection.Database.Type == string(DBTypeSQLite) || detection.Database.Type == string(DBTypePostgres)) && detection.Database.Table == "" {
			return nil, fmt.Errorf("database table is required for detection: %s", detection.Name)
		}
		if detection.DetectionType != string(DetectionTypeFileHash) {
			return nil, fmt.Errorf("invalid detection type: %s", detection.DetectionType)
		}
		if detection.Fields.IOC == "" {
			return nil, fmt.Errorf("IOC field is required for detection: %s", detection.Name)
		}
	}

	return config, nil
}

func SaveConfig(filePath string, config *Config) error {
	return helpers.WriteYAML(filePath, config)
}
