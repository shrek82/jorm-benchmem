package query_bench

import (
	"fmt"
	"testing"

	"github.com/shrek82/jorm"
)

// setupTestData 在每个 benchmark 开始前准备测试数据
func setupTestData(b *testing.B, count int) {
	b.Helper()
	sqlDB, err := NewSQLDB()
	if err != nil {
		b.Fatalf("open sql db: %v", err)
	}
	defer sqlDB.Close()

	// 清空表
	if _, err := sqlDB.Exec("DELETE FROM users"); err != nil {
		b.Fatalf("delete users: %v", err)
	}
	if _, err := sqlDB.Exec("DELETE FROM sqlite_sequence WHERE name='users'"); err != nil {
		// 忽略错误
	}

	// 插入测试数据
	stmt, err := sqlDB.Prepare("INSERT INTO users (username, age) VALUES (?, ?)")
	if err != nil {
		b.Fatalf("prepare insert: %v", err)
	}
	defer stmt.Close()

	for i := 1; i <= count; i++ {
		_, err := stmt.Exec(fmt.Sprintf("user_%d", i), 20+i%30)
		if err != nil {
			b.Fatalf("insert data: %v", err)
		}
	}
}

// BenchmarkJormQueryByID 测试 jorm 根据 ID 查询单条记录
func BenchmarkJormQueryByID(b *testing.B) {
	// 准备 1000 条测试数据
	setupTestData(b, 1000)

	engine, err := NewJormEngine()
	if err != nil {
		b.Fatalf("new jorm engine: %v", err)
	}
	if sqlDB := engine.Connection(); sqlDB != nil {
		defer sqlDB.Close()
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var user User
		// 查询 ID 在 1-1000 之间的随机记录
		queryID := int64(i%1000 + 1)
		err := jorm.Model(&User{}, engine).Where("id = ?", queryID).Find(&user)
		if err != nil {
			b.Fatalf("jorm query: %v", err)
		}
	}
}

// BenchmarkGormQueryByID 测试 gorm 根据 ID 查询单条记录
func BenchmarkGormQueryByID(b *testing.B) {
	setupTestData(b, 1000)

	db, err := NewGormDB()
	if err != nil {
		b.Fatalf("new gorm db: %v", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		b.Fatalf("gorm db.DB(): %v", err)
	}
	defer sqlDB.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var user User
		queryID := int64(i%1000 + 1)
		if err := db.Where("id = ?", queryID).First(&user).Error; err != nil {
			b.Fatalf("gorm query: %v", err)
		}
	}
}

// BenchmarkXormQueryByID 测试 xorm 根据 ID 查询单条记录
func BenchmarkXormQueryByID(b *testing.B) {
	setupTestData(b, 1000)

	engine, err := NewXormEngine()
	if err != nil {
		b.Fatalf("new xorm engine: %v", err)
	}
	defer engine.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var user User
		queryID := int64(i%1000 + 1)
		has, err := engine.Where("id = ?", queryID).Get(&user)
		if err != nil {
			b.Fatalf("xorm query: %v", err)
		}
		if !has {
			b.Fatalf("user not found: %d", queryID)
		}
	}
}

// BenchmarkJormQueryByName 测试 jorm 根据 Name 查询单条记录
func BenchmarkJormQueryByName(b *testing.B) {
	setupTestData(b, 1000)

	engine, err := NewJormEngine()
	if err != nil {
		b.Fatalf("new jorm engine: %v", err)
	}
	if sqlDB := engine.Connection(); sqlDB != nil {
		defer sqlDB.Close()
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var user User
		queryName := fmt.Sprintf("user_%d", i%1000+1)
		err := jorm.Model(&User{}, engine).Where("username = ?", queryName).Find(&user)
		if err != nil {
			b.Fatalf("jorm query: %v", err)
		}
	}
}

// BenchmarkGormQueryByName 测试 gorm 根据 Name 查询单条记录
func BenchmarkGormQueryByName(b *testing.B) {
	setupTestData(b, 1000)

	db, err := NewGormDB()
	if err != nil {
		b.Fatalf("new gorm db: %v", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		b.Fatalf("gorm db.DB(): %v", err)
	}
	defer sqlDB.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var user User
		queryName := fmt.Sprintf("user_%d", i%1000+1)
		if err := db.Where("username = ?", queryName).First(&user).Error; err != nil {
			b.Fatalf("gorm query: %v", err)
		}
	}
}

// BenchmarkXormQueryByName 测试 xorm 根据 Name 查询单条记录
func BenchmarkXormQueryByName(b *testing.B) {
	setupTestData(b, 1000)

	engine, err := NewXormEngine()
	if err != nil {
		b.Fatalf("new xorm engine: %v", err)
	}
	defer engine.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var user User
		queryName := fmt.Sprintf("user_%d", i%1000+1)
		has, err := engine.Where("username = ?", queryName).Get(&user)
		if err != nil {
			b.Fatalf("xorm query: %v", err)
		}
		if !has {
			b.Fatalf("user not found: %s", queryName)
		}
	}
}
