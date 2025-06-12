package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

type Process struct {
	PID         int
	PPID        int
	Name        string
	MemoryUsage int64 // in KB
	Children    []*Process
	PIDs        []int // Used for consolidated processes with same name
}

// For the summary section
type ProcessGroup struct {
	Name        string
	PIDs        []int
	MemoryUsage int64
}

func main() {
	// Get all processes
	processes, err := getAllProcesses()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting processes: %v\n", err)
		os.Exit(1)
	}

	// Create a map for easy lookup
	processMap := make(map[int]*Process)
	for i := range processes {
		processMap[processes[i].PID] = &processes[i]
	}

	// Build process tree
	rootProcesses := buildProcessTree(processes, processMap)

	// Calculate memory usage including children
	totalSystemMemory := int64(0)
	for _, proc := range rootProcesses {
		calculateTotalMemory(proc)
		totalSystemMemory += proc.MemoryUsage
	}

	// Section 1: Print the process tree with memory usage
	fmt.Println("=== PROCESS TREE ===")
	for _, proc := range rootProcesses {
		printProcessTree(proc, 0, proc.MemoryUsage)
	}

	fmt.Println("\n=== MEMORY USAGE SUMMARY ===")
	fmt.Printf("Total monitored memory: %s\n\n", formatMemory(totalSystemMemory))

	// Create process groups for the summary section
	processGroups := createProcessGroups(rootProcesses)

	// Sort by memory usage (descending)
	sort.Slice(processGroups, func(i, j int) bool {
		if processGroups[i].MemoryUsage == processGroups[j].MemoryUsage {
			// If memory usage is the same, sort by lowest PID
			return min(processGroups[i].PIDs) < min(processGroups[j].PIDs)
		}
		return processGroups[i].MemoryUsage > processGroups[j].MemoryUsage
	})

	// Print memory usage summary
	for _, group := range processGroups {
		percentage := float64(group.MemoryUsage) / float64(totalSystemMemory) * 100
		pids := make([]string, len(group.PIDs))
		for i, pid := range group.PIDs {
			pids[i] = strconv.Itoa(pid)
		}
		pidStr := strings.Join(pids, ", ")
		fmt.Printf("%s [%s]: %s (%.1f%%)\n",
			group.Name,
			pidStr,
			formatMemory(group.MemoryUsage),
			percentage)
	}
}

func min(nums []int) int {
	if len(nums) == 0 {
		return 0
	}

	m := nums[0]
	for _, n := range nums {
		if n < m {
			m = n
		}
	}
	return m
}

func getAllProcesses() ([]Process, error) {
	var processes []Process

	// Read /proc directory
	procDir, err := os.Open("/proc")
	if err != nil {
		return nil, err
	}
	defer procDir.Close()

	// Get all directories in /proc that are numbers (PIDs)
	entries, err := procDir.Readdir(-1)
	if err != nil {
		return nil, err
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

		processes = append(processes, proc)
	}

	return processes, nil
}

