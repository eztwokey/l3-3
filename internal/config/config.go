package config

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/wb-go/wbf/config"
)

const (
	LocalPath = "./config/config.yaml"
)

type Config struct {
	Api      ApiCfg      `mapstructure:"api" yaml:"api"`
	Postgres PostgresCfg `mapstructure:"postgres" yaml:"postgres"`
}

type ApiCfg struct {
	GinMode      string `mapstructure:"gin_mode" yaml:"gin_mode"`
	Addr         string `mapstructure:"addr" yaml:"addr"`
	ReadTimeout  int    `mapstructure:"read_timeout" yaml:"read_timeout"`
	WriteTimeout int    `mapstructure:"write_timeout" yaml:"write_timeout"`
	IdleTimeout  int    `mapstructure:"idle_timeout" yaml:"idle_timeout"`
}

type PostgresCfg struct {
	Host     string `mapstructure:"host" yaml:"host"`
	Port     int    `mapstructure:"port" yaml:"port"`
	User     string `mapstructure:"user" yaml:"user"`
	Password string `mapstructure:"password" yaml:"password"`
	DBName   string `mapstructure:"dbname" yaml:"dbname"`
	SSLMode  string `mapstructure:"sslmode" yaml:"sslmode"`
}

func (p PostgresCfg) DSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		p.User, p.Password, p.Host, p.Port, p.DBName, p.SSLMode,
	)
}

func (cfg *Config) Read(paths ...string) error {
	c := config.New()

	if err := c.LoadConfigFiles(paths...); err != nil {
		return fmt.Errorf("{Read 1}: %w", err)
	}

	if err := c.Unmarshal(cfg); err != nil {
		return fmt.Errorf("{Read 2}: %w", err)
	}

	if err := cfg.validate(); err != nil {
		return fmt.Errorf("{Read 3}: %w", err)
	}

	cfg.print()
	return nil
}

func (cfg *Config) print() {
	log.Println("=================== CONFIG ===================")
	log.Println("Gin Mode...............", cfg.Api.GinMode)
	log.Println("API Address............", cfg.Api.Addr)
	log.Println("Read Timeout...........", cfg.Api.ReadTimeout)
	log.Println("Write Timeout..........", cfg.Api.WriteTimeout)
	log.Println("Idle Timeout...........", cfg.Api.IdleTimeout)
	log.Println("Postgres Host..........", cfg.Postgres.Host)
	log.Println("Postgres Port..........", cfg.Postgres.Port)
	log.Println("Postgres DB............", cfg.Postgres.DBName)
	log.Printf("==============================================\n\n")
}

func (cfg *Config) validate() error {
	if cfg.Api.GinMode == "" || (cfg.Api.GinMode != gin.DebugMode && cfg.Api.GinMode != gin.ReleaseMode) {
		return fmt.Errorf("invalid gin mode: %s", cfg.Api.GinMode)
	}
	if cfg.Api.Addr == "" {
		return fmt.Errorf("invalid addr: %s", cfg.Api.Addr)
	}
	if cfg.Api.ReadTimeout <= 0 {
		return fmt.Errorf("invalid read timeout: %d", cfg.Api.ReadTimeout)
	}
	if cfg.Api.WriteTimeout <= 0 {
		return fmt.Errorf("invalid write timeout: %d", cfg.Api.WriteTimeout)
	}
	if cfg.Api.IdleTimeout <= 0 {
		return fmt.Errorf("invalid idle timeout: %d", cfg.Api.IdleTimeout)
	}
	if cfg.Postgres.Host == "" {
		return fmt.Errorf("postgres host is required")
	}
	if cfg.Postgres.Port <= 0 {
		return fmt.Errorf("postgres port is required")
	}
	if cfg.Postgres.DBName == "" {
		return fmt.Errorf("postgres dbname is required")
	}
	return nil
}
