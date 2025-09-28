package services

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"tappmanager/internal/models"
	"tappmanager/internal/storage"

	"github.com/shirou/gopsutil/v3/process"
)

// ProcessService handles process-related operations
type ProcessService struct {
	storage storage.Storage
}

// NewProcessService creates a new process service
func NewProcessService(storage storage.Storage) *ProcessService {
	return &ProcessService{
		storage: storage,
	}
}

// GetProcesses retrieves all processes with detailed information
func (ps *ProcessService) GetProcesses() ([]*models.ProcessInfo, error) {
	procs, err := process.Processes()
	if err != nil {
		return nil, fmt.Errorf("failed to get processes: %w", err)
	}

	var processInfos []*models.ProcessInfo
	for _, p := range procs {
		info, err := ps.getProcessInfo(p)
		if err != nil {
			continue // Skip processes we can't read
		}
		processInfos = append(processInfos, info)
	}

	// Sort by CPU usage to get more accurate data
	sort.Slice(processInfos, func(i, j int) bool {
		return processInfos[i].CPU > processInfos[j].CPU
	})

	return processInfos, nil
}

// getProcessInfo extracts detailed information from a process
func (ps *ProcessService) getProcessInfo(p *process.Process) (*models.ProcessInfo, error) {
	info := &models.ProcessInfo{
		PID: p.Pid,
	}

	// Get basic information
	if name, err := p.Name(); err == nil {
		info.Name = name
	}

	if ppid, err := p.Ppid(); err == nil {
		info.PPID = ppid
	}

	if status, err := p.Status(); err == nil && len(status) > 0 {
		info.Status = status[0]
	}

	// Get CPU percentage - use a more reliable method
	if cpu, err := p.CPUPercent(); err == nil {
		info.CPU = float64(cpu)
	} else {
		// Try alternative method for CPU percentage
		if times, err := p.Times(); err == nil {
			// This is a basic calculation - in a real implementation you'd want to track previous values
			info.CPU = (times.User + times.System) * 100.0
		} else {
			info.CPU = 0.0
		}
	}

	// Get memory percentage
	if mem, err := p.MemoryPercent(); err == nil {
		info.Memory = float64(mem)
	} else {
		info.Memory = 0.0
	}

	if memInfo, err := p.MemoryInfo(); err == nil {
		info.MemoryBytes = memInfo.RSS
	}

	if createTime, err := p.CreateTime(); err == nil {
		info.CreateTime = time.Unix(0, createTime*int64(time.Millisecond))
	}

	if username, err := p.Username(); err == nil {
		info.Username = username
	}

	if cmdline, err := p.Cmdline(); err == nil {
		info.Command = cmdline
	}

	if cwd, err := p.Cwd(); err == nil {
		info.WorkingDir = cwd
	}

	if numThreads, err := p.NumThreads(); err == nil {
		info.NumThreads = numThreads
	} else {
		info.NumThreads = 0
	}

	if nice, err := p.Nice(); err == nil {
		info.Nice = nice
	} else {
		info.Nice = 0
	}

	// Check if process is running
	info.IsRunning = true

	return info, nil
}

// FilterProcesses filters processes based on criteria
func (ps *ProcessService) FilterProcesses(processes []*models.ProcessInfo, filter *models.ProcessFilter) []*models.ProcessInfo {
	var filtered []*models.ProcessInfo

	for _, proc := range processes {
		// Search term filter
		if filter.SearchTerm != "" {
			searchTerm := strings.ToLower(filter.SearchTerm)
			if !strings.Contains(strings.ToLower(proc.Name), searchTerm) &&
				!strings.Contains(strings.ToLower(proc.Command), searchTerm) &&
				!strings.Contains(strings.ToLower(proc.Username), searchTerm) {
				continue
			}
		}

		// CPU filter
		if proc.CPU < filter.MinCPU || proc.CPU > filter.MaxCPU {
			continue
		}

		// Memory filter
		if proc.Memory < filter.MinMemory || proc.Memory > filter.MaxMemory {
			continue
		}

		// Status filter
		if filter.Status != "" && proc.Status != filter.Status {
			continue
		}

		// Username filter
		if filter.Username != "" && proc.Username != filter.Username {
			continue
		}

		// System process filter
		if !filter.ShowSystem && ps.isSystemProcess(proc) {
			continue
		}

		filtered = append(filtered, proc)
	}

	return filtered
}

