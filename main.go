package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/rivo/tview"
	"github.com/shirou/gopsutil/v3/process"
)

func main() {
	// Ensure TERM is set (tview/tcell requirement)
	if os.Getenv("TERM") == "" {
		os.Setenv("TERM", "xterm-256color")
	}

	app := tview.NewApplication()
	table := tview.NewTable().SetBorders(true)

	// Headers
	headers := []string{"PID", "Name", "Status", "CPU%", "Memory%"}
	for i, h := range headers {
		table.SetCell(0, i,
			tview.NewTableCell(fmt.Sprintf("[yellow]%s", h)).
				SetSelectable(false))
	}

	// Function to refresh processes
	refresh := func() {
		procs, err := process.Processes()
		if err != nil {
			log.Printf("Error fetching processes: %v", err)
			return
		}

		// Clear old rows (except header)
		for r := 1; r < table.GetRowCount(); r++ {
			for c := 0; c < table.GetColumnCount(); c++ {
				table.SetCell(r, c, tview.NewTableCell(""))
			}
		}

		for i, p := range procs {
			name, _ := p.Name()

			// FIX: convert []string -> string
			statusSlice, _ := p.Status()
			status := ""
			if len(statusSlice) > 0 {
				status = statusSlice[0]
			}

			cpu, _ := p.CPUPercent()
			mem, _ := p.MemoryPercent()

			row := i + 1
			table.SetCell(row, 0, tview.NewTableCell(strconv.Itoa(int(p.Pid))))
			table.SetCell(row, 1, tview.NewTableCell(name))
			table.SetCell(row, 2, tview.NewTableCell(status))
			table.SetCell(row, 3, tview.NewTableCell(fmt.Sprintf("%.2f", cpu)))
			table.SetCell(row, 4, tview.NewTableCell(fmt.Sprintf("%.2f", mem)))
		}
	}

	// Periodic refresh loop
	go func() {
		for {
			app.QueueUpdateDraw(func() {
				refresh()
			})
			time.Sleep(2 * time.Second)
		}
	}()

	if err := app.SetRoot(table, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}
