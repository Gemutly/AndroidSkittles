package types

import "time"

// PerformanceMetrics represents comprehensive performance data for a device/app
type PerformanceMetrics struct {
	Timestamp   time.Time      `json:"timestamp"`
	DeviceID    string         `json:"deviceId"`
	PackageName string         `json:"packageName"`
	CPU         CPUMetrics     `json:"cpu"`
	Memory      MemoryMetrics  `json:"memory"`
	FPS         FPSMetrics     `json:"fps"`
	Battery     BatteryMetrics `json:"battery"`
	Network     NetworkMetrics `json:"network"`
}

// CPUMetrics represents CPU usage information
type CPUMetrics struct {
	Usage      float64 `json:"usage"`      // percentage
	UserTime   int64   `json:"userTime"`   // milliseconds
	SystemTime int64   `json:"systemTime"` // milliseconds
}

// MemoryMetrics represents memory usage information
type MemoryMetrics struct {
	PSS        int64 `json:"pss"`        // KB - Proportional Set Size
	HeapSize   int64 `json:"heapSize"`   // KB - Total heap size
	HeapAlloc  int64 `json:"heapAlloc"`  // KB - Allocated heap
	NativeHeap int64 `json:"nativeHeap"` // KB - Native heap
}

// FPSMetrics represents frame rate and rendering performance
type FPSMetrics struct {
	Current   float64 `json:"current"`   // Current FPS
	Average   float64 `json:"average"`   // Average FPS over sample period
	JankCount int     `json:"jankCount"` // Number of janky frames
	FrameTime []int64 `json:"frameTime"` // milliseconds per frame
}

// BatteryMetrics represents battery and power consumption
type BatteryMetrics struct {
	PowerMAh float64 `json:"powerMAh"` // Power consumption in mAh
	Percent  float64 `json:"percent"`  // Battery percentage
}

// NetworkMetrics represents network traffic statistics
type NetworkMetrics struct {
	RxBytes   int64 `json:"rxBytes"`   // Received bytes
	TxBytes   int64 `json:"txBytes"`   // Transmitted bytes
	RxPackets int64 `json:"rxPackets"` // Received packets
	TxPackets int64 `json:"txPackets"` // Transmitted packets
}
