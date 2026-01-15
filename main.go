package main

import (
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/shrek82/jorm"
	"github.com/shrek82/jorm/logger"
)

type User struct {
	ID        int64     `jorm:"pk;auto"`
	Name      string    `jorm:"column:username"`
	Age       int       `jorm:"column:age"`
	CreatedAt time.Time `jorm:"auto_time"`
	UpdatedAt time.Time `jorm:"auto_update"`
}

func (u *User) TableName() string {
	return "users"
}

func main() {
	// SQLite 数据库文件路径
	//dsn := "test.db"

	engine, err := jorm.Open("mysql", "root@tcp(127.0.0.1:3306)/jorm?charset=utf8mb4", nil)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	customLogger := logger.NewStdLogger()
	customLogger.SetLevel(logger.LogLevelInfo)
	customLogger.SetFormat(logger.LogFormatText)
	engine.SetLogger(customLogger)

	defer engine.Close()

	// 1. 自动迁移 (创建表)
	err = engine.AutoMigrate(&User{})
	if err != nil {
		log.Fatalf("AutoMigrate failed: %v", err)
	}

	//删除所有数据
	engine.Model(&User{}).Where("id > ?", 0).Delete()

	// 2. 插入用户示例
	newUser := &User{
		Name: "张三",
		Age:  28,
	}

	_, err = engine.Model(newUser).Insert(newUser)
	if err != nil {
		log.Printf("插入用户失败: %v", err)
	} else {
		log.Printf("成功插入用户，ID: %d", newUser.ID)
	}

	// 3. 查询用户示例
	var user User
	err = engine.Model(&User{}).Where("id = ?", newUser.ID).First(&user)
	if err != nil {
		log.Printf("查询用户失败: %v", err)
	} else {
		log.Printf("查询到用户: %+v", user)
	}

	// 4. 更新用户示例
	user.Age = 29
	rows, err := engine.Model(&user).Where("id = ?", user.ID).Update(user)
	if err != nil {
		log.Printf("更新用户失败: %v", err)
	} else {
		log.Printf("成功更新用户，影响行数: %d", rows)
	}

	// 5. 删除用户示例
	// rows, err = engine.Model(&User{}).Where("id = ?", user.ID).Delete()
	// if err != nil {
	// 	log.Printf("删除用户失败: %v", err)
	// } else {
	// 	log.Printf("成功删除用户，影响行数: %d", rows)
	// }
}
