package cmd

import (
	"fmt"
	"sort"
	"strings"

	"github.com/pjkaufman/go-go-gadgets/mem-analyzer/internal/proc"
	"github.com/pjkaufman/go-go-gadgets/pkg/logger"
	"github.com/spf13/cobra"
)

var TreeCmd = &cobra.Command{
	Use:   "tree",
	Short: "Lists the processes and their RAM usage as a kind of tree where children are indented under their parents",
	// Example: heredoc.Doc(`To show a list of all series names that are being tracked:
	// magnum list

	// To include information like publisher, status, series, etc.:
	// magnum list -v
	// `),
	Run: func(cmd *cobra.Command, args []string) {
		const minimumKb int64 = 0
		processes, numIgnoredProcesses, err := proc.GetAllProcesses(minimumKb)
		if err != nil {
			logger.WriteErrorf("Error getting processes: %v\n", err)
		}

		// // Create a map for easy lookup
		processMap := make(map[int]*proc.Process)
		for i := range processes {
			processMap[processes[i].PID] = processes[i]
		}

		// Build process tree
		var rootProcesses []*proc.Process
		for _, process := range processes {
			if process.IsRoot {
				rootProcesses = append(rootProcesses, process)
			}
		}

		// Calculate memory usage including children
		for _, proc := range rootProcesses {
			calculateTotalMemory(proc)
		}

		// Print the process tree with memory usage
		for _, proc := range rootProcesses {
			printProcessTree(proc, 0, proc.MemoryUsage)
		}

		if numIgnoredProcesses != 0 {
			logger.WriteInfof("%d process(es) were ignored because they had less than or equal to %d KB of RAM usage\n", numIgnoredProcesses, minimumKb)
		}
	},
}

func init() {
	rootCmd.AddCommand(TreeCmd)

	// TreeCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "show the publisher and other info about the series")
	// TreeCmd.Flags().StringVarP(&seriesPublisher, "publisher", "p", "", "show series with the specified publisher")
	// TreeCmd.Flags().StringVarP(&seriesType, "type", "t", "", "show series with the specified type")
	// TreeCmd.Flags().StringVarP(&seriesStatus, "status", "r", "", "show series with the specified status")
}

func calculateTotalMemory(proc *proc.Process) int64 {
	totalMem := proc.MemoryUsage

	for _, child := range proc.Children {
		totalMem += calculateTotalMemory(child)
	}

	proc.MemoryUsage = totalMem
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

func printProcessTree(proc *proc.Process, level int, parentMem int64) {
	indent := strings.Repeat("  ", level)

	if level == 0 {
		// Root process
		fmt.Printf("%s%s (PID: %d): %s\n", indent, proc.Name, proc.PID, formatMemory(proc.MemoryUsage))
	} else {
		// Child process
		percentage := 0.0
		if proc.MemoryUsage != 0 && parentMem != 0 {
			percentage = float64(proc.MemoryUsage) / float64(parentMem) * 100
		}

		fmt.Printf("%s%s (PID: %d): %s (%.1f%% of parent)\n", indent, proc.Name, proc.PID, formatMemory(proc.MemoryUsage), percentage)
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
