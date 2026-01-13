# JORM - Golang ORM Framework

> **注意**: 当前项目基于AI开发，功能暂未完善，单管仅供学习使用，请无用于生产环境。

JORM 计划是开发一个轻量级、高性能的 Golang ORM 框架，基于 Go 1.18+ 泛型设计。它旨在提供简单、直观且类型安全的数据库操作体验，融合了 ActiveRecord 的优雅风格。

## 性能对比

JORM 在常见数据库操作上展现出优异的性能表现，以下是与其他主流 ORM 框架的性能对比测试结果：

| 操作类型 | ORM 框架 | QPS (每秒查询次数) | 执行时间 (ns/op) | 内存分配 (B/op) | 分配次数 (allocs/op) |
|---------|---------|------------------|----------------|---------------|-------------------|
| **插入** | JORM | 2,128 | 469,885 | 1,884 | 36 |
| | GORM | 1,940 | 515,512 | 5,831 | 87 |
| | XORM | 2,103 | 475,414 | 2,341 | 46 |
| **查询** | JORM | 36,566 | 27,346 | 2,203 | 53 |
| | GORM | 31,625 | 31,626 | 3,762 | 63 |
| | XORM | 29,352 | 34,082 | 4,295 | 121 |
| **更新** | JORM | 2,243 | 445,765 | 2,162 | 39 |
| | GORM | 1,969 | 507,858 | 6,479 | 80 |
| | XORM | 2,189 | 456,972 | 3,582 | 91 |

### 性能分析

1. **查询性能优势显著**：JORM 在查询操作上性能领先，QPS 比 GORM 高约 15.6%，比 XORM 高约 24.6%。

2. **内存使用效率高**：在所有操作类型中，JORM 的内存分配量都是最低的，特别是在插入操作中，内存分配量仅为 GORM 的约 32%。

3. **GC 压力小**：JORM 的内存分配次数最少，有助于减少垃圾回收压力，提高应用整体性能。

4. **写入性能稳定**：在插入和更新操作中，JORM 性能与 XORM 相当，但明显优于 GORM。

## 测试环境

- **操作系统**: macOS (Darwin 24.6.0)
- **CPU**: Intel(R) Core(TM) i7-8700 CPU @ 3.20GHz
- **架构**: amd64
- **数据库**: SQLite3 (本地文件数据库)
- **测试数据量**: 
  - 插入/更新测试: 每次操作清空表，逐条插入/更新
  - 查询测试: 预先插入 1000 条测试数据
- **Go版本**: Go 1.18+
- **测试方法**: 使用 Go 标准库 testing 包的基准测试功能
- **测试时间**: 2026年1月13日

> 注：所有测试均使用相同的硬件环境和测试数据，确保结果的可比性。SQLite3 作为文件数据库，测试结果可能与其他数据库(如MySQL、PostgreSQL)有所差异。


## 目录

