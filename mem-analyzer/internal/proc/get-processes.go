package proc

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func GetAllProcesses(minKb int64) ([]Process, int, error) {
	var (
		processes           []Process
		numIgnoredProcesses int
		// processMap          = make(map[int]*Process)
	)

	// Read /proc directory
	procDir, err := os.Open("/proc")
	if err != nil {
		return nil, 0, err
	}
	defer procDir.Close()

	// Get all directories in /proc that are numbers (PIDs)
	entries, err := procDir.Readdir(-1)
	if err != nil {
		return nil, 0, err
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		// Check if the directory name is a number (PID)
		pid, err := strconv.Atoi(entry.Name())
		if err != nil {
			continue
		}

		// Get process info
		proc, err := getProcessInfo(pid)
		if err != nil {
			// Skip processes that disappeared
			continue
		}

		if proc.MemoryUsage <= minKb {
			numIgnoredProcesses++
			continue
		}

		processes = append(processes, proc)
		// processMap[proc.PID] = &proc
	}

	// addChildrenToProcesses(processes, processMap)

	return processes, numIgnoredProcesses, nil
}

func getProcessInfo(pid int) (Process, error) {
	proc := Process{PID: pid}

	// Get process name from /proc/[pid]/comm
	commPath := filepath.Join("/proc", strconv.Itoa(pid), "comm")
	commData, err := os.ReadFile(commPath)
	if err != nil {
		return proc, err
	}
	proc.Name = strings.TrimSpace(string(commData))

	// Get PPID from /proc/[pid]/status
	statusPath := filepath.Join("/proc", strconv.Itoa(pid), "status")
	statusData, err := os.ReadFile(statusPath)
	if err != nil {
		return proc, err
	}

	// Parse status file for PPID and memory usage
	lines := strings.Split(string(statusData), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "PPid:") {
			fields := strings.Fields(line)
			if len(fields) >= 2 {
				proc.PPID, _ = strconv.Atoi(fields[1])
			}
		} else if strings.HasPrefix(line, "VmRSS:") {
			// VmRSS is the resident set size, the portion of memory held in RAM
			fields := strings.Fields(line)
			if len(fields) >= 2 {
				proc.MemoryUsage, _ = strconv.ParseInt(fields[1], 10, 64)
			}
		}
	}

	return proc, nil
}

// func addChildrenToProcesses(processes []*Process, processMap map[int]*Process) {
// 	// Group processes by name for orphaned processes
// 	nameToProcess := make(map[string][]*Process)
// 	for i := range processes {
// 		nameToProcess[processes[i].Name] = append(nameToProcess[processes[i].Name], processes[i])
// 	}

// 	processedPIDs := make(map[int]bool)

// 	// First, build the parent-child relationships
// 	for i := range processes {
// 		proc := processes[i]
// 		parent, exists := processMap[proc.PPID]

// 		if exists && proc.PPID != proc.PID { // Avoid self-references
// 			parent.Children = append(parent.Children, proc)
// 			processedPIDs[proc.PID] = true
// 		}
// 	}

// 	// Now add orphaned processes (those without parents in our list)
// 	// or those that have themselves as parent
// 	for i := range processes {
// 		proc := processes[i]
// 		if !processedPIDs[proc.PID] || proc.PPID == proc.PID {
// 			// Group by name if possible
// 			if len(nameToProcess[proc.Name]) > 1 {
// 				// If this is the first occurrence of this name we're processing
// 				isFirst := true
// 				for _, p := range nameToProcess[proc.Name] {
// 					if p.PID < proc.PID && !processedPIDs[p.PID] {
// 						isFirst = false
// 						break
// 					}
// 				}

// 				if isFirst {
// 					processedPIDs[proc.PID] = true

// 					// Add other processes with the same name as children
// 					for _, p := range nameToProcess[proc.Name] {
// 						if p.PID != proc.PID && !processedPIDs[p.PID] {
// 							proc.Children = append(proc.Children, p)
// 							processedPIDs[p.PID] = true
// 						}
// 					}
// 				} else {
// 					proc.IsChild = true
// 				}
// 			} else {
// 				processedPIDs[proc.PID] = true
// 			}
// 		}
// 	}
// }
