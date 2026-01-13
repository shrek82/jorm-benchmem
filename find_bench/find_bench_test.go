package find_bench

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

// BenchmarkJormFindByID 测试 jorm Find 查询单条记录的 QPS
func BenchmarkJormFindByID(b *testing.B) {
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
			b.Fatalf("jorm find: %v", err)
		}
	}
}

// BenchmarkGormFindByID 测试 gorm Find 查询单条记录的 QPS
func BenchmarkGormFindByID(b *testing.B) {
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
			b.Fatalf("gorm find: %v", err)
		}
	}
}

// BenchmarkXormFindByID 测试 xorm Find 查询单条记录的 QPS
func BenchmarkXormFindByID(b *testing.B) {
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
		has, err := engine.ID(queryID).Get(&user)
		if err != nil {
			b.Fatalf("xorm find: %v", err)
		}
		if !has {
			b.Fatalf("user not found: %d", queryID)
		}
	}
}

// BenchmarkJormFindLimit 测试 jorm Find 查询限制数量记录的 QPS
func BenchmarkJormFindLimit(b *testing.B) {
	// 准备 5000 条测试数据
	setupTestData(b, 5000)

	engine, err := NewJormEngine()
	if err != nil {
		b.Fatalf("new jorm engine: %v", err)
	}
	if sqlDB := engine.Connection(); sqlDB != nil {
		defer sqlDB.Close()
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var users []User
		err := jorm.Model(&User{}, engine).Limit(100).FindAll(&users)
		if err != nil {
			b.Fatalf("jorm find: %v", err)
		}
		if len(users) != 100 {
			b.Fatalf("expected 100 users, got %d", len(users))
		}
	}
}

// BenchmarkGormFindLimit 测试 gorm Find 查询限制数量记录的 QPS
func BenchmarkGormFindLimit(b *testing.B) {
	setupTestData(b, 5000)

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
		var users []User
		if err := db.Limit(100).Find(&users).Error; err != nil {
			b.Fatalf("gorm find: %v", err)
		}
		if len(users) != 100 {
			b.Fatalf("expected 100 users, got %d", len(users))
		}
	}
}

// BenchmarkXormFindLimit 测试 xorm Find 查询限制数量记录的 QPS
func BenchmarkXormFindLimit(b *testing.B) {
	setupTestData(b, 5000)

	engine, err := NewXormEngine()
	if err != nil {
		b.Fatalf("new xorm engine: %v", err)
	}
	defer engine.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var users []User
		if err := engine.Limit(100).Find(&users); err != nil {
			b.Fatalf("xorm find: %v", err)
		}
		if len(users) != 100 {
			b.Fatalf("expected 100 users, got %d", len(users))
		}
	}
}

// BenchmarkJormFindAll 测试 jorm FindAll 查询所有记录的 QPS
func BenchmarkJormFindAll(b *testing.B) {
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
		var users []User
		err := jorm.Model(&User{}, engine).FindAll(&users)
		if err != nil {
			b.Fatalf("jorm find: %v", err)
		}
	}
}

// BenchmarkGormFindAll 测试 gorm Find 查询所有记录的 QPS
func BenchmarkGormFindAll(b *testing.B) {
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
		var users []User
		if err := db.Find(&users).Error; err != nil {
			b.Fatalf("gorm find: %v", err)
		}
	}
}

// BenchmarkXormFindAll 测试 xorm Find 查询所有记录的 QPS
func BenchmarkXormFindAll(b *testing.B) {
	setupTestData(b, 1000)

	engine, err := NewXormEngine()
	if err != nil {
		b.Fatalf("new xorm engine: %v", err)
	}
	defer engine.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var users []User
		if err := engine.Find(&users); err != nil {
			b.Fatalf("xorm find: %v", err)
		}
	}
}
