> [!Note]
> This content is translated by LLM. Original text can be found [here](README.zh.md)

# MySQL Pool
> A Golang MySQL wrapper supporting chain calls, read-write separation, query builder, and complete connection management.<br>
>
> Node.js version [here](https://github.com/pardnchiu/node-mysql-pool) | PHP version [here](https://github.com/pardnchiu/php-mysql-pool)

[![lang](https://img.shields.io/badge/lang-Go-blue)](README.zh.md) 
[![license](https://img.shields.io/github/license/pardnchiu/go-mysql)](LICENSE)
[![version](https://img.shields.io/github/v/tag/pardnchiu/go-mysql)](https://github.com/pardnchiu/go-mysql/releases)
![card](https://goreportcard.com/badge/github.com/pardnchiu/go-mysql)<br>
[![readme](https://img.shields.io/badge/readme-EN-white)](README.md)
[![readme](https://img.shields.io/badge/readme-ZH-white)](README.zh.md) 

## Key Features

### Read-Write Separation
Supports read-write connection pool configuration, enabling pre-connections to improve efficiency.

### Query Builder
Provides a chainable SQL query builder interface to prevent SQL injection attacks.

### CRUD Operations
Complete support for create, read, update, and delete operations.

## Dependencies

- [`github.com/go-sql-driver/mysql`](https://github.com/go-sql-driver/mysql)
- [`github.com/pardnchiu/go-logger`](https://github.com/pardnchiu/go-logger)<br>
  If not needed, you can fork and replace it with your preferred logger. You can also vote [here](https://forms.gle/EvNLwzpHfxWR2gmP6) to share your feedback.

## Usage

### Installation
```bash
go get github.com/pardnchiu/go-mysql
```

### Initialization
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
  
  fmt.Printf("Inserted user ID: %d\n", lastID)
  
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

## Configuration Overview

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

## Supported Operations

### Query
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

## Available Functions

### Connection Pool Management
- **New** - Create a new connection pool
  ```go
  pool, err := mp.New(config)
  ```
  - Initializes read-write separation connection pool
  - Validates database connection availability

- **Close** - Close the connection pool
  ```go
  err := pool.Close()
  ```
  - Closes all connections
  - Waits for ongoing queries to complete
  - Releases system resources

### Query Builder
- **DB** - Specify database
  ```go
  builder := pool.Read.DB("database_name")
  ```

- **Table** - Specify table
  ```go
  builder := builder.Table("table_name")
  ```

- **Select** - Select columns
  ```go
  builder := builder.Select("col1", "col2", "col3")
  ```

- **Where** - Filter conditions
  ```go
  builder := builder.Where("column", "value")
  builder := builder.Where("column", ">", "value")
  builder := builder.Where("column", "LIKE", "pattern")
  ```

- **Join** - Table joins
  ```go
  builder := builder.LeftJoin("table2", "table1.id", "table2.foreign_id")
  builder := builder.RightJoin("table2", "table1.id", "table2.foreign_id")
  builder := builder.InnerJoin("table2", "table1.id", "table2.foreign_id")
  ```

### Data Operations
- **Insert** - Insert data
  ```go
  lastID, err := builder.Insert(data)
  ```

- **Update** - Update data
  ```go
  result, err := builder.Update(data)
  ```

- **Upsert** - Insert or update
  ```go
  lastID, err := builder.Upsert(insertData, updateData)
  ```

## License

This project is licensed under the [MIT](LICENSE) license.

## Author

<img src="https://avatars.githubusercontent.com/u/25631760" align="left" width="96" height="96" style="margin-right: 0.5rem;">

<h4 style="padding-top: 0">邱敬幃 Pardn Chiu</h4>

<a href="mailto:dev@pardn.io" target="_blank">
  <img src="https://pardn.io/image/email.svg" width="48" height="48">
</a> <a href="https://linkedin.com/in/pardnchiu" target="_blank">
  <img src="https://pardn.io/image/linkedin.svg" width="48" height="48">
</a>

***

©️ 2025 [邱敬幃 Pardn Chiu](https://pardn.io)
