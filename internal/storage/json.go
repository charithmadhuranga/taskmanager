package storage

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"tappmanager/internal/models"
)

// JSONStorage implements Storage interface using JSON files
type JSONStorage struct {
	dataDir    string
	backupDir  string
	config     *models.AppConfig
	processes  []*models.ProcessInfo
}

// NewJSONStorage creates a new JSON storage instance
func NewJSONStorage(dataDir string) *JSONStorage {
	backupDir := filepath.Join(dataDir, "backups")
	return &JSONStorage{
		dataDir:   dataDir,
		backupDir: backupDir,
		config:    models.NewAppConfig(),
		processes: []*models.ProcessInfo{},
	}
}

// ensureDirectories creates necessary directories if they don't exist
func (s *JSONStorage) ensureDirectories() error {
	if err := os.MkdirAll(s.dataDir, 0755); err != nil {
		return fmt.Errorf("failed to create data directory: %w", err)
	}
	if err := os.MkdirAll(s.backupDir, 0755); err != nil {
		return fmt.Errorf("failed to create backup directory: %w", err)
	}
	return nil
}

// LoadConfig loads the application configuration
func (s *JSONStorage) LoadConfig() (*models.AppConfig, error) {
	if err := s.ensureDirectories(); err != nil {
		return nil, err
	}

	configFile := filepath.Join(s.dataDir, "config.json")
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		// Create default config if file doesn't exist
		s.config = models.NewAppConfig()
		return s.config, nil
	}

	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	var config models.AppConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	s.config = &config
	return s.config, nil
}

// SaveConfig saves the application configuration
func (s *JSONStorage) SaveConfig(config *models.AppConfig) error {
	if err := s.ensureDirectories(); err != nil {
		return err
	}

	config.UpdatedAt = time.Now()
	configFile := filepath.Join(s.dataDir, "config.json")
	jsonData, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := ioutil.WriteFile(configFile, jsonData, 0644); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	s.config = config
	return nil
}

// SaveProcessSnapshot saves a snapshot of current processes
func (s *JSONStorage) SaveProcessSnapshot(processes []*models.ProcessInfo) error {
	if err := s.ensureDirectories(); err != nil {
		return err
	}

	snapshotFile := filepath.Join(s.dataDir, "process_snapshot.json")
	jsonData, err := json.MarshalIndent(processes, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal process snapshot: %w", err)
	}

	if err := ioutil.WriteFile(snapshotFile, jsonData, 0644); err != nil {
		return fmt.Errorf("failed to write process snapshot: %w", err)
	}

	s.processes = processes
	return nil
}

// LoadProcessSnapshot loads the last saved process snapshot
func (s *JSONStorage) LoadProcessSnapshot() ([]*models.ProcessInfo, error) {
	if err := s.ensureDirectories(); err != nil {
		return nil, err
	}

	snapshotFile := filepath.Join(s.dataDir, "process_snapshot.json")
	if _, err := os.Stat(snapshotFile); os.IsNotExist(err) {
		return []*models.ProcessInfo{}, nil
	}

	data, err := ioutil.ReadFile(snapshotFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read process snapshot: %w", err)
	}

	var processes []*models.ProcessInfo
	if err := json.Unmarshal(data, &processes); err != nil {
		return nil, fmt.Errorf("failed to unmarshal process snapshot: %w", err)
	}

	s.processes = processes
	return processes, nil
}

// CreateBackup creates a backup of the current data
func (s *JSONStorage) CreateBackup() error {
	if err := s.ensureDirectories(); err != nil {
		return err
	}

	backupFile := filepath.Join(s.backupDir, fmt.Sprintf("backup_%s.json", time.Now().Format("20060102_150405")))
	
	backupData := map[string]interface{}{
		"config":    s.config,
		"processes": s.processes,
		"timestamp": time.Now(),
	}

	jsonData, err := json.MarshalIndent(backupData, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal backup data: %w", err)
	}

	if err := ioutil.WriteFile(backupFile, jsonData, 0644); err != nil {
		return fmt.Errorf("failed to write backup file: %w", err)
	}

	return nil
}

// RestoreBackup restores data from a backup file
func (s *JSONStorage) RestoreBackup(backupPath string) error {
	data, err := ioutil.ReadFile(backupPath)
	if err != nil {
		return fmt.Errorf("failed to read backup file: %w", err)
	}

	var backupData map[string]interface{}
	if err := json.Unmarshal(data, &backupData); err != nil {
		return fmt.Errorf("failed to unmarshal backup data: %w", err)
	}

	// Restore config
	if configData, ok := backupData["config"]; ok {
		configBytes, err := json.Marshal(configData)
		if err != nil {
			return fmt.Errorf("failed to marshal config from backup: %w", err)
		}
		var config models.AppConfig
		if err := json.Unmarshal(configBytes, &config); err != nil {
			return fmt.Errorf("failed to unmarshal config from backup: %w", err)
		}
		s.config = &config
	}

	// Restore processes
	if processesData, ok := backupData["processes"]; ok {
		processesBytes, err := json.Marshal(processesData)
		if err != nil {
			return fmt.Errorf("failed to marshal processes from backup: %w", err)
		}
		var processes []*models.ProcessInfo
		if err := json.Unmarshal(processesBytes, &processes); err != nil {
			return fmt.Errorf("failed to unmarshal processes from backup: %w", err)
		}
		s.processes = processes
	}

	return nil
}

