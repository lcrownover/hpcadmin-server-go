package config

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"

	"gopkg.in/yaml.v3"
)

type ServerConfig struct {
	Host  string         `yaml:"host"`
	Port  int            `yaml:"port"`
	Oauth OauthConfig    `yaml:"oauth"`
	DB    DatabaseConfig `yaml:"database"`
}

type OauthConfig struct {
	TenantID     string `yaml:"tenant_id"`
	ClientID     string `yaml:"client_id"`
	ClientSecret string `yaml:"client_secret"`
}

type DatabaseConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
}

// Load loads the configuration from the given path
// If the path is empty, it will load the default configuration
// file from /etc/hpcadmin-server/config.yaml
func LoadFile(configPath string) (*ServerConfig, error) {
	var err error
	if configPath == "" {
		configPath = "/etc/hpcadmin-server/config.yaml"
	}
	slog.Debug("Configuration path found", "method", "Load", "path", configPath)

	slog.Debug("Reading config file", "method", "Load", "path", configPath)
	configData, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read configuration file: %v", err)
	}

	slog.Debug("Parsing YAML", "method", "Load", "path", configPath)
	cfg := &ServerConfig{}
	err = yaml.Unmarshal(configData, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to load configuration: %v", err)
	}

	return cfg, nil

}

func LoadEnvironment(cfg *ServerConfig) *ServerConfig {
	// HPCADMIN_SERVER_HOST
	if host, found := os.LookupEnv("HPCADMIN_SERVER_HOST"); found {
		slog.Debug("Found host override", "method", "LoadEnvironment", "host", host)
		cfg.Host = host
	}
	// HPCADMIN_SERVER_PORT
	if port, found := os.LookupEnv("HPCADMIN_SERVER_PORT"); found {
		slog.Debug("Found port override", "method", "LoadEnvironment", "port", port)
		iport, err := strconv.Atoi(port)
		if err != nil {
			slog.Warn("Invalid port number", "method", "LoadEnvironment", "port", port)
		} else {
			cfg.Port = iport
		}
	}
	// HPCADMIN_SERVER_DATABASE_HOST
	if dbhost, found := os.LookupEnv("HPCADMIN_SERVER_DATABASE_HOST"); found {
		slog.Debug("Found database host override", "method", "LoadEnvironment", "host", dbhost)
		cfg.DB.Host = dbhost
	}
	// HPCADMIN_SERVER_DATABASE_PORT
	if dbport, found := os.LookupEnv("HPCADMIN_SERVER_DATABASE_PORT"); found {
		slog.Debug("Found database port override", "method", "LoadEnvironment", "port", dbport)
		idbport, err := strconv.Atoi(dbport)
		if err != nil {
			slog.Warn("Invalid database port number", "method", "LoadEnvironment", "port", dbport)
		} else {
			cfg.DB.Port = idbport
		}
	}
	// HPCADMIN_SERVER_DATABASE_USER
	if dbuser, found := os.LookupEnv("HPCADMIN_SERVER_DATABASE_USER"); found {
		slog.Debug("Found database user override", "method", "LoadEnvironment", "user", dbuser)
		cfg.DB.User = dbuser
	}
	// HPCADMIN_SERVER_DATABASE_PASSWORD
	if dbpassword, found := os.LookupEnv("HPCADMIN_SERVER_DATABASE_USER"); found {
		slog.Debug("Found database user override", "method", "LoadEnvironment", "password", "REDACTED")
		cfg.DB.Password = dbpassword
	}
	// HPCADMIN_SERVER_DATABASE_DBNAME
	if dbname, found := os.LookupEnv("HPCADMIN_SERVER_DATABASE_DBNAME"); found {
		slog.Debug("Found database user override", "method", "LoadEnvironment", "dbname", dbname)
		cfg.DB.DBName = dbname
	}
	// HPCADMIN_SERVER_OAUTH_TENANT_ID
	if tenantID, found := os.LookupEnv("HPCADMIN_SERVER_OAUTH_TENANT_ID"); found {
		slog.Debug("Found oauth tenantID override", "method", "LoadEnvironment", "tenantID", tenantID)
		cfg.Oauth.TenantID = tenantID
	}
	// HPCADMIN_SERVER_OAUTH_CLIENT_ID
	if clientID, found := os.LookupEnv("HPCADMIN_SERVER_OAUTH_CLIENT_ID"); found {
		slog.Debug("Found oauth clientID override", "method", "LoadEnvironment", "clientID", clientID)
		cfg.Oauth.ClientID = clientID
	}
	// HPCADMIN_SERVER_OAUTH_CLIENT_SECRET
	if clientSecret, found := os.LookupEnv("HPCADMIN_SERVER_OAUTH_CLIENT_SECRET"); found {
		slog.Debug("Found oauth clientSecret override", "method", "LoadEnvironment", "clientSecret", "REDACTED")
		cfg.Oauth.ClientSecret = clientSecret
	}
	return cfg
}

func Validate(cfg *ServerConfig) error {
	if cfg.Host == "" {
		return fmt.Errorf("missing host")
	}
	if cfg.Port == 0 {
		return fmt.Errorf("missing port")
	}
	if cfg.DB.Host == "" {
		return fmt.Errorf("missing database host")
	}
	if cfg.DB.Port == 0 {
		return fmt.Errorf("missing database port")
	}
	if cfg.DB.User == "" {
		return fmt.Errorf("missing database user")
	}
	if cfg.DB.Password == "" {
		return fmt.Errorf("missing database password")
	}
	if cfg.DB.DBName == "" {
		return fmt.Errorf("missing database name")
	}
	if cfg.Oauth.TenantID == "" {
		return fmt.Errorf("missing oauth tenant ID")
	}
	if cfg.Oauth.ClientID == "" {
		return fmt.Errorf("missing oauth client ID")
	}
	if cfg.Oauth.ClientSecret == "" {
		return fmt.Errorf("missing oauth client secret")
	}
	return nil
}
