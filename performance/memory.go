package performance

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/mobile-next/mobilecli/devices"
	"github.com/mobile-next/mobilecli/types"
)

// CollectMemoryMetrics gathers memory usage information for a specific package
func CollectMemoryMetrics(device devices.Device, packageName string) (types.MemoryMetrics, error) {
	// Execute dumpsys meminfo for the package
	output, err := device.ExecuteCommand("shell", "dumpsys", "meminfo", packageName)
	if err != nil {
		return types.MemoryMetrics{}, fmt.Errorf("failed to execute dumpsys meminfo: %w", err)
	}

	metrics, err := parseDumpsysMeminfo(string(output))
	if err != nil {
		return types.MemoryMetrics{}, fmt.Errorf("failed to parse meminfo: %w", err)
	}

	return metrics, nil
}

// parseDumpsysMeminfo parses dumpsys meminfo output
// Example output:
// Applications Memory Usage (in Kilobytes):
// Uptime: 12345678 Realtime: 12345678
//
// ** MEMINFO in pid 12345 [com.example.app] **
//                    Pss  Private  Private  SwapPss     Heap     Heap     Heap
//                  Total    Dirty    Clean    Dirty     Size    Alloc     Free
//               --------  -------  -------  -------  -------  -------  -------
//   Native Heap    10264    10264        0       16   20480    15800     4680
//   Dalvik Heap     2136     2136        0        0    8192     3500     4692
// ...
// TOTAL PSS:      45678
func parseDumpsysMeminfo(output string) (types.MemoryMetrics, error) {
	metrics := types.MemoryMetrics{}
	lines := strings.Split(output, "\n")

	// Look for PSS (Proportional Set Size)
	pssRe := regexp.MustCompile(`TOTAL\s+PSS:\s+(\d+)`)
	// Look for heap information
	nativeHeapRe := regexp.MustCompile(`Native Heap\s+(\d+)`)
	dalvikHeapRe := regexp.MustCompile(`Dalvik Heap\s+(\d+)`)
	// Look for heap size/alloc in the header section
	heapSizeRe := regexp.MustCompile(`Heap\s+Size[:\s]+(\d+)`)
	heapAllocRe := regexp.MustCompile(`Heap\s+Alloc[:\s]+(\d+)`)

	for i, line := range lines {
		// Try to match TOTAL PSS
		if matches := pssRe.FindStringSubmatch(line); len(matches) > 1 {
			if pss, err := strconv.ParseInt(matches[1], 10, 64); err == nil {
				metrics.PSS = pss
			}
		}

		// Try to match Native Heap
		if matches := nativeHeapRe.FindStringSubmatch(line); len(matches) > 1 {
			if nativeHeap, err := strconv.ParseInt(matches[1], 10, 64); err == nil {
				metrics.NativeHeap = nativeHeap
			}
		}

		// Try to match Dalvik Heap and extract heap size/alloc from same line
		if strings.Contains(line, "Dalvik Heap") {
			fields := strings.Fields(line)
			// Format: "Dalvik Heap    2136     2136        0        0    8192     3500     4692"
			// Indices:      [0]       [1]      [2]       [3]      [4]     [5]      [6]      [7]
			if len(fields) >= 8 {
				// Field[5] is HeapSize, Field[6] is HeapAlloc
				if heapSize, err := strconv.ParseInt(fields[5], 10, 64); err == nil {
					metrics.HeapSize = heapSize
				}
				if heapAlloc, err := strconv.ParseInt(fields[6], 10, 64); err == nil {
					metrics.HeapAlloc = heapAlloc
				}
			}
		}

		// Alternative: Look for heap summary section
		if strings.Contains(line, "Heap") && i+1 < len(lines) {
			nextLine := lines[i+1]
			if matches := heapSizeRe.FindStringSubmatch(nextLine); len(matches) > 1 {
				if heapSize, err := strconv.ParseInt(matches[1], 10, 64); err == nil && metrics.HeapSize == 0 {
					metrics.HeapSize = heapSize
				}
			}
			if matches := heapAllocRe.FindStringSubmatch(nextLine); len(matches) > 1 {
				if heapAlloc, err := strconv.ParseInt(matches[1], 10, 64); err == nil && metrics.HeapAlloc == 0 {
					metrics.HeapAlloc = heapAlloc
				}
			}
		}
	}

	if metrics.PSS == 0 {
		return metrics, fmt.Errorf("failed to extract PSS from meminfo output")
	}

	return metrics, nil
}
