package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

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
		log.Printf("MYSQL_HOST not set, using default: %s", mysqlHost)
	}

	mysqlPort := os.Getenv("MYSQL_PORT")
	if mysqlPort == "" {
		mysqlPort = "3306"
		log.Printf("MYSQL_PORT not set, using default: %s", mysqlPort)
	}

	mysqlUser := os.Getenv("MYSQL_USER")
	if mysqlUser == "" {
		mysqlUser = "root"
		log.Printf("MYSQL_USER not set, using default: %s", mysqlUser)
	}

	mysqlPassword := os.Getenv("MYSQL_PASSWORD")
	if mysqlPassword == "" {
		log.Printf("MYSQL_PASSWORD not set, using empty password")
	}

	mysqlDBName := os.Getenv("MYSQL_DBNAME")
	if mysqlDBName == "" {
		mysqlDBName = "okblog"
		log.Printf("MYSQL_DBNAME not set, using default: %s", mysqlDBName)
	}

	// Format: username:password@tcp(host:port)/dbname
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		mysqlUser, mysqlPassword, mysqlHost, mysqlPort, mysqlDBName)

	var err error
	mysqlDB, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Printf("Failed to open MySQL connection: %v", err)
		return err
	}

	// Set connection pool settings
	mysqlDB.SetMaxOpenConns(25)
	mysqlDB.SetMaxIdleConns(5)
	mysqlDB.SetConnMaxLifetime(5 * time.Minute)

	// Test the connection
	if err := mysqlDB.Ping(); err != nil {
		log.Printf("Failed to ping MySQL: %v", err)
		return err
	}

	log.Println("Successfully connected to MySQL database")
	return nil
}

// GetMySQLDB returns the MySQL database connection
func GetMySQLDB() *sql.DB {
	if mysqlDB == nil {
		log.Fatal("MySQL connection not initialized. Call InitMySQLDB first.")
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
