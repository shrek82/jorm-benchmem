database:
  driver: "mysql"           # 数据库驱动
  host: "127.0.0.1"       # MySQL服务器地址
  port: "3306"           # MySQL端口，默认3306
  username: "root"        # 数据库用户名
  password: ""             # 数据库密码，留空表示无密码
  database: "test"  # 数据库名称
  charset: "utf8mb4"      # 字符集
  parseTime: true          # 解析时间类型
  loc: "Local"            # 时区
  debug_sql: true         # 是否打印执行的SQL语句，默认false，设置为true开启调试

user表结构
  type User struct {
	ID        int    `db:"id"`
	Name      string `db:"name"`
	Email     string `db:"email"`
	Age       int    `db:"age"`
	Active    bool   `db:"active"`
	CreatedAt string `db:"created_at"`
}