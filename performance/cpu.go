package performance

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/mobile-next/mobilecli/devices"
	"github.com/mobile-next/mobilecli/types"
)

// CollectCPUMetrics gathers CPU usage information for a specific package
func CollectCPUMetrics(device devices.Device, packageName string) (types.CPUMetrics, error) {
	// Get PID for the package
	pidOutput, err := device.ExecuteCommand("shell", "pidof", packageName)
	if err != nil {
		return types.CPUMetrics{}, fmt.Errorf("failed to get PID: %w", err)
	}

	pid := strings.TrimSpace(string(pidOutput))
	if pid == "" {
		return types.CPUMetrics{}, fmt.Errorf("process not running for package: %s", packageName)
	}

	// Get CPU stats from /proc/[pid]/stat
	statOutput, err := device.ExecuteCommand("shell", "cat", fmt.Sprintf("/proc/%s/stat", pid))
	if err != nil {
		return types.CPUMetrics{}, fmt.Errorf("failed to read proc stat: %w", err)
	}

	cpuMetrics, err := parseProcStat(string(statOutput))
	if err != nil {
		return types.CPUMetrics{}, fmt.Errorf("failed to parse proc stat: %w", err)
	}

	// Get total CPU usage percentage using top
	topOutput, err := device.ExecuteCommand("shell", "top", "-b", "-n", "1", "-p", pid)
	if err == nil {
		// Try to extract CPU percentage from top output
		if usage := parseTopCPUUsage(string(topOutput)); usage >= 0 {
			cpuMetrics.Usage = usage
		}
	}

	return cpuMetrics, nil
}

// parseProcStat parses /proc/[pid]/stat format
// Format: pid (comm) state ... utime stime cutime cstime ...
func parseProcStat(stat string) (types.CPUMetrics, error) {
	// Remove process name (enclosed in parentheses) to simplify parsing
	re := regexp.MustCompile(`\(.*?\)`)
	stat = re.ReplaceAllString(stat, "")

	fields := strings.Fields(stat)
	if len(fields) < 15 {
		return types.CPUMetrics{}, fmt.Errorf("invalid /proc/stat format")
	}

	// utime is at index ~11, stime at ~12 (after removing process name)
	// These are in clock ticks, need to convert to milliseconds
	// Clock ticks per second is typically 100 on Android
	const clockTicksPerSecond = 100

	utime, err := strconv.ParseInt(fields[11], 10, 64)
	if err != nil {
		return types.CPUMetrics{}, fmt.Errorf("failed to parse utime: %w", err)
	}

	stime, err := strconv.ParseInt(fields[12], 10, 64)
	if err != nil {
		return types.CPUMetrics{}, fmt.Errorf("failed to parse stime: %w", err)
	}

	// Convert clock ticks to milliseconds
	userTimeMs := (utime * 1000) / clockTicksPerSecond
	systemTimeMs := (stime * 1000) / clockTicksPerSecond

	return types.CPUMetrics{
		Usage:      0, // Will be filled by top command
		UserTime:   userTimeMs,
		SystemTime: systemTimeMs,
	}, nil
}

// parseTopCPUUsage extracts CPU usage percentage from top output
func parseTopCPUUsage(topOutput string) float64 {
	// top output format varies, but typically contains CPU% column
	// Example: "12345 u0_a123   10 -10  ... 2.5% ... com.example.app"
	lines := strings.Split(topOutput, "\n")
	for _, line := range lines {
		if strings.Contains(line, "%") {
			// Look for percentage values
			re := regexp.MustCompile(`(\d+\.?\d*)%`)
			matches := re.FindStringSubmatch(line)
			if len(matches) > 1 {
				usage, err := strconv.ParseFloat(matches[1], 64)
				if err == nil {
					return usage
				}
			}
		}
	}
	return -1
}
