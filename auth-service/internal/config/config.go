package config

import (
	"fmt"
	"log"
	"time"

	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Config struct {
	DBHost         string `mapstructure:"DB_HOST"`
	DBPort         int    `mapstructure:"DB_PORT"`
	DBUser         string `mapstructure:"DB_USER"`
	DBPassword     string `mapstructure:"DB_PASSWORD"`
	DBName         string `mapstructure:"DB_NAME"`
	JWTIssuer      string `mapstructure:"JWT_ISSUER"`
	JWTAudience    string `mapstructure:"JWT_AUDIENCE"`
	JWTSigningKey  string `mapstructure:"JWT_SIGNING_KEY"`
	JWTValidity    int    `mapstructure:"JWT_VALIDITY"`
	AuthPort       int    `mapstructure:"AUTH_PORT"`
	ProbePort      int    `mapstructure:"PROBE_PORT"`
}

func (c *Config) AuthAddress() string {
	return fmt.Sprintf(":%d", c.AuthPort)
}

func (c *Config) ProbeAddress() string {
	return fmt.Sprintf(":%d", c.ProbePort)
}

func LoadConfig() (*Config, error) {
	viper.AddConfigPath(".")
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	// Set default values
	viper.SetDefault("AUTH_PORT", 9090)
	viper.SetDefault("PROBE_PORT", 9091)

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

func InitDB(cfg *Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	log.Println("Connected to MySQL database")
	return db, nil
}