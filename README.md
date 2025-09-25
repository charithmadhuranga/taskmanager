# Terminal Process Manager

A simple terminal-based process manager written in **Go** using [`tview`](https://github.com/rivo/tview) for the UI and [`gopsutil`](https://github.com/shirou/gopsutil) for process information.

It allows you to **view running processes**, their **PID, name, status, CPU%, and memory usage** in real-time.

---

## Features

- Terminal-based, interactive UI.
- Auto-refresh every 2 seconds.
- Displays:
    - PID
    - Process Name
    - Status
    - CPU%
    - Memory%
- Cross-platform (Linux, macOS, Windows).
- Mouse support enabled for scrolling.

---

## Requirements

- Go 1.18+
- Terminal supporting `xterm-256color` or similar.

Dependencies:

```bash
go get github.com/rivo/tview
go get github.com/shirou/gopsutil/v3