// ListBackups returns a list of available backup files
func (s *JSONStorage) ListBackups() ([]string, error) {
	if err := s.ensureDirectories(); err != nil {
		return nil, err
	}

	files, err := ioutil.ReadDir(s.backupDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read backup directory: %w", err)
	}

	var backups []string
	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".json" {
			backups = append(backups, filepath.Join(s.backupDir, file.Name()))
		}
	}

	return backups, nil
}

// ExportProcesses exports processes in the specified format
func (s *JSONStorage) ExportProcesses(format string) (string, error) {
	if err := s.ensureDirectories(); err != nil {
		return "", err
	}

	timestamp := time.Now().Format("20060102_150405")
	
	switch format {
	case "json":
		filename := filepath.Join(s.dataDir, fmt.Sprintf("processes_export_%s.json", timestamp))
		jsonData, err := json.MarshalIndent(s.processes, "", "  ")
		if err != nil {
			return "", fmt.Errorf("failed to marshal processes for export: %w", err)
		}
		if err := ioutil.WriteFile(filename, jsonData, 0644); err != nil {
			return "", fmt.Errorf("failed to write export file: %w", err)
		}
		return filename, nil
		
	case "csv":
		filename := filepath.Join(s.dataDir, fmt.Sprintf("processes_export_%s.csv", timestamp))
		file, err := os.Create(filename)
		if err != nil {
			return "", fmt.Errorf("failed to create CSV file: %w", err)
		}
		defer file.Close()

		writer := csv.NewWriter(file)
		defer writer.Flush()

		// Write header
		header := []string{"PID", "PPID", "Name", "Status", "CPU%", "Memory%", "MemoryBytes", "Username", "Command", "WorkingDir", "NumThreads", "Nice", "CreateTime"}
		if err := writer.Write(header); err != nil {
			return "", fmt.Errorf("failed to write CSV header: %w", err)
		}

		// Write data
		for _, proc := range s.processes {
			record := []string{
				strconv.Itoa(int(proc.PID)),
				strconv.Itoa(int(proc.PPID)),
				proc.Name,
				proc.Status,
				fmt.Sprintf("%.2f", proc.CPU),
				fmt.Sprintf("%.2f", proc.Memory),
				strconv.FormatUint(proc.MemoryBytes, 10),
				proc.Username,
				proc.Command,
				proc.WorkingDir,
				strconv.Itoa(int(proc.NumThreads)),
				strconv.Itoa(int(proc.Nice)),
				proc.CreateTime.Format(time.RFC3339),
			}
			if err := writer.Write(record); err != nil {
				return "", fmt.Errorf("failed to write CSV record: %w", err)
			}
		}
		return filename, nil
		
	default:
		return "", fmt.Errorf("unsupported export format: %s", format)
	}
}

// ImportProcesses imports processes from the specified data and format
func (s *JSONStorage) ImportProcesses(data string, format string) error {
	switch format {
	case "json":
		var processes []*models.ProcessInfo
		if err := json.Unmarshal([]byte(data), &processes); err != nil {
			return fmt.Errorf("failed to unmarshal JSON data: %w", err)
		}
		s.processes = processes
		return s.SaveProcessSnapshot(processes)
		
	case "csv":
		reader := csv.NewReader(strings.NewReader(data))
		records, err := reader.ReadAll()
		if err != nil {
			return fmt.Errorf("failed to read CSV data: %w", err)
		}
		
		if len(records) < 2 {
			return fmt.Errorf("CSV data must have at least a header and one data row")
		}
		
		var processes []*models.ProcessInfo
		for i, record := range records[1:] { // Skip header
			if len(record) < 13 {
				return fmt.Errorf("CSV record %d has insufficient columns", i+1)
			}
			
			pid, _ := strconv.Atoi(record[0])
			ppid, _ := strconv.Atoi(record[1])
			cpu, _ := strconv.ParseFloat(record[4], 64)
			memory, _ := strconv.ParseFloat(record[5], 64)
			memoryBytes, _ := strconv.ParseUint(record[6], 10, 64)
			numThreads, _ := strconv.Atoi(record[10])
			nice, _ := strconv.Atoi(record[11])
			createTime, _ := time.Parse(time.RFC3339, record[12])
			
			process := &models.ProcessInfo{
				PID:         int32(pid),
				PPID:        int32(ppid),
				Name:        record[2],
				Status:      record[3],
				CPU:         cpu,
				Memory:      memory,
				MemoryBytes: memoryBytes,
				Username:    record[7],
				Command:     record[8],
				WorkingDir:  record[9],
				NumThreads:  int32(numThreads),
				Nice:        int32(nice),
				CreateTime:  createTime,
				IsRunning:   true,
			}
			processes = append(processes, process)
		}
		
		s.processes = processes
		return s.SaveProcessSnapshot(processes)
		
	default:
		return fmt.Errorf("unsupported import format: %s", format)
	}
}