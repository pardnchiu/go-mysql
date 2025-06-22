# MySQL Pool
> 一個支持鏈式呼叫的 Golang MySQL 包裝器，具備讀寫分離設置、查詢建構器等功能，提供完整的連線管理。<br>
> A Golang MySQL wrapper with chainable calls, featuring read-write separation, query builder, and complete connection management.<br>
>
> Node.js version [here](https://github.com/pardnchiu/node-mysql-pool) |  PHP version [here](https://github.com/pardnchiu/php-mysql-pool)

![lang](https://img.shields.io/github/languages/top/pardnchiu/go-mysql)
[![license](https://img.shields.io/github/license/pardnchiu/go-mysql)](LICENSE) 
[![version](https://img.shields.io/github/v/tag/pardnchiu/go-mysql)](https://github.com/pardnchiu/go-mysql/releases) 

## 三大主軸 / Three Core Features

### 讀寫分離配置 / Read-Write Separation
支援讀寫連線池配置，增加資料庫預連接，提高連接效率<br>
Supports read-write connection pool configuration with database pre-connections to improve connection efficiency

### 查詢建構器 / Query Builder
支持鏈式語法的 SQL 查詢建構介面，防止 SQL 注入攻擊<br>
Chainable SQL query interface that prevents SQL injection attacks

### CRUD 操作 / CRUD Operation
完整的新增、查詢、更新、刪除操作支援<br>
Complete support for Create, Read, Update, and Delete operations

## 依賴套件 / Dependencies

- [`github.com/go-sql-driver/mysql`](https://github.com/go-sql-driver/mysql)
- [`github.com/pardnchiu/go-logger`](https://github.com/pardnchiu/go-logger)<br>
  如果你不需要，你可以 fork 然後使用你熟悉的取代。更可以到[這裡](https://forms.gle/EvNLwzpHfxWR2gmP6)進行投票讓我知道。<br>
  If you don't need this, you can fork the project and replace it. You can also vote [here](https://forms.gle/EvNLwzpHfxWR2gmP6) to let me know your thought.

## 使用方法 / How to use

### 安裝 / Installation
```bash
go get github.com/pardnchiu/go-mysql
```

### 初始化 / Initialization
```go
package main

import (
  "fmt"
  "log"
  
  mp "github.com/pardnchiu/go-mysql"
)

func main() {
  config := mp.Config{
    Read: &mp.DBConfig{
      Host:       "localhost",
      Port:       3306,
      User:       "root",
      Password:   "password",
      Charset:    "utf8mb4",
      Connection: 10,
    },
    Write: &mp.DBConfig{
      Host:       "localhost",
      Port:       3306,
      User:       "root",
      Password:   "password",
      Charset:    "utf8mb4",
      Connection: 5,
    },
  }
  
  // Initialize
  pool, err := mp.New(config)
  if err != nil {
    log.Fatal(err)
  }
  defer pool.Close()
  
  // Insert data
  userData := map[string]interface{}{
    "name":  "John Doe",
    "email": "john@example.com",
    "age":   30,
  }
  
  lastID, err := pool.Write.
    DB("myapp").
    Table("users").
    Insert(userData)
  if err != nil {
    log.Fatal(err)
  }
  
  fmt.Printf("插入使用者 ID: %d\n", lastID)
  
  // Query data
  rows, err := pool.Read.
    DB("myapp").
    Table("users").
    Select("id", "name", "email").
    Where("age", ">", 18).
    OrderBy("created_at", "DESC").
    Limit(10).
    Get()
  if err != nil {
    log.Fatal(err)
  }
  defer rows.Close()
  
  for rows.Next() {
    var id int
    var name, email string
    err := rows.Scan(&id, &name, &email)
    if err != nil {
      log.Fatal(err)
    }
    fmt.Printf("User: %s (%s)\n", name, email)
  }
}
```

## 配置介紹 / Configuration

```go
type Config struct {
  Read  *DBConfig
  Write *DBConfig
  Log   *Log
}

type DBConfig struct {
  Host       string // Database host address
  Port       int    // Database port
  User       string // Database username
  Password   string // Database password
  Charset    string // Character set (default: utf8mb4)
  Connection int    // Maximum connections
}

type Log struct {
  Path      string // Log directory path (default: ./logs/goMysql)
  Stdout    bool   // Enable console output (default: false)
  MaxSize   int64  // Maximum size before file rotation (default: 16*1024*1024)
  MaxBackup int    // Number of log files to retain (default: 5)
  Type      string // Output format: "json" for slog standard, "text" for tree format (default: "text")
}
```

## 支持操作 / Operations

### 查詢 / Queries
```go
// Basic query
rows, err := pool.Read.
  DB("database_name").
  Table("users").
  Select("id", "name", "email").
  Where("status", "active").
  Get()

// Complex conditional query
rows, err := pool.Read.
  DB("database_name").
  Table("users").
  Select("*").
  Where("age", ">", 18).
  Where("status", "active").
  Where("name", "LIKE", "John").
  OrderBy("created_at", "DESC").
  Limit(10).
  Offset(20).
  Get()

// JOIN query
rows, err := pool.Read.
  DB("database_name").
  Table("users").
  Select("users.name", "profiles.bio").
  LeftJoin("profiles", "users.id", "profiles.user_id").
  Where("users.status", "active").
  Get()

// Count total
rows, err := pool.Read.
  DB("database_name").
  Table("users").
  Select("id", "name").
  Where("status", "active").
  Total().
  Limit(10).
  Get()
```

### CRUD
```go
// Insert data
data := map[string]interface{}{
  "name":  "Jane Doe",
  "email": "jane@example.com",
  "age":   25,
}

lastID, err := pool.Write.
  DB("database_name").
  Table("users").
  Insert(data)

// Update data
updateData := map[string]interface{}{
  "age":    26,
  "status": "updated",
}

result, err := pool.Write.
  DB("database_name").
  Table("users").
  Where("id", 1).
  Update(updateData)

// Upsert operation
data := map[string]interface{}{
  "email": "unique@example.com",
  "name":  "New User",
}

updateData := map[string]interface{}{
  "name": "Updated User",
  "last_login": "NOW()",
}

lastID, err := pool.Write.
  DB("database_name").
  Table("users").
  Upsert(data, updateData)

// Increment values
result, err := pool.Write.
  DB("database_name").
  Table("users").
  Where("id", 1).
  Increase("view_count", 1).
  Update()
```

### SQL
```go
// Direct query
rows, err := pool.Read.Query("SELECT * FROM users WHERE age > ?", 18)

// Direct execution
result, err := pool.Write.Exec("UPDATE users SET last_login = NOW() WHERE id = ?", userID)
```

## 可用函式 / Functions

### 連線池管理 / Pool Management
- **New** - 建立新的連線池 / Create new connection
  ```go
  pool, err := mp.New(config)
  ```
  - 初始化讀寫分離的連線池<br>
    Initialize read-write separated connection pool
  - 驗證資料庫連線可用性<br>
    Validate database connection availability

- **Close** - 關閉連線池 / Close connection
  ```go
  err := pool.Close()
  ```
  - 關閉所有連線<br>
    Close all connections
  - 等待進行中的查詢完成<br>
    Wait for ongoing queries to complete
  - 釋放系統資源<br>
    Release system resources

### 查詢建構 / Query Building
- **DB** - 指定資料庫 / Select database
  ```go
  builder := pool.Read.DB("database_name")
  ```

- **Table** - 指定資料表 / Select table
  ```go
  builder := builder.Table("table_name")
  ```

- **Select** - 選擇欄位 / Select columns
  ```go
  builder := builder.Select("col1", "col2", "col3")
  ```

- **Where** - 篩選條件 / Filter data
  ```go
  builder := builder.Where("column", "value")
  builder := builder.Where("column", ">", "value")
  builder := builder.Where("column", "LIKE", "pattern")
  ```

- **Join** - 資料表聯結 / Table joins
  ```go
  builder := builder.LeftJoin("table2", "table1.id", "table2.foreign_id")
  builder := builder.RightJoin("table2", "table1.id", "table2.foreign_id")
  builder := builder.InnerJoin("table2", "table1.id", "table2.foreign_id")
  ```

### 資料操作 / Data
- **Insert** - 插入資料 / Insert data
  ```go
  lastID, err := builder.Insert(data)
  ```

- **Update** - 更新資料 / Update data
  ```go
  result, err := builder.Update(data)
  ```

- **Upsert** - 插入或更新 / Insert or update
  ```go
  lastID, err := builder.Upsert(insertData, updateData)
  ```

## 授權條款 / License

此原始碼專案採用 [MIT](LICENSE) 授權條款。<br>
This source code project is licensed under the [MIT](LICENSE) license.

## 作者 / Author

<img src="https://avatars.githubusercontent.com/u/25631760" align="left" width="96" height="96" style="margin-right: 0.5rem;">

<h4 style="padding-top: 0">邱敬幃 Pardn Chiu</h4>

<a href="mailto:dev@pardn.io" target="_blank">
  <img src="https://pardn.io/image/email.svg" width="48" height="48">
</a> <a href="https://linkedin.com/in/pardnchiu" target="_blank">
  <img src="https://pardn.io/image/linkedin.svg" width="48" height="48">
</a>

***

©️ 2025 [邱敬幃 Pardn Chiu](https://pardn.io)
