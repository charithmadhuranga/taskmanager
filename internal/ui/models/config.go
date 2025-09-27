package models

import "time"

// AppConfig represents the application configuration for Bubble Tea
type AppConfig struct {
	RefreshRate    int           `json:"refresh_rate"`
	ShowSystem     bool          `json:"show_system"`
	DefaultSort    ProcessSort   `json:"default_sort"`
	DefaultFilter  ProcessFilter `json:"default_filter"`
	AutoRefresh    bool          `json:"auto_refresh"`
	Theme          string        `json:"theme"`
	DataDir        string        `json:"data_dir"`
	Version        string        `json:"version"`
	CreatedAt      time.Time     `json:"created_at"`
	UpdatedAt      time.Time     `json:"updated_at"`
}

// ProcessSort represents sorting options for processes
type ProcessSort struct {
	Field string `json:"field"` // cpu, memory, pid, name, status
	Order string `json:"order"` // asc, desc
}

// ProcessFilter represents filtering options for processes
type ProcessFilter struct {
	SearchTerm string  `json:"search_term"`
	MinCPU     float64 `json:"min_cpu"`
	MaxCPU     float64 `json:"max_cpu"`
	MinMemory  float64 `json:"min_memory"`
	MaxMemory  float64 `json:"max_memory"`
	Status     string  `json:"status"`
	Username   string  `json:"username"`
	ShowSystem bool    `json:"show_system"`
}

// NewAppConfig creates a new AppConfig instance with default values
func NewAppConfig() *AppConfig {
	return &AppConfig{
		RefreshRate: 2,
		ShowSystem:  false,
		DefaultSort: ProcessSort{
			Field: "cpu",
			Order: "desc",
		},
		DefaultFilter: ProcessFilter{
			SearchTerm: "",
			MinCPU:     0,
			MaxCPU:     100,
			MinMemory:  0,
			MaxMemory:  100,
			Status:     "",
			Username:   "",
			ShowSystem: false,
		},
		AutoRefresh: true,
		Theme:       "default",
		DataDir:     "~/.tappmanager",
		Version:     "1.0.0",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}
