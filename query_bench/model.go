package query_bench

// User 用于三种 ORM 统一对比的模型
// 注意：struct tag 同时包含 jorm / gorm / xorm 的配置

type User struct {
	ID   int64  `jorm:"primaryKey;autoIncrement" gorm:"primaryKey;autoIncrement" xorm:"'id' pk autoincr"`
	Name string `jorm:"column:username" gorm:"column:username" xorm:"'username'"`
	Age  int    `jorm:"column:age" gorm:"column:age" xorm:"'age'"`
}

// TableName 使三种 ORM 使用同一张表
func (User) TableName() string {
	return "users"
}
