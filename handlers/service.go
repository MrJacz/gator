package handlers

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"text/template"
)

const systemdServiceTemplate = `[Unit]
Description=Gator RSS Feed Aggregator
After=network.target postgresql.service
Wants=postgresql.service

[Service]
Type=simple
User={{.User}}
WorkingDirectory={{.WorkingDir}}
ExecStart={{.ExecPath}} agg {{.Interval}} --concurrency={{.Concurrency}}
Restart=always
RestartSec=10
StandardOutput=journal
StandardError=journal
SyslogIdentifier=gator

[Install]
WantedBy=multi-user.target
`

type ServiceConfig struct {
	User        string
	WorkingDir  string
	ExecPath    string
	Interval    string
	Concurrency string
}

func serviceInstall(s *State, cmd Command) error {
	if runtime.GOOS != "linux" {
		return fmt.Errorf("service installation is only supported on Linux with systemd")
	}

	interval := "5m"
	concurrency := "3"

	if len(cmd.Args) > 0 {
		interval = cmd.Args[0]
	}
	if len(cmd.Args) > 1 {
		concurrency = cmd.Args[1]
	}

	// Get current user
	currentUser := os.Getenv("USER")
	if currentUser == "" {
		return fmt.Errorf("could not determine current user")
	}

	// Get executable path
	execPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("could not get executable path: %w", err)
	}
	execPath, err = filepath.Abs(execPath)
	if err != nil {
		return fmt.Errorf("could not get absolute path: %w", err)
	}

	// Get working directory
	workingDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("could not get working directory: %w", err)
	}

	config := ServiceConfig{
		User:        currentUser,
		WorkingDir:  workingDir,
		ExecPath:    execPath,
		Interval:    interval,
		Concurrency: concurrency,
	}

	// Parse and execute template
	tmpl, err := template.New("service").Parse(systemdServiceTemplate)
	if err != nil {
		return fmt.Errorf("could not parse service template: %w", err)
	}

	// Create temporary file for service
	tmpFile, err := os.CreateTemp("", "gator-*.service")
	if err != nil {
		return fmt.Errorf("could not create temporary file: %w", err)
	}
	defer os.Remove(tmpFile.Name())

	if err := tmpl.Execute(tmpFile, config); err != nil {
		tmpFile.Close()
		return fmt.Errorf("could not execute template: %w", err)
	}
	tmpFile.Close()

	// Define service file path
	servicePath := "/etc/systemd/system/gator.service"

	fmt.Printf("Installing gator service...\n")
	fmt.Printf("Configuration:\n")
	fmt.Printf("  User: %s\n", config.User)
	fmt.Printf("  Executable: %s\n", config.ExecPath)
	fmt.Printf("  Working Directory: %s\n", config.WorkingDir)
	fmt.Printf("  Interval: %s\n", config.Interval)
	fmt.Printf("  Concurrency: %s\n\n", config.Concurrency)

	// Copy service file (requires sudo)
	fmt.Printf("Installing service file to %s (requires sudo)...\n", servicePath)
	cpCmd := exec.Command("sudo", "cp", tmpFile.Name(), servicePath)
	cpCmd.Stdout = os.Stdout
	cpCmd.Stderr = os.Stderr
	if err := cpCmd.Run(); err != nil {
		return fmt.Errorf("failed to copy service file: %w", err)
	}

	// Reload systemd
	fmt.Println("Reloading systemd daemon...")
	reloadCmd := exec.Command("sudo", "systemctl", "daemon-reload")
	reloadCmd.Stdout = os.Stdout
	reloadCmd.Stderr = os.Stderr
	if err := reloadCmd.Run(); err != nil {
		return fmt.Errorf("failed to reload systemd: %w", err)
	}

	fmt.Println("\nService installed successfully!")
	fmt.Println("\nUseful Commands:")
	fmt.Println("  sudo systemctl start gator       # Start the service")
	fmt.Println("  sudo systemctl stop gator        # Stop the service")
	fmt.Println("  sudo systemctl status gator      # Check service status")
	fmt.Println("  sudo systemctl enable gator      # Enable service on boot")
	fmt.Println("  sudo systemctl disable gator     # Disable service on boot")
	fmt.Println("  sudo journalctl -u gator -f      # View service logs")
	fmt.Println("  gator service uninstall          # Uninstall the service")

	return nil
}

