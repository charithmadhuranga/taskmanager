package storage

import "tappmanager/internal/models"

// Storage defines the interface for data persistence
type Storage interface {
	// Configuration operations
	LoadConfig() (*models.AppConfig, error)
	SaveConfig(config *models.AppConfig) error
	
	// Process data operations
	SaveProcessSnapshot(processes []*models.ProcessInfo) error
	LoadProcessSnapshot() ([]*models.ProcessInfo, error)
	
	// Backup operations
	CreateBackup() error
	RestoreBackup(backupPath string) error
	ListBackups() ([]string, error)
	
	// Export operations
	ExportProcesses(format string) (string, error) // json, csv, xml
	ImportProcesses(data string, format string) error
}
