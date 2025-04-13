package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Config struct {
	DBHost         string
	DBPort         int
	DBUser         string
	DBPassword     string
	DBName         string
	JWTIssuer      string
	JWTAudience    string
	JWTSigningKey  string
	JWTValidity    int
	APIPort        int
	ProbePort      int
}

func (c *Config) APIAddress() string {
	return fmt.Sprintf(":%d", c.APIPort)
}

func (c *Config) ProbeAddress() string {
	return fmt.Sprintf(":%d", c.ProbePort)
}

func LoadConfig() *Config {

	viper.SetConfigFile(".env")
	_ = viper.ReadInConfig() 

	viper.AutomaticEnv()

	getInt := func(key string, defaultValue int) int {
		if value := viper.GetString(key); value != "" {
			if intVal, err := strconv.Atoi(value); err == nil {
				return intVal
			}
		}
		return defaultValue
	}

	return &Config{
		DBHost:        getEnv("DB_HOST", "localhost"),
		DBPort:        getInt("DB_PORT", 3306),
		DBUser:        getEnv("DB_USER", "root"),
		DBPassword:    getEnv("DB_PASSWORD", ""),
		DBName:        getEnv("DB_NAME", "taskmango_db"),
		JWTIssuer:     getEnv("JWT_ISSUER", "auth-service"),
		JWTAudience:   getEnv("JWT_AUDIENCE", "task-manager"),
		JWTSigningKey: getEnv("JWT_SIGNING_KEY", "secret-key-change-me"),
		JWTValidity:   getInt("JWT_VALIDITY", 3600),
		APIPort:       getInt("API_PORT", 8080),
		ProbePort:     getInt("PROBE_PORT", 8081),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	if value := viper.GetString(key); value != "" {
		return value
	}
	return defaultValue
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