func serviceUninstall(s *State, cmd Command) error {
	if runtime.GOOS != "linux" {
		return fmt.Errorf("service management is only supported on Linux with systemd")
	}

	servicePath := "/etc/systemd/system/gator.service"

	fmt.Println("Uninstalling gator service...")

	// Stop service if running
	fmt.Println("Stopping service...")
	stopCmd := exec.Command("sudo", "systemctl", "stop", "gator")
	stopCmd.Stdout = os.Stdout
	stopCmd.Stderr = os.Stderr
	stopCmd.Run() // Ignore error if service not running

	// Disable service
	fmt.Println("Disabling service...")
	disableCmd := exec.Command("sudo", "systemctl", "disable", "gator")
	disableCmd.Stdout = os.Stdout
	disableCmd.Stderr = os.Stderr
	disableCmd.Run() // Ignore error if service not enabled

	// Remove service file
	fmt.Printf("Removing service file %s...\n", servicePath)
	rmCmd := exec.Command("sudo", "rm", "-f", servicePath)
	rmCmd.Stdout = os.Stdout
	rmCmd.Stderr = os.Stderr
	if err := rmCmd.Run(); err != nil {
		return fmt.Errorf("failed to remove service file: %w", err)
	}

	// Reload systemd
	fmt.Println("Reloading systemd daemon...")
	reloadCmd := exec.Command("sudo", "systemctl", "daemon-reload")
	reloadCmd.Stdout = os.Stdout
	reloadCmd.Stderr = os.Stderr
	if err := reloadCmd.Run(); err != nil {
		return fmt.Errorf("failed to reload systemd: %w", err)
	}

	fmt.Println("\nService uninstalled successfully!")
	return nil
}

func serviceStatus(s *State, cmd Command) error {
	if runtime.GOOS != "linux" {
		return fmt.Errorf("service management is only supported on Linux with systemd")
	}

	statusCmd := exec.Command("systemctl", "status", "gator")
	statusCmd.Stdout = os.Stdout
	statusCmd.Stderr = os.Stderr
	statusCmd.Run() // Ignore exit code, status Command returns non-zero when service is stopped

	return nil
}

func serviceLogs(s *State, cmd Command) error {
	if runtime.GOOS != "linux" {
		return fmt.Errorf("service management is only supported on Linux with systemd")
	}

	lines := "50"
	follow := false

	for _, arg := range cmd.Args {
		if arg == "-f" || arg == "--follow" {
			follow = true
		} else if strings.HasPrefix(arg, "-n") {
			lines = strings.TrimPrefix(arg, "-n")
		}
	}

	args := []string{"-u", "gator", "-n", lines}
	if follow {
		args = append(args, "-f")
	}

	logsCmd := exec.Command("journalctl", args...)
	logsCmd.Stdout = os.Stdout
	logsCmd.Stderr = os.Stderr
	logsCmd.Stdin = os.Stdin

	return logsCmd.Run()
}

func Service(s *State, cmd Command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("usage: %s <install|uninstall|status|logs> [args...]", cmd.Name)
	}

	subcommand := cmd.Args[0]
	subArgs := cmd.Args[1:]

	switch subcommand {
	case "install":
		return serviceInstall(s, Command{Name: "service install", Args: subArgs})
	case "uninstall":
		return serviceUninstall(s, Command{Name: "service uninstall", Args: subArgs})
	case "status":
		return serviceStatus(s, Command{Name: "service status", Args: subArgs})
	case "logs":
		return serviceLogs(s, Command{Name: "service logs", Args: subArgs})
	default:
		return fmt.Errorf("unknown subcommand: %s\nAvailable: install, uninstall, status, logs", subcommand)
	}
}