// SortProcesses sorts processes based on criteria
func (ps *ProcessService) SortProcesses(processes []*models.ProcessInfo, sortConfig *models.ProcessSort) {
	switch sortConfig.Field {
	case "cpu":
		if sortConfig.Order == "asc" {
			sort.Slice(processes, func(i, j int) bool {
				return processes[i].CPU < processes[j].CPU
			})
		} else {
			sort.Slice(processes, func(i, j int) bool {
				return processes[i].CPU > processes[j].CPU
			})
		}
	case "memory":
		if sortConfig.Order == "asc" {
			sort.Slice(processes, func(i, j int) bool {
				return processes[i].Memory < processes[j].Memory
			})
		} else {
			sort.Slice(processes, func(i, j int) bool {
				return processes[i].Memory > processes[j].Memory
			})
		}
	case "pid":
		if sortConfig.Order == "asc" {
			sort.Slice(processes, func(i, j int) bool {
				return processes[i].PID < processes[j].PID
			})
		} else {
			sort.Slice(processes, func(i, j int) bool {
				return processes[i].PID > processes[j].PID
			})
		}
	case "name":
		if sortConfig.Order == "asc" {
			sort.Slice(processes, func(i, j int) bool {
				return processes[i].Name < processes[j].Name
			})
		} else {
			sort.Slice(processes, func(i, j int) bool {
				return processes[i].Name > processes[j].Name
			})
		}
	case "status":
		if sortConfig.Order == "asc" {
			sort.Slice(processes, func(i, j int) bool {
				return processes[i].Status < processes[j].Status
			})
		} else {
			sort.Slice(processes, func(i, j int) bool {
				return processes[i].Status > processes[j].Status
			})
		}
	case "threads":
		if sortConfig.Order == "asc" {
			sort.Slice(processes, func(i, j int) bool {
				return processes[i].NumThreads < processes[j].NumThreads
			})
		} else {
			sort.Slice(processes, func(i, j int) bool {
				return processes[i].NumThreads > processes[j].NumThreads
			})
		}
	case "nice":
		if sortConfig.Order == "asc" {
			sort.Slice(processes, func(i, j int) bool {
				return processes[i].Nice < processes[j].Nice
			})
		} else {
			sort.Slice(processes, func(i, j int) bool {
				return processes[i].Nice > processes[j].Nice
			})
		}
	case "user":
		if sortConfig.Order == "asc" {
			sort.Slice(processes, func(i, j int) bool {
				return processes[i].Username < processes[j].Username
			})
		} else {
			sort.Slice(processes, func(i, j int) bool {
				return processes[i].Username > processes[j].Username
			})
		}
	}
}

// isSystemProcess determines if a process is a system process
func (ps *ProcessService) isSystemProcess(proc *models.ProcessInfo) bool {
	// Common system process names
	systemProcesses := []string{
		"kernel_task", "launchd", "kextd", "mds", "mdworker",
		"WindowServer", "loginwindow", "UserEventAgent", "configd",
		"syslogd", "kdc", "distnoted", "notifyd", "securityd",
		"coreaudiod", "coreduetd", "fseventsd", "locationd",
		"powerd", "thermalmonitord", "wifid", "bluetoothd",
		"hidd", "pboard", "sharingd", "usbmuxd", "com.apple",
	}

	for _, sysProc := range systemProcesses {
		if strings.Contains(proc.Name, sysProc) {
			return true
		}
	}

	// Check for system users
	systemUsers := []string{"root", "daemon", "nobody", "system"}
	for _, sysUser := range systemUsers {
		if proc.Username == sysUser {
			return true
		}
	}

	return false
}

// KillProcess attempts to kill a process
func (ps *ProcessService) KillProcess(pid int32) error {
	proc, err := process.NewProcess(pid)
	if err != nil {
		return fmt.Errorf("failed to get process %d: %w", pid, err)
	}

	if err := proc.Kill(); err != nil {
		return fmt.Errorf("failed to kill process %d: %w", pid, err)
	}

	return nil
}

// GetProcessTree returns a hierarchical view of processes
func (ps *ProcessService) GetProcessTree(processes []*models.ProcessInfo) map[int32][]*models.ProcessInfo {
	tree := make(map[int32][]*models.ProcessInfo)
	
	for _, proc := range processes {
		tree[proc.PPID] = append(tree[proc.PPID], proc)
	}
	
	return tree
}

// GetProcessStats returns statistics about the processes
func (ps *ProcessService) GetProcessStats(processes []*models.ProcessInfo) map[string]interface{} {
	stats := make(map[string]interface{})
	
	totalProcesses := len(processes)
	runningProcesses := 0
	totalCPU := 0.0
	totalMemory := 0.0
	
	statusCounts := make(map[string]int)
	userCounts := make(map[string]int)
	
	for _, proc := range processes {
		if proc.IsRunning {
			runningProcesses++
		}
		
		totalCPU += proc.CPU
		totalMemory += proc.Memory
		
		statusCounts[proc.Status]++
		userCounts[proc.Username]++
	}
	
	stats["total_processes"] = totalProcesses
	stats["running_processes"] = runningProcesses
	stats["total_cpu"] = totalCPU
	stats["total_memory"] = totalMemory
	stats["status_counts"] = statusCounts
	stats["user_counts"] = userCounts
	
	return stats
}
