package main

import (
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

func (u *User) TableName() string {
	return "users"
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

	// 插入用户示例
	newUser := User{
		Name:      "张三",
		Email:     "zhangsan@example.com",
		Age:       28,
		Active:    true,
		CreatedAt: "2025-06-25 15:04:05",
	}

	err = jorm.Model(&User{}, engine).Create(&newUser)
	if err != nil {
		log.Printf("插入用户失败: %v", err)
	} else {
		log.Printf("成功插入用户，ID: %d", newUser.ID)
	}

}
