package main

import (
	"log"

	"github.com/shrek82/jorm"
	"github.com/shrek82/jorm/driver/sqlite"
)

type User struct {
	ID   int64  `jorm:"primaryKey;autoIncrement"`
	Name string `jorm:"column:username"`
	Age  int
}

func (u *User) TableName() string {
	return "users"
}

func main() {
	// SQLite 数据库文件路径
	dsn := "test.db"

	engine, err := jorm.Open(sqlite.Open(dsn))
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// 设置为全局默认数据库
	jorm.SetDefault(engine)

	// Verify connection
	if sqlDB := engine.Connection(); sqlDB != nil {
		if err := sqlDB.Ping(); err != nil {
			log.Fatalf("Failed to ping database: %v", err)
		}
		defer sqlDB.Close()
	}

	// 插入用户示例
	newUser := &User{
		Name: "张三",
		Age:  28,
	}

	err = jorm.Model(&User{}).Create(newUser)
	if err != nil {
		log.Printf("插入用户失败: %v", err)
	} else {
		log.Printf("成功插入用户，ID: %d", newUser.ID)
	}
}
