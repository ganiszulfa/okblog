package database

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	"okblog/tag/pkg/logger"

	_ "github.com/go-sql-driver/mysql"
)

var (
	mysqlDB *sql.DB
)

// InitMySQLDB initializes the MySQL database connection
func InitMySQLDB() error {
	mysqlHost := os.Getenv("MYSQL_HOST")
	if mysqlHost == "" {
		mysqlHost = "localhost"
		logger.Info("MYSQL_HOST not set, using default", map[string]string{"host": mysqlHost})
	}

	mysqlPort := os.Getenv("MYSQL_PORT")
	if mysqlPort == "" {
		mysqlPort = "3306"
		logger.Info("MYSQL_PORT not set, using default", map[string]string{"port": mysqlPort})
	}

	mysqlUser := os.Getenv("MYSQL_USER")
	if mysqlUser == "" {
		mysqlUser = "root"
		logger.Info("MYSQL_USER not set, using default", map[string]string{"user": mysqlUser})
	}

	mysqlPassword := os.Getenv("MYSQL_PASSWORD")
	if mysqlPassword == "" {
		logger.Warn("MYSQL_PASSWORD not set, using empty password", nil)
	}

	mysqlDBName := os.Getenv("MYSQL_DBNAME")
	if mysqlDBName == "" {
		mysqlDBName = "okblog"
		logger.Info("MYSQL_DBNAME not set, using default", map[string]string{"dbname": mysqlDBName})
	}

	// Format: username:password@tcp(host:port)/dbname
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		mysqlUser, mysqlPassword, mysqlHost, mysqlPort, mysqlDBName)

	var err error
	mysqlDB, err = sql.Open("mysql", dsn)
	if err != nil {
		logger.Error("Failed to open MySQL connection", err)
		return err
	}

	// Set connection pool settings
	mysqlDB.SetMaxOpenConns(25)
	mysqlDB.SetMaxIdleConns(5)
	mysqlDB.SetConnMaxLifetime(5 * time.Minute)

	// Test the connection
	if err := mysqlDB.Ping(); err != nil {
		logger.Error("Failed to ping MySQL", err)
		return err
	}

	logger.Info("Successfully connected to MySQL database", nil)
	return nil
}

// GetMySQLDB returns the MySQL database connection
func GetMySQLDB() *sql.DB {
	if mysqlDB == nil {
		logger.Fatal("MySQL connection not initialized. Call InitMySQLDB first.", nil)
	}
	return mysqlDB
}

// CloseMySQLDB closes the MySQL database connection
func CloseMySQLDB() error {
	if mysqlDB != nil {
		return mysqlDB.Close()
	}
	return nil
}
