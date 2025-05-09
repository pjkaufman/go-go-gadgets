package proc

type Process struct {
	PID         int
	PPID        int
	Name        string
	IsRoot      bool
	MemoryUsage int64 // in KB
	Children    []*Process
}
