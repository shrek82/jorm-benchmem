package multi_query_bench

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

// BenchmarkJormQueryLimit 测试 jorm 查询多条记录（LIMIT）
func BenchmarkJormQueryLimit(b *testing.B) {
	// 准备 10000 条测试数据
	setupTestData(b, 10000)

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
		limit := 100
		err := jorm.Model(&User{}, engine).Limit(limit).FindAll(&users)
		if err != nil {
			b.Fatalf("jorm query: %v", err)
		}
		if len(users) != limit {
			b.Fatalf("expected %d users, got %d", limit, len(users))
		}
	}
}

// BenchmarkGormQueryLimit 测试 gorm 查询多条记录（LIMIT）
func BenchmarkGormQueryLimit(b *testing.B) {
	setupTestData(b, 10000)

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
		limit := 100
		if err := db.Limit(limit).Find(&users).Error; err != nil {
			b.Fatalf("gorm query: %v", err)
		}
		if len(users) != limit {
			b.Fatalf("expected %d users, got %d", limit, len(users))
		}
	}
}

// BenchmarkXormQueryLimit 测试 xorm 查询多条记录（LIMIT）
func BenchmarkXormQueryLimit(b *testing.B) {
	setupTestData(b, 10000)

	engine, err := NewXormEngine()
	if err != nil {
		b.Fatalf("new xorm engine: %v", err)
	}
	defer engine.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var users []User
		limit := 100
		if err := engine.Limit(limit).Find(&users); err != nil {
			b.Fatalf("xorm query: %v", err)
		}
		if len(users) != limit {
			b.Fatalf("expected %d users, got %d", limit, len(users))
		}
	}
}

// BenchmarkJormQueryByAgeRange 测试 jorm 查询年龄范围内的多条记录
func BenchmarkJormQueryByAgeRange(b *testing.B) {
	// 准备 10000 条测试数据
	setupTestData(b, 10000)

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
		ageMin := 20 + i%10
		ageMax := ageMin + 10
		err := jorm.Model(&User{}, engine).Where("age BETWEEN ? AND ?", ageMin, ageMax).FindAll(&users)
		if err != nil {
			b.Fatalf("jorm query: %v", err)
		}
	}
}

// BenchmarkGormQueryByAgeRange 测试 gorm 查询年龄范围内的多条记录
func BenchmarkGormQueryByAgeRange(b *testing.B) {
	setupTestData(b, 10000)

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
		ageMin := 20 + i%10
		ageMax := ageMin + 10
		if err := db.Where("age BETWEEN ? AND ?", ageMin, ageMax).Find(&users).Error; err != nil {
			b.Fatalf("gorm query: %v", err)
		}
	}
}

// BenchmarkXormQueryByAgeRange 测试 xorm 查询年龄范围内的多条记录
func BenchmarkXormQueryByAgeRange(b *testing.B) {
	setupTestData(b, 10000)

	engine, err := NewXormEngine()
	if err != nil {
		b.Fatalf("new xorm engine: %v", err)
	}
	defer engine.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var users []User
		ageMin := 20 + i%10
		ageMax := ageMin + 10
		if err := engine.Where("age BETWEEN ? AND ?", ageMin, ageMax).Find(&users); err != nil {
			b.Fatalf("xorm query: %v", err)
		}
	}
}

// BenchmarkJormQueryAll 测试 jorm 查询所有记录
func BenchmarkJormQueryAll(b *testing.B) {
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
		err := jorm.Model(&User{}, engine).FindAll(&users)
		if err != nil {
			b.Fatalf("jorm query: %v", err)
		}
	}
}

// BenchmarkGormQueryAll 测试 gorm 查询所有记录
func BenchmarkGormQueryAll(b *testing.B) {
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
		if err := db.Find(&users).Error; err != nil {
			b.Fatalf("gorm query: %v", err)
		}
	}
}

// BenchmarkXormQueryAll 测试 xorm 查询所有记录
func BenchmarkXormQueryAll(b *testing.B) {
	setupTestData(b, 5000)

	engine, err := NewXormEngine()
	if err != nil {
		b.Fatalf("new xorm engine: %v", err)
	}
	defer engine.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var users []User
		if err := engine.Find(&users); err != nil {
			b.Fatalf("xorm query: %v", err)
		}
	}
}
