package main

import (
	"fmt"
	"log"
	"os/exec"
	"sync"
)

type Process struct {
	ID      int       `json:"id"`
	Command string    `json:"command"`
	Args    []string  `json:"args"`
	Cmd     *exec.Cmd `json:"-"`
	Running bool      `json:"running"`
	PID     int       `json:"pid"`
}

type ProcessManager struct {
	mu        sync.Mutex
	processes map[int]*Process
	counter   int
}

var processManager = &ProcessManager{
	processes: make(map[int]*Process),
}

func (pm *ProcessManager) StartProcess(command string, args ...string) (*Process, error) {
	// Get mutex lock and unlock after end of func or error
	pm.mu.Lock()
	defer pm.mu.Unlock()
	cmd := exec.Command(command, args...)
	if err := cmd.Start(); err != nil {
		return nil, err
	}

	process := &Process{
		ID:      pm.counter,
		Command: command,
		Args:    args,
		Cmd:     cmd,
		Running: true,
		PID:     cmd.Process.Pid,
	}
	pm.processes[pm.counter] = process
	pm.counter++
	go runProc(pm, process.ID)

	return process, nil
}

func (pm *ProcessManager) StopProcess(id int) error {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	process, exists := pm.processes[id]
	if !exists {
		return fmt.Errorf("Process with id %d not found", id)
	}
	if err := process.Cmd.Process.Kill(); err != nil {
		if err.Error() == "os: process already finished" {
			delete(pm.processes, id)
			return nil
		}
		return err
	}
	delete(pm.processes, id)
	return nil
}

func (pm *ProcessManager) ListProcesses() []*Process {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	var processes []*Process
	for _, process := range pm.processes {
		processes = append(processes, process)
	}
	if processes == nil {
		return make([]*Process, 0)
	}
	return processes
}

func (pm *ProcessManager) RestartProcess(id int) (*Process, error) {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	process, exists := pm.processes[id]
	if !exists {
		return nil, fmt.Errorf("Process with id %d not found", id)
	}
	if process.Running {
		return nil, fmt.Errorf("Process with id %d is already running", id)
	}
	cmd := exec.Command(process.Command, process.Args...)
	if err := cmd.Start(); err != nil {
		return nil, err
	}
	process.Cmd = cmd
	process.PID = cmd.Process.Pid
	process.Running = true
	go runProc(pm, process.ID)
	return process, nil
}

func runProc(pm *ProcessManager, pid int) {
	process := pm.processes[pid]
	err := process.Cmd.Wait()
	pm.mu.Lock()
	defer pm.mu.Unlock()
	if err != nil {
		log.Printf("Error: %v", err)
		process.Running = false
	}
}
