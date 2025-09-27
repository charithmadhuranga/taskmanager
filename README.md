# Terminal Process Manager

A powerful terminal-based process manager for **macOS, Linux, and Windows** written in **Go** using [`Bubble Tea`](https://github.com/charmbracelet/bubbletea) for the UI and [`gopsutil`](https://github.com/shirou/gopsutil) for cross-platform process information.

## Features

### Process Monitoring
- **Real-time process monitoring** with auto-refresh
- **Cross-platform support** (macOS, Linux, Windows)
- **Detailed process information** (PID, name, status, CPU%, memory%, user, threads, etc.)
- **Process filtering** by CPU usage, memory usage, status, user, and search terms
- **Advanced sorting** by CPU, memory, PID, name, or status
- **System process toggle** to show/hide system processes

### Process Management
- **Kill processes** with confirmation
- **Process details view** with comprehensive information
- **Process search** by name or PID
- **Process tree visualization** (parent-child relationships)
- **Process statistics** and reporting

### Data Management
- **JSON-based data storage** for process snapshots
- **Automatic backups** with configurable retention
- **Data export** in JSON and CSV formats
- **Configuration management** with YAML files
- **Cross-platform data persistence**

## Architecture

The application follows a modular architecture:

```
tappmanager/
├── cmd/                    # Entry point
├── internal/
│   ├── app/               # Application core
│   ├── ui/                # User interface
│   │   ├── views/         # Different UI views
│   │   └── components/    # Reusable components
│   ├── models/            # Data models
│   ├── storage/           # Data persistence
│   ├── services/          # Business logic
│   └── utils/             # Utility functions
├── config/                # Configuration files
└── data/                  # Data storage
```

## Requirements

- Go 1.24+
- Terminal supporting modern TUI features
- Bubble Tea framework for terminal UI

## Installation

```bash
git clone <repository-url>
cd tappmanager
go mod tidy
go build -o tappmanager
```

## Usage

```bash
./tappmanager
```

### Keyboard Shortcuts

- **Ctrl+P** - Switch to Processes view
- **Ctrl+D** - Switch to Details view
- **Ctrl+S** - Switch to Statistics view
- **Ctrl+H** - Show help
- **Ctrl+Q** - Quit application

### Processes View
- **Ctrl+R** - Refresh process list
- **Ctrl+K** - Kill selected process
- **Ctrl+D** - Show process details
- **Ctrl+F** - Filter processes
- **Ctrl+S** - Toggle system processes
- **Ctrl+E** - Export process list
- **Ctrl+B** - Create backup
- **Ctrl+O** - Sort by CPU usage
- **Ctrl+M** - Sort by memory usage
- **Ctrl+P** - Sort by PID
- **Ctrl+N** - Sort by name
- **Ctrl+T** - Sort by status

### Details View
- **Ctrl+R** - Refresh process details
- **Ctrl+K** - Kill selected process
- **↑/↓** - Select previous/next process
- **Ctrl+F** - Search processes

### Statistics View
- **Ctrl+R** - Refresh statistics
- **Ctrl+E** - Export statistics

## Configuration

Configuration is stored in `~/.tappmanager/config.yaml`:

```yaml
data_dir: "~/.tappmanager"
theme: "default"
refresh_rate: 2
auto_backup: true
backup_count: 10
show_system: false
auto_refresh: true
```

## Data Storage

All data is stored in JSON format in the configured data directory:
- `config.json` - Application configuration
- `process_snapshot.json` - Current process snapshot
- `backups/` - Automatic backup files

## Cross-Platform Support

This process manager works seamlessly across:
- **macOS** - Full process monitoring and management
- **Linux** - Complete system process visibility
- **Windows** - Native process handling and monitoring

The application automatically adapts to each platform's process management capabilities and system APIs.
