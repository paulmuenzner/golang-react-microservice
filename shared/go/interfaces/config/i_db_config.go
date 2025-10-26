package interfaceconfig

import "time"

// Parameter structure for database configuration
type IPostgreConnectionConfig struct {
	Host            string
	Port            int
	User            string
	Password        string
	DBName          string
	SSLMode         string
	Timeout         time.Duration
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxIdleTime time.Duration
	ConnMaxLifetime time.Duration
}