- [为什么选择 JORM (Why JORM)](#为什么选择-jorm-why-jorm)
- [快速开始 (Quick Start)](#快速开始-quick-start)
- [支持的数据库 (Supported Databases)](#支持的数据库-supported-databases)
- [核心功能 (Core Features)](#核心功能-core-features)
  - [连接数据库 (Connection)](#连接数据库-connection)
  - [定义模型 (Model Definition)](#定义模型-model-definition)
  - [自动迁移 (Auto Migration)](#自动迁移-auto-migration)
  - [CRUD 操作](#crud-操作)
  - [链式查询 (Chain Query)](#链式查询-chain-query)
- [进阶使用 (Advanced Usage)](#进阶使用-advanced-usage)
  - [多数据库支持 (Multi-Database)](#多数据库支持-multi-database)
  - [事务 (Transaction)](#事务-transaction)
  - [作用域 (Scopes)](#作用域-scopes)
  - [原生 SQL (Raw SQL)](#原生-sql-raw-sql)

## 安装 (Installation)

```bash
go get github.com/shrek82/jorm
```

## 快速开始 (Quick Start)

### 1. 定义模型

```go
package main

import "time"

type User struct {
    ID        int64     `jorm:"primaryKey;autoIncrement"`
    Name      string    `jorm:"column:username;size:100;not null"`
    Age       int       `jorm:"default:18"`
    CreatedAt time.Time
}
```

### 2. 初始化与使用

以 SQLite 为例：

```go
package main

import (
    "fmt"
    "github.com/shrek82/jorm"
    "github.com/shrek82/jorm/driver/sqlite"
)

func main() {
    // 1. 连接数据库
    db, err := jorm.Open(sqlite.Open("test.db"))
    if err != nil {
        panic("failed to connect database")
    }

    // 2. 设置为全局默认数据库 (推荐)
    jorm.SetDefault(db)

    // 3. 自动迁移表结构
    // 自动创建 users 表
    err = db.AutoMigrate(&User{})
    if err != nil {
        panic(err)
    }

    // 4. 插入数据
    user := &User{Name: "Tom", Age: 20}
    err = jorm.Model(&User{}).Create(user)
    if err != nil {
        panic(err)
    }
    fmt.Printf("User ID: %d\n", user.ID)

    // 5. 查询数据
    var result User
    // SELECT * FROM users WHERE username = 'Tom' LIMIT 1
    err = jorm.Model(&User{}).Where("username = ?", "Tom").Find(&result)
    if err != nil {
        fmt.Println("User not found")
    } else {
        fmt.Printf("Found: %+v\n", result)
    }
}
```

## 支持的数据库 (Supported Databases)

JORM 目前支持以下数据库驱动：

- **SQLite**: `github.com/shrek82/jorm/driver/sqlite`
- **MySQL**: `github.com/shrek82/jorm/driver/mysql`
- **PostgreSQL**: `github.com/shrek82/jorm/driver/postgres`
- **SQL Server**: `github.com/shrek82/jorm/driver/sqlserver`
- **Oracle**: `github.com/shrek82/jorm/driver/oracle`

## 核心功能 (Core Features)

### 连接数据库 (Connection)

使用 `jorm.Open` 配合相应的驱动连接数据库。

#### MySQL
```go
import "github.com/shrek82/jorm/driver/mysql"

dsn := "user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
db, err := jorm.Open(mysql.Open(dsn))
```

#### PostgreSQL
```go
import "github.com/shrek82/jorm/driver/postgres"

dsn := "postgres://user:password@localhost:5432/dbname?sslmode=disable"
db, err := jorm.Open(postgres.Open(dsn))
```

#### SQL Server
```go
import "github.com/shrek82/jorm/driver/sqlserver"

dsn := "sqlserver://user:password@localhost:1433?database=dbname"
db, err := jorm.Open(sqlserver.Open(dsn))
```

#### Oracle
```go
import "github.com/shrek82/jorm/driver/oracle"

dsn := "oracle://user:password@localhost:1521/service_name"
db, err := jorm.Open(oracle.Open(dsn))
```

### 定义模型 (Model Definition)

使用 Struct Tag 定义表结构映射。支持的 Tag 如下：

| Tag | 描述 | 示例 |
| --- | --- | --- |
| `primaryKey` | 标记为主键 | `jorm:"primaryKey"` |
| `autoIncrement` | 标记为自增列 | `jorm:"autoIncrement"` |
| `column` | 自定义列名 | `jorm:"column:user_name"` |
| `type` | 自定义数据库类型 | `jorm:"type:varchar(100)"` |
| `size` | 字段长度 | `jorm:"size:255"` |
| `default` | 默认值 | `jorm:"default:'active'"` |
| `not null` | 非空约束 | `jorm:"not null"` |
| `unique` | 唯一索引 | `jorm:"unique"` |
| `-` | 忽略该字段 | `jorm:"-"` |
| `->` | 只读字段 | `jorm:"->"` |
| `<-` | 只写字段 | `jorm:"<-"` |

**自定义表名**

可以通过实现 `TableName()` 方法来自定义表名：

```go
type User struct {
    ID   int64  `jorm:"primaryKey;autoIncrement"`
    Name string `jorm:"column:username"`
    Age  int
}

// TableName 自定义表名
func (u *User) TableName() string {
    return "my_users"
}
```

### 自动迁移 (Auto Migration)

`AutoMigrate` 会自动创建表、添加缺失的列。它**不会**删除列或更改现有列的类型以保护数据。

```go
db.AutoMigrate(&User{}, &Product{}, &Order{})
```

### CRUD 操作

JORM 使用泛型 `Session[T]` 进行类型安全的操作。

#### 创建 (Create)

**插入单条记录**
```go
user := &User{Name: "Jerry", Age: 18}
jorm.Model(&User{}).Create(user)
// 此时 user.ID 会被自动回填
fmt.Printf("Inserted user ID: %d\n", user.ID)
```

**批量插入 (结构体切片)**
```go
users := []User{
    {Name: "Tom", Age: 20},
    {Name: "Jerry", Age: 22},
    {Name: "Mike", Age: 25},
}
jorm.Model(&User{}).Create(users)
```

**批量插入 (指针切片)**
```go
u1 := &User{Name: "Alice", Age: 18}
u2 := &User{Name: "Bob", Age: 20}
users := []*User{u1, u2}
jorm.Model(&User{}).Create(users)
```

**使用 Map 插入**
```go
userMap := map[string]interface{}{
    "Name": "Charlie",
    "Age":  30,
}
jorm.Model(&User{}).Create(userMap)
```

#### 查询 (Read)

**查询单条 (Find)**
```go
var user User
// 获取第一条记录，按主键排序
jorm.Model(&User{}).Find(&user)

// 带字符串条件
jorm.Model(&User{}).Where("name = ?", "Jerry").Find(&user)

// 使用主键值作为条件
jorm.Model(&User{}).Where(123).Find(&user)

// 使用结构体作为条件 (自动映射非零值字段)
jorm.Model(&User{}).Where(User{Name: "Jerry"}).Find(&user)

// 使用 Map 作为条件
jorm.Model(&User{}).Where(map[string]interface{}{"age": 20}).Find(&user)
```

**查询列表 (FindAll)**
```go
var users []User
// 获取所有记录
jorm.Model(&User{}).FindAll(&users)

// 带条件和排序
jorm.Model(&User{}).Where("age > ?", 18).Order("age DESC").FindAll(&users)

// 使用 Limit 和 Offset 分页
jorm.Model(&User{}).Order("id DESC").Limit(10).Offset(20).FindAll(&users)

// 指定查询字段
jorm.Model(&User{}).Select("name", "age").FindAll(&users)

// 多次 Select 会追加字段
jorm.Model(&User{}).Select("name").Select("age").FindAll(&users)

// 使用 Joins 连表查询
jorm.Model(&User{}).
    Joins("LEFT JOIN orders ON orders.user_id = users.id").
    Where("orders.status = ?", "completed").
    FindAll(&users)
```

**Count 计数**
```go
// 统计所有记录数
count, err := jorm.Model(&User{}).Count()
fmt.Printf("Total users: %d\n", count)

// 带条件统计
count, err := jorm.Model(&User{}).Where("age > ?", 18).Count()

// 统计特定字段
count, err := jorm.Model(&User{}).Count("id")
```

**Sum 求和**
```go
// 计算年龄总和
totalAge, err := jorm.Model(&User{}).Sum("Age")
fmt.Printf("Total age: %.2f\n", totalAge)

// 带条件求和
totalAge, err := jorm.Model(&User{}).Where("age > ?", 18).Sum("Age")
```

**FindOrCreate 查找或创建**
```go
var user User
// 如果找到返回现有记录，否则创建新记录
jorm.Model(&User{}).
    Where("name = ?", "Tom").
    FindOrCreate(&user, &User{Name: "Tom", Age: 20})

// 使用 Map 创建
jorm.Model(&User{}).
    Where("name = ?", "Alice").
    FindOrCreate(&user, map[string]interface{}{
        "Name": "Alice",
        "Age":  25,
    })
```

#### 更新 (Update)

**Update (更新单条)**

强制要求带 WHERE 条件，返回影响的行数。

```go
// 1. 使用 Map 更新指定字段
// UPDATE users SET age = 20 WHERE id = 1
affected, err := jorm.Model(&User{}).Where("id = ?", 1).Update(map[string]interface{}{"age": 20})
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Updated %d rows\n", affected)

// 2. 使用 Struct 更新 (仅更新非零值字段)
// UPDATE users SET age = 100, name = 'Jerry' WHERE id = 1
updateUser := &User{Age: 100, Name: "Jerry"}
affected, err = jorm.Model(&User{}).Where("id = ?", 1).Update(updateUser)

// 3. 零值字段不会被更新 (如果需要更新零值，请使用 Save 或 Map)
updateUser := &User{Age: 0, Name: "NewName"}
affected, err = jorm.Model(&User{}).Where("id = ?", 1).Update(updateUser)
// 只有 name 会被更新，age 保持原值
```

**UpdateAll (批量更新)**

用于更新多条记录，必须带 WHERE 条件以防止全局更新。

```go
// UPDATE users SET age = 25 WHERE age > 18
affected, err := jorm.Model(&User{}).Where("age > ?", 18).UpdateAll(map[string]interface{}{"age": 25})
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Updated %d rows\n", affected)

// 使用 Struct 批量更新
affected, err = jorm.Model(&User{}).Where("status = ?", "inactive").UpdateAll(&User{Status: "active"})
```

**Save (全量更新)**

`Save` 用于保存模型对象的所有字段（包括零值）。该方法要求 struct 必须包含有效的主键值。

```go
user := &User{ID: 111, Name: "UpdatedName", Age: 0}
// UPDATE users SET name = 'UpdatedName', age = 0 WHERE id = 111
affected, err := jorm.Model(&User{}).Save(user)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Saved %d rows\n", affected)

// Save 会自动使用主键作为 WHERE 条件，无需手动指定
```

#### 删除 (Delete)

强制要求带 WHERE 条件，以防止误删全表数据。

```go
// 根据主键删除
var user User
jorm.Model(&User{}).Where(123).Delete(nil)

// 根据条件删除
jorm.Model(&User{}).Where("name = ?", "Jerry").Delete(nil)

// 批量删除
jorm.Model(&User{}).Where("age < ?", 18).Delete(nil)

// 使用结构体条件删除
jorm.Model(&User{}).Where(User{Name: "OldUser"}).Delete(nil)

// 使用 Map 条件删除
jorm.Model(&User{}).Where(map[string]interface{}{"status": "deleted"}).Delete(nil)
```

### 链式查询 (Chain Query)

支持灵活的 SQL 构造，所有条件通过链式调用叠加。

```go
var users []User
jorm.Model(&User{}).
    Where("age > ?", 18).
    Where("name LIKE ?", "T%").
    Order("age DESC").
    Limit(10).
    Offset(20).
    FindAll(&users)
```

**Where 条件**
```go
// 字符串条件
jorm.Model(&User{}).Where("age > ?", 18).Find(&user)

// 多个 AND 条件
jorm.Model(&User{}).
    Where("age > ?", 18).
    Where("name = ?", "Tom").
    Find(&user)

// 结构体条件 (自动转换为 AND)
jorm.Model(&User{}).Where(User{Age: 20, Name: "Tom"}).Find(&user)

// Map 条件
jorm.Model(&User{}).Where(map[string]interface{}{
    "age":  20,
    "name": "Tom",
}).Find(&user)
```

**Order 排序**
```go
// 单字段排序
jorm.Model(&User{}).Order("age DESC").FindAll(&users)

// 多字段排序
jorm.Model(&User{}).Order("age DESC, name ASC").FindAll(&users)

// 多次 Order 调用
jorm.Model(&User{}).Order("age DESC").Order("name ASC").FindAll(&users)
```

**Select 查询字段**
```go
// 查询特定字段
jorm.Model(&User{}).Select("name", "age").Find(&user)

// 使用 SQL 表达式
jorm.Model(&User{}).Select("COUNT(*) as count").Find(&user)

// 追加字段
jorm.Model(&User{}).Select("name").Select("age").Find(&user)
```

**Omit 忽略字段**
```go
// 查询时忽略特定字段
jorm.Model(&User{}).Omit("password", "secret").Find(&user)
```

**Limit 和 Offset**
```go
// 限制返回数量
jorm.Model(&User{}).Limit(10).FindAll(&users)

// 分页查询
jorm.Model(&User{}).Limit(10).Offset(20).FindAll(&users)
```

**Joins 连表查询**
```go
// LEFT JOIN
jorm.Model(&User{}).
    Joins("LEFT JOIN orders ON orders.user_id = users.id").
    Where("orders.status = ?", "completed").
    FindAll(&users)

// 带参数的 JOIN
jorm.Model(&User{}).
    Joins("LEFT JOIN orders ON orders.user_id = users.id AND orders.status = ?", "completed").
    FindAll(&users)
```

**Table 指定表名**
```go
// 覆盖默认表名
jorm.Model(&User{}).Table("custom_users").Find(&user)
```

**With 上下文**
```go
ctx := context.Background()
ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
defer cancel()

jorm.Model(&User{}).WithContext(ctx).Find(&user)
```

**Clone 克隆会话**
```go
// 克隆会话用于不同的查询
baseQuery := jorm.Model(&User{}).Where("age > ?", 18)

// 使用克隆执行不同查询
var users1 []User
baseQuery.Clone().Order("age DESC").FindAll(&users1)

var users2 []User
baseQuery.Clone().Order("name ASC").FindAll(&users2)
```

**Rows 流式处理**
```go
// 使用 Rows 迭代处理大结果集，避免内存溢出
rows, err := jorm.Model(&User{}).Where("age > ?", 18).Rows()
if err != nil {
    log.Fatal(err)
}
defer rows.Close()

for rows.Next() {
    var user User
    if err := rows.Scan(&user); err != nil {
        log.Fatal(err)
    }
    fmt.Printf("User: %+v\n", user)
}
```

**错误处理**

JORM 的 Session 对象包含 `Error` 字段，可以检查操作是否出错：

```go
// 检查 Error 字段
session := jorm.Model(&User{}).Where("invalid_column = ?", "value")
if session.Error != nil {
    log.Fatal(session.Error)
}

// 在链式调用中，Error 会累积
session = jorm.Model(&User{}).Where("age > ?", 18)
if session.Error != nil {
    log.Fatal(session.Error)
}

// 执行查询时返回错误
err := jorm.Model(&User{}).Where("name = ?", "Tom").Find(&user)
if err != nil {
    log.Fatal(err)
}
```

## 进阶使用 (Advanced Usage)

### 多数据库支持 (Multi-Database)

JORM 设计灵活，支持同时操作多个数据库。

**设置全局默认数据库**：
```go
jorm.SetDefault(db)
// 后续操作默认使用 db
jorm.Model(&User{}).Find(&user)
```

**显示 SQL 日志**：
```go
// 开启 SQL 日志输出
jorm.ShowSql(true)
```

**指定数据库操作**：
```go
// 方式一：在 Model 初始化时指定
jorm.Model(&User{}, db2).Find(&user)

// 方式二：中途切换数据库 (读写分离)
jorm.Model(&User{}).SetDB(dbRead).Find(&user)
```

### 事务 (Transaction)

使用 `db.Transaction` 闭包自动管理事务提交与回滚。如果在闭包内返回 error，事务会自动回滚；返回 nil 则自动提交。

```go
err := db.Transaction(func(tx *jorm.DB) error {
    // 注意：在事务中必须将 tx 传递给 Model
    
    // 插入数据
    user := &User{Name: "TxUser", Age: 20}
    if err := jorm.Model(user, tx).Create(user); err != nil {
        return err
    }
    
    // 更新数据
    if _, err := jorm.Model(&User{}, tx).Where("id = ?", user.ID).Update(map[string]interface{}{"age": 25}); err != nil {
        return err
    }
    
    // 删除数据
    if err := jorm.Model(&User{}, tx).Where("name = ?", "OldName").Delete(nil); err != nil {
        return err
    }
    
    return nil
})
```

**手动事务管理**
```go
// 开启事务
tx := db.Begin()
if tx.Error != nil {
    log.Fatal(tx.Error)
}

// 执行操作
user := &User{Name: "ManualTx", Age: 30}
if err := jorm.Model(user, tx).Create(user); err != nil {
    tx.Rollback()
    log.Fatal(err)
}

// 提交事务
tx.Commit()
if tx.Error != nil {
    log.Fatal(tx.Error)
}
```

### 作用域 (Scopes)

定义可复用的查询逻辑，保持代码整洁。

```go
// 定义 Scope
func Active(s *jorm.Session[User]) *jorm.Session[User] {
    return s.Where("status = ?", "active")
}

func Paginate(page, pageSize int) func(*jorm.Session[User]) *jorm.Session[User] {
    return func(s *jorm.Session[User]) *jorm.Session[User] {
        return s.Limit(pageSize).Offset((page - 1) * pageSize)
    }
}

// 使用 Scope
var users []User
jorm.Model(&User{}).
    Scope(Active, Paginate(1, 10)).
    FindAll(&users)
```

### 原生 SQL (Raw SQL)

如果 ORM 无法满足需求，可以使用 `db.Exec` 和 `db.Query` 执行原生 SQL。

```go
// 执行非查询 SQL (INSERT, UPDATE, DELETE)
result, err := db.Exec(context.Background(), "UPDATE users SET age = ? WHERE name = ?", 20, "Tom")
if err != nil {
    log.Fatal(err)
}
affected, _ := result.RowsAffected()
fmt.Printf("Affected rows: %d\n", affected)

// 执行查询 SQL
rows, err := db.Query(context.Background(), "SELECT id, name, age FROM users WHERE age > ?", 18)
if err != nil {
    log.Fatal(err)
}
defer rows.Close()

for rows.Next() {
    var id int64
    var name string
    var age int
    if err := rows.Scan(&id, &name, &age); err != nil {
        log.Fatal(err)
    }
    fmt.Printf("User: ID=%d, Name=%s, Age=%d\n", id, name, age)
}
```

---

**License**

MIT
