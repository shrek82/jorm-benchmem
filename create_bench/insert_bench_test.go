package create_bench

import (
	"fmt"
	"testing"

	"github.com/shrek82/jorm"
)

// prepareUser 生成一条测试数据，index 用于避免完全相同的数据
func prepareUser(index int) *User {
	return &User{
		Name: fmt.Sprintf("user_%d", index),
		Age:  20 + index%30,
	}
}

// truncateUsers 在每个 benchmark 开始前清空 users 表，保证数据量一致
func truncateUsers(b *testing.B) {
	b.Helper()
	sqlDB, err := NewSQLDB()
	if err != nil {
		b.Fatalf("open sql db: %v", err)
	}
	defer sqlDB.Close()

	// SQLite 使用 DELETE 代替 TRUNCATE
	if _, err := sqlDB.Exec("DELETE FROM users"); err != nil {
		b.Fatalf("delete users: %v", err)
	}
	// 重置自增 ID
	if _, err := sqlDB.Exec("DELETE FROM sqlite_sequence WHERE name='users'"); err != nil {
		// 忽略错误，因为表可能还没有插入过数据
	}
}

func BenchmarkJormInsert(b *testing.B) {
	truncateUsers(b)

	engine, err := NewJormEngine()
	if err != nil {
		b.Fatalf("new jorm engine: %v", err)
	}
	// 关闭底层连接
	if sqlDB := engine.Connection(); sqlDB != nil {
		defer sqlDB.Close()
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		u := prepareUser(i)
		if err := jorm.Model(&User{}, engine).Create(u); err != nil {
			b.Fatalf("jorm insert: %v", err)
		}
	}
}

func BenchmarkGormInsert(b *testing.B) {
	truncateUsers(b)

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
		u := prepareUser(i)
		if err := db.Create(u).Error; err != nil {
			b.Fatalf("gorm insert: %v", err)
		}
	}
}

func BenchmarkXormInsert(b *testing.B) {
	truncateUsers(b)

	engine, err := NewXormEngine()
	if err != nil {
		b.Fatalf("new xorm engine: %v", err)
	}
	defer engine.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		u := prepareUser(i)
		if _, err := engine.Insert(u); err != nil {
			b.Fatalf("xorm insert: %v", err)
		}
	}
}
