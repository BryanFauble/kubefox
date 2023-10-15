package config

import "time"

var (
	Instance  string
	Platform  string
	Namespace string

	CertDir string

	NumWorkers        int
	EventTTL          time.Duration
	TelemetryInterval time.Duration

	LogFormat string
	LogLevel  string

	GRPCSrvAddr   string
	HTTPSrvAddr   string
	HTTPSSrvAddr  string
	HealthSrvAddr string

	VaultAddr     string
	NATSAddr      string
	TelemetryAddr string
)