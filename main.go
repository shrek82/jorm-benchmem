package main

import (
	"fmt"
	"log"

	"github.com/shrek82/jorm"
	"github.com/shrek82/jorm/driver/mysql"
)

type User struct {
	ID        int    `jorm:"column:id;primaryKey"`
	Name      string `jorm:"column:name"`
	Email     string `jorm:"column:email"`
	Age       int    `jorm:"column:age"`
	Active    bool   `jorm:"column:active"`
	CreatedAt string `jorm:"column:created_at"`
}

func main() {
	// DSN: username:password@tcp(host:port)/database?charset=utf8mb4&parseTime=True&loc=Local
	dsn := "root:@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=true&loc=Local"

	engine, err := jorm.Open(mysql.Open(dsn))
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	// Verify connection
	if sqlDB := engine.Connection(); sqlDB != nil {
		if err := sqlDB.Ping(); err != nil {
			log.Fatalf("Failed to ping database: %v", err)
		}
		defer sqlDB.Close()
	}

	var users []User
	// Query users
	err = jorm.Model(&users, engine).Table("user").Find(&users)
	if err != nil {
		// Try without Table if it infers from struct name "User" -> "users" or "user"
		// But let's keep Table("user") to be safe as per mysql.md
		log.Printf("Query failed: %v", err)
	}

	fmt.Printf("Found %d users:\n", len(users))
	for _, u := range users {
		fmt.Printf("%+v\n", u)
	}
}