func getProcessInfo(pid int) (Process, error) {
	proc := Process{
		PID:  pid,
		PIDs: []int{pid}, // Initialize PIDs with the process's own PID
	}

	// Get process name from /proc/[pid]/comm
	commPath := filepath.Join("/proc", strconv.Itoa(pid), "comm")
	commData, err := ioutil.ReadFile(commPath)
	if err != nil {
		return proc, err
	}
	proc.Name = strings.TrimSpace(string(commData))

	// Get PPID from /proc/[pid]/status
	statusPath := filepath.Join("/proc", strconv.Itoa(pid), "status")
	statusData, err := ioutil.ReadFile(statusPath)
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

func buildProcessTree(processes []Process, processMap map[int]*Process) []*Process {
	// Track processed PIDs
	processedPIDs := make(map[int]bool)

	// First pass: build the tree based on parent-child relationships
	for i := range processes {
		proc := &processes[i]
		parent, exists := processMap[proc.PPID]

		if exists && proc.PPID != proc.PID { // Avoid self-references
			if parent.Name == proc.Name {
				// Same name process: merge with parent
				shouldMerge := true

				// Check if this process has children with different names
				for j := range processes {
					childProc := &processes[j]
					if childProc.PPID == proc.PID && childProc.Name != proc.Name {
						// This process has a child with a different name, don't merge
						shouldMerge = false
						break
					}
				}

				if shouldMerge {
					// Merge with parent
					parent.MemoryUsage += proc.MemoryUsage
					parent.PIDs = append(parent.PIDs, proc.PID)
					processedPIDs[proc.PID] = true

					// If this process has children with the same name, also merge them
					for j := range processes {
						childProc := &processes[j]
						if childProc.PPID == proc.PID && childProc.Name == proc.Name {
							parent.MemoryUsage += childProc.MemoryUsage
							parent.PIDs = append(parent.PIDs, childProc.PID)
							processedPIDs[childProc.PID] = true
						}
					}
				} else {
					// Add as regular child
					parent.Children = append(parent.Children, proc)
					processedPIDs[proc.PID] = true
				}
			} else {
				// Different name: add as child
				parent.Children = append(parent.Children, proc)
				processedPIDs[proc.PID] = true
			}
		}
	}

	// Second pass: collect root processes (those not processed yet)
	var rootProcesses []*Process
	nameToRoots := make(map[string][]*Process)

	for i := range processes {
		proc := &processes[i]
		if !processedPIDs[proc.PID] {
			nameToRoots[proc.Name] = append(nameToRoots[proc.Name], proc)
		}
	}

	// Consolidate root processes with the same name
	for _, procs := range nameToRoots {
		if len(procs) == 0 {
			continue
		}

		// Sort by PID to ensure deterministic behavior
		sort.Slice(procs, func(i, j int) bool {
			return procs[i].PID < procs[j].PID
		})

		// Use first process as the representative
		rootProc := procs[0]

		// Merge all processes with the same name
		if len(procs) > 1 {
			for i := 1; i < len(procs); i++ {
				rootProc.MemoryUsage += procs[i].MemoryUsage
				rootProc.PIDs = append(rootProc.PIDs, procs[i].PID)

				// Merge the children too
				rootProc.Children = append(rootProc.Children, procs[i].Children...)
			}
		}

		rootProcesses = append(rootProcesses, rootProc)
	}

	// Sort root processes by lowest PID in their PID list
	sort.Slice(rootProcesses, func(i, j int) bool {
		return min(rootProcesses[i].PIDs) < min(rootProcesses[j].PIDs)
	})

	return rootProcesses
}

func calculateTotalMemory(proc *Process) int64 {
	totalMem := proc.MemoryUsage

	for _, child := range proc.Children {
		totalMem += calculateTotalMemory(child)
	}

	return totalMem
}

func formatMemory(kb int64) string {
	if kb < 1024 {
		return fmt.Sprintf("%d KB", kb)
	} else if kb < 1024*1024 {
		return fmt.Sprintf("%.2f MB", float64(kb)/1024)
	} else {
		return fmt.Sprintf("%.2f GB", float64(kb)/(1024*1024))
	}
}

func printProcessTree(proc *Process, level int, parentMem int64) {
	indent := strings.Repeat("  ", level)

	// Sort PIDs for consistent output
	sort.Ints(proc.PIDs)

	// Convert PIDs to string
	pids := make([]string, len(proc.PIDs))
	for i, pid := range proc.PIDs {
		pids[i] = strconv.Itoa(pid)
	}
	pidStr := strings.Join(pids, ", ")

	if level == 0 {
		// Root process
		fmt.Printf("%s%s [%s]: %s\n", indent, proc.Name, pidStr, formatMemory(proc.MemoryUsage))
	} else {
		// Child process
		percentage := float64(proc.MemoryUsage) / float64(parentMem) * 100
		fmt.Printf("%s%s [%s]: %s (%.1f%% of parent)\n",
			indent, proc.Name, pidStr, formatMemory(proc.MemoryUsage), percentage)
	}

	// Sort children by memory usage (descending)
	sort.Slice(proc.Children, func(i, j int) bool {
		return proc.Children[i].MemoryUsage > proc.Children[j].MemoryUsage
	})

	// Print children
	for _, child := range proc.Children {
		printProcessTree(child, level+1, proc.MemoryUsage)
	}
}

func createProcessGroups(rootProcesses []*Process) []ProcessGroup {
	// Create a map to group processes by name
	nameToGroup := make(map[string]*ProcessGroup)

	// Helper function to recursively collect all processes
	var collectProcesses func(proc *Process)
	collectProcesses = func(proc *Process) {
		group, exists := nameToGroup[proc.Name]
		if !exists {
			group = &ProcessGroup{
				Name: proc.Name,
			}
			nameToGroup[proc.Name] = group
		}

		// Add this process's PIDs and memory to the group
		group.PIDs = append(group.PIDs, proc.PIDs...)
		group.MemoryUsage += proc.MemoryUsage

		// Process children
		for _, child := range proc.Children {
			collectProcesses(child)
		}
	}

	// Collect all processes
	for _, root := range rootProcesses {
		collectProcesses(root)
	}

	// Convert map to slice
	var groups []ProcessGroup
	for _, group := range nameToGroup {
		// Sort PIDs for consistent output
		sort.Ints(group.PIDs)
		groups = append(groups, *group)
	}

	return groups
}

// package main v2

// import (
// 	"fmt"
// 	"io/ioutil"
// 	"os"
// 	"path/filepath"
// 	"sort"
// 	"strconv"
// 	"strings"
// )

// type Process struct {
// 	PID         int
// 	PPID        int
// 	Name        string
// 	MemoryUsage int64 // in KB
// 	Children    []*Process
// }

// func main() {
// 	// Get all processes
// 	processes, err := getAllProcesses()
// 	if err != nil {
// 		fmt.Fprintf(os.Stderr, "Error getting processes: %v\n", err)
// 		os.Exit(1)
// 	}

// 	// Create a map for easy lookup
// 	processMap := make(map[int]*Process)
// 	for i := range processes {
// 		processMap[processes[i].PID] = &processes[i]
// 	}

// 	// Build process tree
// 	rootProcesses := buildProcessTree(processes, processMap)

// 	// Calculate memory usage including children
// 	for _, proc := range rootProcesses {
// 		calculateTotalMemory(proc)
// 	}

// 	// Print the process tree with memory usage
// 	for _, proc := range rootProcesses {
// 		printProcessTree(proc, 0, proc.MemoryUsage)
// 	}
// }

// func getAllProcesses() ([]Process, error) {
// 	var processes []Process

// 	// Read /proc directory
// 	procDir, err := os.Open("/proc")
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer procDir.Close()

// 	// Get all directories in /proc that are numbers (PIDs)
// 	entries, err := procDir.Readdir(-1)
// 	if err != nil {
// 		return nil, err
// 	}

// 	for _, entry := range entries {
// 		if !entry.IsDir() {
// 			continue
// 		}

// 		// Check if the directory name is a number (PID)
// 		pid, err := strconv.Atoi(entry.Name())
// 		if err != nil {
// 			continue
// 		}

// 		// Get process info
// 		proc, err := getProcessInfo(pid)
// 		if err != nil {
// 			// Skip processes that disappeared
// 			continue
// 		}

// 		processes = append(processes, proc)
// 	}

// 	return processes, nil
// }

// func getProcessInfo(pid int) (Process, error) {
// 	proc := Process{PID: pid}

// 	// Get process name from /proc/[pid]/comm
// 	commPath := filepath.Join("/proc", strconv.Itoa(pid), "comm")
// 	commData, err := ioutil.ReadFile(commPath)
// 	if err != nil {
// 		return proc, err
// 	}
// 	proc.Name = strings.TrimSpace(string(commData))

// 	// Get PPID from /proc/[pid]/status
// 	statusPath := filepath.Join("/proc", strconv.Itoa(pid), "status")
// 	statusData, err := ioutil.ReadFile(statusPath)
// 	if err != nil {
// 		return proc, err
// 	}

// 	// Parse status file for PPID and memory usage
// 	lines := strings.Split(string(statusData), "\n")
// 	for _, line := range lines {
// 		if strings.HasPrefix(line, "PPid:") {
// 			fields := strings.Fields(line)
// 			if len(fields) >= 2 {
// 				proc.PPID, _ = strconv.Atoi(fields[1])
// 			}
// 		} else if strings.HasPrefix(line, "VmRSS:") {
// 			// VmRSS is the resident set size, the portion of memory held in RAM
// 			fields := strings.Fields(line)
// 			if len(fields) >= 2 {
// 				proc.MemoryUsage, _ = strconv.ParseInt(fields[1], 10, 64)
// 			}
// 		}
// 	}

// 	return proc, nil
// }

// func buildProcessTree(processes []Process, processMap map[int]*Process) []*Process {
// 	// Group processes by name
// 	nameToProcess := make(map[string][]*Process)
// 	for i := range processes {
// 		nameToProcess[processes[i].Name] = append(nameToProcess[processes[i].Name], &processes[i])
// 	}

// 	// Create a set of root processes
// 	var rootProcesses []*Process
// 	processedPIDs := make(map[int]bool)

// 	// First pass: identify parent-child relationships where names differ
// 	for i := range processes {
// 		proc := &processes[i]
// 		parent, exists := processMap[proc.PPID]

// 		if exists && proc.PPID != proc.PID && parent.Name != proc.Name {
// 			// Only add as child if names are different
// 			parent.Children = append(parent.Children, proc)
// 			processedPIDs[proc.PID] = true
// 		}
// 	}

// 	// Second pass: consolidate processes with the same name
// 	// and handle orphaned processes
// 	// Create groups of processes with the same name that haven't been processed yet
// 	sameNameGroups := make(map[string][]*Process)
// 	for i := range processes {
// 		proc := &processes[i]
// 		if !processedPIDs[proc.PID] {
// 			sameNameGroups[proc.Name] = append(sameNameGroups[proc.Name], proc)
// 		}
// 	}

// 	// Process each group
// 	for name, procs := range sameNameGroups {
// 		if len(procs) == 0 {
// 			continue
// 		}

// 		// Sort processes by PID to ensure consistent results
// 		sort.Slice(procs, func(i, j int) bool {
// 			return procs[i].PID < procs[j].PID
// 		})

// 		// Use the first process as the "representative" for this name
// 		rootProc := procs[0]
// 		processedPIDs[rootProc.PID] = true

// 		// Consolidate memory from all processes with the same name
// 		for i := 1; i < len(procs); i++ {
// 			rootProc.MemoryUsage += procs[i].MemoryUsage
// 			processedPIDs[procs[i].PID] = true
// 		}

// 		// Store all PIDs in the Name field (we'll parse this later)
// 		pids := make([]string, len(procs))
// 		for i, p := range procs {
// 			pids[i] = strconv.Itoa(p.PID)
// 		}
// 		rootProc.Name = name + " [" + strings.Join(pids, ", ") + "]"

// 		rootProcesses = append(rootProcesses, rootProc)
// 	}

// 	// Sort root processes by memory usage (descending)
// 	sort.Slice(rootProcesses, func(i, j int) bool {
// 		return rootProcesses[i].MemoryUsage > rootProcesses[j].MemoryUsage
// 	})

// 	return rootProcesses
// }

// func calculateTotalMemory(proc *Process) int64 {
// 	totalMem := proc.MemoryUsage

// 	for _, child := range proc.Children {
// 		totalMem += calculateTotalMemory(child)
// 	}

// 	proc.MemoryUsage = totalMem
// 	return totalMem
// }

// func formatMemory(kb int64) string {
// 	if kb < 1024 {
// 		return fmt.Sprintf("%d KB", kb)
// 	} else if kb < 1024*1024 {
// 		return fmt.Sprintf("%.2f MB", float64(kb)/1024)
// 	} else {
// 		return fmt.Sprintf("%.2f GB", float64(kb)/(1024*1024))
// 	}
// }

// func printProcessTree(proc *Process, level int, parentMem int64) {
// 	indent := strings.Repeat("  ", level)

// 	if level == 0 {
// 		// Root process
// 		fmt.Printf("%s%s (PID: %d): %s\n", indent, proc.Name, proc.PID, formatMemory(proc.MemoryUsage))
// 	} else {
// 		// Child process
// 		percentage := float64(proc.MemoryUsage) / float64(parentMem) * 100
// 		fmt.Printf("%s%s (PID: %d): %s (%.1f%% of parent)\n", indent, proc.Name, proc.PID, formatMemory(proc.MemoryUsage), percentage)
// 	}

// 	// Sort children by memory usage (descending)
// 	sort.Slice(proc.Children, func(i, j int) bool {
// 		return proc.Children[i].MemoryUsage > proc.Children[j].MemoryUsage
// 	})

// 	// Print children
// 	for _, child := range proc.Children {
// 		printProcessTree(child, level+1, proc.MemoryUsage)
// 	}
// }

// package main v1

// import (
// 	"fmt"
// 	"io/ioutil"
// 	"os"
// 	"path/filepath"
// 	"sort"
// 	"strconv"
// 	"strings"
// )

// type Process struct {
// 	PID         int
// 	PPID        int
// 	Name        string
// 	MemoryUsage int64 // in KB
// 	Children    []*Process
// }

// func main() {
// 	// Get all processes
// 	processes, err := getAllProcesses()
// 	if err != nil {
// 		fmt.Fprintf(os.Stderr, "Error getting processes: %v\n", err)
// 		os.Exit(1)
// 	}

// 	// Create a map for easy lookup
// 	processMap := make(map[int]*Process)
// 	for i := range processes {
// 		processMap[processes[i].PID] = &processes[i]
// 	}

// 	// Build process tree
// 	rootProcesses := buildProcessTree(processes, processMap)

// 	// Calculate memory usage including children
// 	for _, proc := range rootProcesses {
// 		calculateTotalMemory(proc)
// 	}

// 	// Print the process tree with memory usage
// 	for _, proc := range rootProcesses {
// 		printProcessTree(proc, 0, proc.MemoryUsage)
// 	}
// }

// func getAllProcesses() ([]Process, error) {
// 	var processes []Process

// 	// Read /proc directory
// 	procDir, err := os.Open("/proc")
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer procDir.Close()

// 	// Get all directories in /proc that are numbers (PIDs)
// 	entries, err := procDir.Readdir(-1)
// 	if err != nil {
// 		return nil, err
// 	}

// 	for _, entry := range entries {
// 		if !entry.IsDir() {
// 			continue
// 		}

// 		// Check if the directory name is a number (PID)
// 		pid, err := strconv.Atoi(entry.Name())
// 		if err != nil {
// 			continue
// 		}

// 		// Get process info
// 		proc, err := getProcessInfo(pid)
// 		if err != nil {
// 			// Skip processes that disappeared
// 			continue
// 		}

// 		processes = append(processes, proc)
// 	}

// 	return processes, nil
// }

// func getProcessInfo(pid int) (Process, error) {
// 	proc := Process{PID: pid}

// 	// Get process name from /proc/[pid]/comm
// 	commPath := filepath.Join("/proc", strconv.Itoa(pid), "comm")
// 	commData, err := ioutil.ReadFile(commPath)
// 	if err != nil {
// 		return proc, err
// 	}
// 	proc.Name = strings.TrimSpace(string(commData))

// 	// Get PPID from /proc/[pid]/status
// 	statusPath := filepath.Join("/proc", strconv.Itoa(pid), "status")
// 	statusData, err := ioutil.ReadFile(statusPath)
// 	if err != nil {
// 		return proc, err
// 	}

// 	// Parse status file for PPID and memory usage
// 	lines := strings.Split(string(statusData), "\n")
// 	for _, line := range lines {
// 		if strings.HasPrefix(line, "PPid:") {
// 			fields := strings.Fields(line)
// 			if len(fields) >= 2 {
// 				proc.PPID, _ = strconv.Atoi(fields[1])
// 			}
// 		} else if strings.HasPrefix(line, "VmRSS:") {
// 			// VmRSS is the resident set size, the portion of memory held in RAM
// 			fields := strings.Fields(line)
// 			if len(fields) >= 2 {
// 				proc.MemoryUsage, _ = strconv.ParseInt(fields[1], 10, 64)
// 			}
// 		}
// 	}

// 	return proc, nil
// }

// func buildProcessTree(processes []Process, processMap map[int]*Process) []*Process {
// 	// Group processes by name for orphaned processes
// 	nameToProcess := make(map[string][]*Process)
// 	for i := range processes {
// 		nameToProcess[processes[i].Name] = append(nameToProcess[processes[i].Name], &processes[i])
// 	}

// 	// Create a set of root processes
// 	var rootProcesses []*Process
// 	processedPIDs := make(map[int]bool)

// 	// First, build the tree based on parent-child relationships
// 	for i := range processes {
// 		proc := &processes[i]
// 		parent, exists := processMap[proc.PPID]

// 		if exists && proc.PPID != proc.PID { // Avoid self-references
// 			parent.Children = append(parent.Children, proc)
// 			processedPIDs[proc.PID] = true
// 		}
// 	}

// 	// Now add orphaned processes (those without parents in our list)
// 	// or those that have themselves as parent
// 	for i := range processes {
// 		proc := &processes[i]
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
// 					rootProcesses = append(rootProcesses, proc)
// 					processedPIDs[proc.PID] = true

// 					// Add other processes with the same name as children
// 					for _, p := range nameToProcess[proc.Name] {
// 						if p.PID != proc.PID && !processedPIDs[p.PID] {
// 							proc.Children = append(proc.Children, p)
// 							processedPIDs[p.PID] = true
// 						}
// 					}
// 				}
// 			} else {
// 				// No other process with the same name
// 				rootProcesses = append(rootProcesses, proc)
// 				processedPIDs[proc.PID] = true
// 			}
// 		}
// 	}

// 	// Sort root processes by memory usage (descending)
// 	sort.Slice(rootProcesses, func(i, j int) bool {
// 		return rootProcesses[i].MemoryUsage > rootProcesses[j].MemoryUsage
// 	})

// 	return rootProcesses
// }

// func calculateTotalMemory(proc *Process) int64 {
// 	totalMem := proc.MemoryUsage

// 	for _, child := range proc.Children {
// 		totalMem += calculateTotalMemory(child)
// 	}

// 	proc.MemoryUsage = totalMem
// 	return totalMem
// }

// func formatMemory(kb int64) string {
// 	if kb < 1024 {
// 		return fmt.Sprintf("%d KB", kb)
// 	} else if kb < 1024*1024 {
// 		return fmt.Sprintf("%.2f MB", float64(kb)/1024)
// 	} else {
// 		return fmt.Sprintf("%.2f GB", float64(kb)/(1024*1024))
// 	}
// }

// func printProcessTree(proc *Process, level int, parentMem int64) {
// 	indent := strings.Repeat("  ", level)

// 	if level == 0 {
// 		// Root process
// 		fmt.Printf("%s%s (PID: %d): %s\n", indent, proc.Name, proc.PID, formatMemory(proc.MemoryUsage))
// 	} else {
// 		// Child process
// 		percentage := float64(proc.MemoryUsage) / float64(parentMem) * 100
// 		fmt.Printf("%s%s (PID: %d): %s (%.1f%% of parent)\n", indent, proc.Name, proc.PID, formatMemory(proc.MemoryUsage), percentage)
// 	}

// 	// Sort children by memory usage (descending)
// 	sort.Slice(proc.Children, func(i, j int) bool {
// 		return proc.Children[i].MemoryUsage > proc.Children[j].MemoryUsage
// 	})

// 	// Print children
// 	for _, child := range proc.Children {
// 		printProcessTree(child, level+1, proc.MemoryUsage)
// 	}
// }
