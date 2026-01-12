package update_bench

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

// BenchmarkJormUpdateByID 测试 jorm 根据 ID 更新单条记录
func BenchmarkJormUpdateByID(b *testing.B) {
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
		// 更新 ID 在 1-1000 之间的随机记录
		queryID := int64(i%1000 + 1)
		updatedUser := User{
			Name: fmt.Sprintf("updated_user_%d", i),
			Age:  30 + i%20,
		}

		// 使用 Where 条件进行更新
		result, err := jorm.Model(&User{}, engine).Where("id = ?", queryID).Update(updatedUser)
		if err != nil {
			b.Fatalf("jorm update: %v", err)
		}
		if result == 0 {
			b.Fatalf("jorm update affected 0 rows")
		}
	}
}

// BenchmarkGormUpdateByID 测试 gorm 根据 ID 更新单条记录
func BenchmarkGormUpdateByID(b *testing.B) {
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
		queryID := int64(i%1000 + 1)
		updatedUser := User{
			Name: fmt.Sprintf("updated_user_%d", i),
			Age:  30 + i%20,
		}

		result := db.Where("id = ?", queryID).Updates(&updatedUser)
		if result.Error != nil {
			b.Fatalf("gorm update: %v", result.Error)
		}
		if result.RowsAffected == 0 {
			b.Fatalf("gorm update affected 0 rows")
		}
	}
}

// BenchmarkXormUpdateByID 测试 xorm 根据 ID 更新单条记录
func BenchmarkXormUpdateByID(b *testing.B) {
	setupTestData(b, 1000)

	engine, err := NewXormEngine()
	if err != nil {
		b.Fatalf("new xorm engine: %v", err)
	}
	defer engine.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		queryID := int64(i%1000 + 1)
		updatedUser := User{
			ID:   queryID, // Xorm 需要 ID 来定位记录
			Name: fmt.Sprintf("updated_user_%d", i),
			Age:  30 + i%20,
		}

		affected, err := engine.ID(queryID).Update(&updatedUser)
		if err != nil {
			b.Fatalf("xorm update: %v", err)
		}
		if affected == 0 {
			b.Fatalf("xorm update affected 0 rows")
		}
	}
}

// BenchmarkJormUpdateByCondition 测试 jorm 根据条件更新多条记录
func BenchmarkJormUpdateByCondition(b *testing.B) {
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
		// 更新年龄在某个范围内的用户
		ageMin := 20 + i%10
		ageMax := ageMin + 5

		updatedUser := User{
			Age: 30 + i%20,
		}

		result, err := jorm.Model(&User{}, engine).Where("age BETWEEN ? AND ?", ageMin, ageMax).Update(updatedUser)
		if err != nil {
			b.Fatalf("jorm update: %v", err)
		}
		if result == 0 {
			// 允许影响 0 行的情况，因为可能没有匹配的记录
		}
	}
}

// BenchmarkGormUpdateByCondition 测试 gorm 根据条件更新多条记录
func BenchmarkGormUpdateByCondition(b *testing.B) {
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
		ageMin := 20 + i%10
		ageMax := ageMin + 5

		updatedUser := User{
			Age: 30 + i%20,
		}

		result := db.Where("age BETWEEN ? AND ?", ageMin, ageMax).Updates(&updatedUser)
		if result.Error != nil {
			b.Fatalf("gorm update: %v", result.Error)
		}
		if result.RowsAffected == 0 {
			// 允许影响 0 行的情况，因为可能没有匹配的记录
		}
	}
}

// BenchmarkXormUpdateByCondition 测试 xorm 根据条件更新多条记录
func BenchmarkXormUpdateByCondition(b *testing.B) {
	setupTestData(b, 5000)

	engine, err := NewXormEngine()
	if err != nil {
		b.Fatalf("new xorm engine: %v", err)
	}
	defer engine.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ageMin := 20 + i%10
		ageMax := ageMin + 5

		updatedUser := User{
			Age: 30 + i%20,
		}

		affected, err := engine.Where("age BETWEEN ? AND ?", ageMin, ageMax).Update(&updatedUser)
		if err != nil {
			b.Fatalf("xorm update: %v", err)
		}
		if affected == 0 {
			// 允许影响 0 行的情况，因为可能没有匹配的记录
		}
	}
}

// BenchmarkJormUpdateAll 测试 jorm 更新所有记录
func BenchmarkJormUpdateAll(b *testing.B) {
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
		updatedUser := User{
			Age: 25 + i%10,
		}

		// 使用 WHERE 1=1 来更新所有记录
		result, err := jorm.Model(&User{}, engine).Where("1=1").Update(updatedUser)
		if err != nil {
			b.Fatalf("jorm update: %v", err)
		}
		if result == 0 {
			b.Fatalf("jorm update affected 0 rows")
		}
	}
}

// BenchmarkGormUpdateAll 测试 gorm 更新所有记录
func BenchmarkGormUpdateAll(b *testing.B) {
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
		updatedUser := User{
			Age: 25 + i%10,
		}

		// 使用 WHERE 1=1 来更新所有记录
		result := db.Model(&User{}).Where("1=1").Updates(&updatedUser)
		if result.Error != nil {
			b.Fatalf("gorm update: %v", result.Error)
		}
		if result.RowsAffected == 0 {
			b.Fatalf("gorm update affected 0 rows")
		}
	}
}

// BenchmarkXormUpdateAll 测试 xorm 更新所有记录
func BenchmarkXormUpdateAll(b *testing.B) {
	setupTestData(b, 1000)

	engine, err := NewXormEngine()
	if err != nil {
		b.Fatalf("new xorm engine: %v", err)
	}
	defer engine.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		updatedUser := User{
			Age: 25 + i%10,
		}

		affected, err := engine.Table("users").Update(&updatedUser)
		if err != nil {
			b.Fatalf("xorm update: %v", err)
		}
		if affected == 0 {
			b.Fatalf("xorm update affected 0 rows")
		}
	}
}
