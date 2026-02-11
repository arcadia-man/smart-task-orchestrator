package handlers

import (
	"net/http"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
)

type MonitoringHandlers struct{}

func NewMonitoringHandlers() *MonitoringHandlers {
	return &MonitoringHandlers{}
}

type SystemMetrics struct {
	CPU struct {
		Usage int `json:"usage"`
		Cores int `json:"cores"`
	} `json:"cpu"`
	Memory struct {
		Used       float64 `json:"used"`
		Total      float64 `json:"total"`
		Percentage int     `json:"percentage"`
	} `json:"memory"`
	Disk struct {
		Used       int `json:"used"`
		Total      int `json:"total"`
		Percentage int `json:"percentage"`
	} `json:"disk"`
	Network struct {
		Inbound  float64 `json:"inbound"`
		Outbound float64 `json:"outbound"`
	} `json:"network"`
}

type ServiceStatus struct {
	Name         string `json:"name"`
	Status       string `json:"status"`
	Uptime       string `json:"uptime"`
	ResponseTime string `json:"responseTime"`
}

type Alert struct {
	ID        string    `json:"id"`
	Type      string    `json:"type"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
	Resolved  bool      `json:"resolved"`
}

type MonitoringResponse struct {
	SystemMetrics SystemMetrics   `json:"systemMetrics"`
	Services      []ServiceStatus `json:"services"`
	Alerts        []Alert         `json:"alerts"`
}

func (h *MonitoringHandlers) GetSystemMetrics(c *gin.Context) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	// Convert bytes to MB for memory
	memUsedMB := float64(m.Alloc) / 1024 / 1024
	memTotalMB := float64(m.Sys) / 1024 / 1024
	memPercentage := int((memUsedMB / memTotalMB) * 100)

	metrics := SystemMetrics{
		CPU: struct {
			Usage int `json:"usage"`
			Cores int `json:"cores"`
		}{
			Usage: 45, // Mock CPU usage
			Cores: runtime.NumCPU(),
		},
		Memory: struct {
			Used       float64 `json:"used"`
			Total      float64 `json:"total"`
			Percentage int     `json:"percentage"`
		}{
			Used:       memUsedMB,
			Total:      memTotalMB,
			Percentage: memPercentage,
		},
		Disk: struct {
			Used       int `json:"used"`
			Total      int `json:"total"`
			Percentage int `json:"percentage"`
		}{
			Used:       120, // Mock disk usage in GB
			Total:      500,
			Percentage: 24,
		},
		Network: struct {
			Inbound  float64 `json:"inbound"`
			Outbound float64 `json:"outbound"`
		}{
			Inbound:  1.2, // Mock network usage in MB/s
			Outbound: 0.8,
		},
	}

	c.JSON(http.StatusOK, metrics)
}

func (h *MonitoringHandlers) GetServices(c *gin.Context) {
	services := []ServiceStatus{
		{
			Name:         "API Server",
			Status:       "healthy",
			Uptime:       "99.9%",
			ResponseTime: "45ms",
		},
		{
			Name:         "Database",
			Status:       "healthy",
			Uptime:       "99.8%",
			ResponseTime: "12ms",
		},
		{
			Name:         "Redis Cache",
			Status:       "healthy",
			Uptime:       "100%",
			ResponseTime: "3ms",
		},
		{
			Name:         "Scheduler Engine",
			Status:       "warning",
			Uptime:       "98.5%",
			ResponseTime: "120ms",
		},
		{
			Name:         "Log Aggregator",
			Status:       "healthy",
			Uptime:       "99.7%",
			ResponseTime: "8ms",
		},
	}

	c.JSON(http.StatusOK, services)
}

func (h *MonitoringHandlers) GetAlerts(c *gin.Context) {
	alerts := []Alert{
		{
			ID:        "1",
			Type:      "warning",
			Message:   "High CPU usage detected on scheduler engine",
			Timestamp: time.Now().Add(-25 * time.Minute),
			Resolved:  false,
		},
		{
			ID:        "2",
			Type:      "info",
			Message:   "Database backup completed successfully",
			Timestamp: time.Now().Add(-8 * time.Hour),
			Resolved:  true,
		},
		{
			ID:        "3",
			Type:      "error",
			Message:   "Failed to connect to external API",
			Timestamp: time.Now().Add(-18 * time.Hour),
			Resolved:  true,
		},
	}

	c.JSON(http.StatusOK, alerts)
}

func (h *MonitoringHandlers) GetFullMonitoring(c *gin.Context) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	// Convert bytes to MB for memory
	memUsedMB := float64(m.Alloc) / 1024 / 1024
	memTotalMB := float64(m.Sys) / 1024 / 1024
	memPercentage := int((memUsedMB / memTotalMB) * 100)

	response := MonitoringResponse{
		SystemMetrics: SystemMetrics{
			CPU: struct {
				Usage int `json:"usage"`
				Cores int `json:"cores"`
			}{
				Usage: 45,
				Cores: runtime.NumCPU(),
			},
			Memory: struct {
				Used       float64 `json:"used"`
				Total      float64 `json:"total"`
				Percentage int     `json:"percentage"`
			}{
				Used:       memUsedMB,
				Total:      memTotalMB,
				Percentage: memPercentage,
			},
			Disk: struct {
				Used       int `json:"used"`
				Total      int `json:"total"`
				Percentage int `json:"percentage"`
			}{
				Used:       120,
				Total:      500,
				Percentage: 24,
			},
			Network: struct {
				Inbound  float64 `json:"inbound"`
				Outbound float64 `json:"outbound"`
			}{
				Inbound:  1.2,
				Outbound: 0.8,
			},
		},
		Services: []ServiceStatus{
			{
				Name:         "API Server",
				Status:       "healthy",
				Uptime:       "99.9%",
				ResponseTime: "45ms",
			},
			{
				Name:         "Database",
				Status:       "healthy",
				Uptime:       "99.8%",
				ResponseTime: "12ms",
			},
			{
				Name:         "Redis Cache",
				Status:       "healthy",
				Uptime:       "100%",
				ResponseTime: "3ms",
			},
			{
				Name:         "Scheduler Engine",
				Status:       "warning",
				Uptime:       "98.5%",
				ResponseTime: "120ms",
			},
			{
				Name:         "Log Aggregator",
				Status:       "healthy",
				Uptime:       "99.7%",
				ResponseTime: "8ms",
			},
		},
		Alerts: []Alert{
			{
				ID:        "1",
				Type:      "warning",
				Message:   "High CPU usage detected on scheduler engine",
				Timestamp: time.Now().Add(-25 * time.Minute),
				Resolved:  false,
			},
			{
				ID:        "2",
				Type:      "info",
				Message:   "Database backup completed successfully",
				Timestamp: time.Now().Add(-8 * time.Hour),
				Resolved:  true,
			},
			{
				ID:        "3",
				Type:      "error",
				Message:   "Failed to connect to external API",
				Timestamp: time.Now().Add(-18 * time.Hour),
				Resolved:  true,
			},
		},
	}

	c.JSON(http.StatusOK, response)
}
