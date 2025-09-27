package app

import (
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// Config holds the application configuration
type Config struct {
	DataDir     string `mapstructure:"data_dir"`
	Theme       string `mapstructure:"theme"`
	RefreshRate int    `mapstructure:"refresh_rate"`
	AutoBackup  bool   `mapstructure:"auto_backup"`
	BackupCount int    `mapstructure:"backup_count"`
	ShowSystem  bool   `mapstructure:"show_system"`
	AutoRefresh bool   `mapstructure:"auto_refresh"`
}

// DefaultConfig returns the default configuration
func DefaultConfig() *Config {
	homeDir, _ := os.UserHomeDir()
	dataDir := filepath.Join(homeDir, ".tappmanager")
	
	return &Config{
		DataDir:     dataDir,
		Theme:       "default",
		RefreshRate: 2, // seconds
		AutoBackup:  true,
		BackupCount: 10,
		ShowSystem:  false,
		AutoRefresh: true,
	}
}

// LoadConfig loads configuration from file and environment variables
func LoadConfig() (*Config, error) {
	config := DefaultConfig()

	// Set default values
	viper.SetDefault("data_dir", config.DataDir)
	viper.SetDefault("theme", config.Theme)
	viper.SetDefault("refresh_rate", config.RefreshRate)
	viper.SetDefault("auto_backup", config.AutoBackup)
	viper.SetDefault("backup_count", config.BackupCount)

	// Set config file
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("$HOME/.tappmanager")
	viper.AddConfigPath("/etc/tappmanager")

	// Read config file
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
		// Config file not found, use defaults
	}

	// Bind environment variables
	viper.AutomaticEnv()
	viper.BindEnv("data_dir", "TAPPMANAGER_DATA_DIR")
	viper.BindEnv("theme", "TAPPMANAGER_THEME")
	viper.BindEnv("refresh_rate", "TAPPMANAGER_REFRESH_RATE")
	viper.BindEnv("auto_backup", "TAPPMANAGER_AUTO_BACKUP")
	viper.BindEnv("backup_count", "TAPPMANAGER_BACKUP_COUNT")

	// Unmarshal into struct
	if err := viper.Unmarshal(config); err != nil {
		return nil, err
	}

	return config, nil
}

// SaveConfig saves the current configuration to file
func SaveConfig(config *Config) error {
	viper.Set("data_dir", config.DataDir)
	viper.Set("theme", config.Theme)
	viper.Set("refresh_rate", config.RefreshRate)
	viper.Set("auto_backup", config.AutoBackup)
	viper.Set("backup_count", config.BackupCount)

	configDir := filepath.Dir(config.DataDir)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return err
	}

	configFile := filepath.Join(configDir, "config.yaml")
	return viper.WriteConfigAs(configFile)
}
