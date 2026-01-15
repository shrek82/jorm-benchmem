package update_bench

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/shrek82/jorm"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"xorm.io/xorm"

	_ "github.com/mattn/go-sqlite3"
)

// DefaultDSN 从环境变量读取 DSN，找不到则使用本地 SQLite 数据库
func DefaultDSN() string {
	if dsn := os.Getenv("BENCH_DSN"); dsn != "" {
		return dsn
	}
	// 使用本地 SQLite 文件数据库
	return "test.db"
}

// NewSQLDB 返回底层 *sql.DB，主要用于在 benchmark 中做 DELETE 等操作
func NewSQLDB() (*sql.DB, error) {
	return sql.Open("sqlite3", DefaultDSN())
}

// NewJormEngine 初始化 jorm 引擎（返回 *jorm.DB）
func NewJormEngine() (*jorm.DB, error) {
	db, err := jorm.Open("sqlite3", DefaultDSN(), nil)
	if err != nil {
		return nil, fmt.Errorf("open jorm: %w", err)
	}
	return db, nil
}

// NewGormDB 初始化 gorm DB
func NewGormDB() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(DefaultDSN()), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("open gorm: %w", err)
	}
	return db, nil
}

// NewXormEngine 初始化 xorm Engine
func NewXormEngine() (*xorm.Engine, error) {
	engine, err := xorm.NewEngine("sqlite3", DefaultDSN())
	if err != nil {
		return nil, fmt.Errorf("open xorm: %w", err)
	}
	return engine, nil
}
