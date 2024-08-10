GORM（Go Object Relational Mapping）是一个流行的Go语言 ORM（对象关系映射）库，用于简化与关系型数据库的交互。

### 1. 连接数据库

首先，要使用GORM，需要配置数据库连接。GORM支持多种数据库，例如MySQL、PostgreSQL、SQLite等。以下是连接MySQL数据库的示例：

```
import (
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
)

dsn := "user:password@tcp(localhost:3306)/database_name?charset=utf8mb4&parseTime=True&loc=Local"
db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
if err != nil {
    panic("failed to connect database")
}
```

### 2. 定义模型

在GORM中，模型是通过Go语言的结构体来定义的。每个模型对应数据库中的一张表，结构体的字段对应表中的列。可以使用标签（tags）来定义字段名、主键、外键等属性。

```
type User struct {
    ID        uint   `gorm:"primaryKey"`
    Name      string
    Email     string
    CreatedAt time.Time
    UpdatedAt time.Time
}
```

### 3. 自动迁移

GORM提供了自动迁移（AutoMigrate）功能，可以根据模型自动创建或更新数据库表结构。

```
err := db.AutoMigrate(&User{})
if err != nil {
    panic("failed to migrate database")
}
```

### 4. 创建记录

使用 `Create` 方法创建新的记录：

```
user := User{Name: "Alice", Email: "alice@example.com"}
result := db.Create(&user)
if result.Error != nil {
    panic("failed to create user")
}
```

### 5. 查询记录

GORM提供了丰富的查询方法，支持链式调用和条件构造。以下是一些常见的查询示例：

- 根据主键查询：

```
var user User
db.First(&user, 1) // 查询ID为1的用户
```

- 条件查询：

```
var users []User
db.Where("name = ?", "Alice").Find(&users) // 查询名为Alice的用户
```

- 链式查询：

```
db.Where("name LIKE ?", "%A%").Order("created_at desc").Limit(10).Find(&users)
```

### 6. 更新记录

使用 `Update` 方法更新记录：

```
db.Model(&user).Update("Name", "Alice Smith")
```

### 7. 删除记录

使用 `Delete` 方法删除记录：

```
db.Delete(&user)
```

### 8. 关联

GORM支持定义和处理模型之间的关联关系，如一对一、一对多、多对多等。通过预加载和懒加载可以有效管理和查询关联数据。

### 9. 事务管理

GORM允许在事务中执行多个操作，以确保操作的原子性和数据的一致性。

```
tx := db.Begin()
defer func() {
    if r := recover(); r != nil {
        tx.Rollback()
    }
}()

// 执行事务操作
tx.Create(&User{Name: "Bob"})
tx.Delete(&User{}, "age < ?", 18)

tx.Commit()
```

### 10. 回调和钩子

通过注册回调函数，可以在模型的生命周期中的不同阶段执行特定的逻辑，如创建前、更新后等。

```
func (user *User) BeforeCreate(tx *gorm.DB) (err error) {
    // 在创建前执行逻辑
    return
}

func (user *User) AfterUpdate(tx *gorm.DB) (err error) {
    // 在更新后执行逻辑
    return
}
```

### 11. 插件和扩展

GORM支持使用插件来扩展其功能，例如缓存支持、日志记录等。

### 12. 性能优化

在使用GORM时，可以注意一些性能优化技巧，如批量操作、预加载数据、使用原生SQL等，以提升数据库操作的效率和性能。