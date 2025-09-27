package models

import (
	"time"
)

// ProcessInfo represents system process information
type ProcessInfo struct {
	PID         int32     `json:"pid"`
	PPID        int32     `json:"ppid"`
	Name        string    `json:"name"`
	Status      string    `json:"status"`
	CPU         float64   `json:"cpu"`
	Memory      float64   `json:"memory"`
	MemoryBytes uint64    `json:"memory_bytes"`
	CreateTime  time.Time `json:"create_time"`
	Username    string    `json:"username"`
	Command     string    `json:"command"`
	WorkingDir  string    `json:"working_dir"`
	NumThreads  int32     `json:"num_threads"`
	Nice        int32     `json:"nice"`
	IsRunning   bool      `json:"is_running"`
}

// ProcessFilter represents filtering options for processes
type ProcessFilter struct {
	SearchTerm string `json:"search_term"`
	MinCPU     float64 `json:"min_cpu"`
	MaxCPU     float64 `json:"max_cpu"`
	MinMemory  float64 `json:"min_memory"`
	MaxMemory  float64 `json:"max_memory"`
	Status     string  `json:"status"`
	Username   string  `json:"username"`
	ShowSystem bool    `json:"show_system"`
}

// ProcessSort represents sorting options for processes
type ProcessSort struct {
	Field string `json:"field"` // cpu, memory, pid, name, status
	Order string `json:"order"` // asc, desc
}

// AppConfig represents the application configuration
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
