# mysql

```
package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// 连接到 MySQL 数据库
	dsn := "user:password@tcp(127.0.0.1:3306)/mytest"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("连接数据库失败: %v", err)
	}
	defer db.Close()

	// 创建数据库表
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS users (
		id INT AUTO_INCREMENT PRIMARY KEY,
		name VARCHAR(100) NOT NULL,
		email VARCHAR(100) NOT NULL UNIQUE,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`)
	if err != nil {
		log.Fatalf("创建表失败: %v", err)
	}

	// 插入数据
	_, err = db.Exec(`INSERT INTO users (name, email) VALUES (?, ?)`, "Alice", "alice@example.com")
	if err != nil {
		log.Fatalf("插入数据失败: %v", err)
	}

	// 预处理语句插入数据
	stmt, err := db.Prepare(`INSERT INTO users (name, email) VALUES (?, ?)`)
	if err != nil {
		log.Fatalf("准备语句失败: %v", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec("Bob", "bob@example.com")
	if err != nil {
		log.Fatalf("执行预处理语句失败: %v", err)
	}

	// 查询数据
	rows, err := db.Query(`SELECT id, name, email, created_at FROM users`)
	if err != nil {
		log.Fatalf("查询数据失败: %v", err)
	}
	defer rows.Close()

	fmt.Println("查询到的用户:")
	for rows.Next() {
		var id int
		var name, email string
		var createdAt string
		err = rows.Scan(&id, &name, &email, &createdAt)
		if err != nil {
			log.Fatalf("扫描行数据失败: %v", err)
		}
		fmt.Printf("ID: %d, Name: %s, Email: %s, CreatedAt: %s\n", id, name, email, createdAt)
	}

	// 更新数据
	_, err = db.Exec(`UPDATE users SET email = ? WHERE name = ?`, "newalice@example.com", "Alice")
	if err != nil {
		log.Fatalf("更新数据失败: %v", err)
	}

	// 删除数据
	_, err = db.Exec(`DELETE FROM users WHERE name = ?`, "Bob")
	if err != nil {
		log.Fatalf("删除数据失败: %v", err)
	}

	// 事务处理
	tx, err := db.Begin()
	if err != nil {
		log.Fatalf("开始事务失败: %v", err)
	}

	_, err = tx.Exec(`INSERT INTO users (name, email) VALUES (?, ?)`, "Charlie", "charlie@example.com")
	if err != nil {
		tx.Rollback()
		log.Fatalf("事务插入失败: %v", err)
	}

	_, err = tx.Exec(`UPDATE users SET email = ? WHERE name = ?`, "newcharlie@example.com", "Charlie")
	if err != nil {
		tx.Rollback()
		log.Fatalf("事务更新失败: %v", err)
	}

	err = tx.Commit()
	if err != nil {
		log.Fatalf("提交事务失败: %v", err)
	}

	// 使用查询单行数据的方法
	var email string
	err = db.QueryRow(`SELECT email FROM users WHERE name = ?`, "Charlie").Scan(&email)
	if err != nil {
		log.Fatalf("查询单行数据失败: %v", err)
	}
	fmt.Printf("查询到的 Charlie 的邮箱: %s\n", email)
}

```



# kafka

**下载的MinGW64**：x86_64-13.2.0-release-posix-seh-msvcrt-rt_v11-rev0

https://github.com/niXman/mingw-builds-binaries/releases

**框架采用：**"github.com/confluentinc/confluent-kafka-go/v2/kafka"

**开启：** `$env:CGO_ENABLED = "1"`，让go可以引用c的编译。，否则可能一直找不到库

```
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

// 生产者
func produce() {
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "localhost:9092"})
	if err != nil {
		log.Fatalf("无法创建生产者: %v", err)
	}
	defer p.Close()

	topic := "myTopic"

	// 监听生产者的交付报告
	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					fmt.Printf("交付失败: %v\n", ev.TopicPartition)
				} else {
					fmt.Printf("成功交付到 %v\n", ev.TopicPartition)
				}
			}
		}
	}()

	// 发送消息
	for i := 0; i < 10; i++ {
		value := fmt.Sprintf("消息 %d", i)
		p.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
			Value:          []byte(value),
		}, nil)

		// 确保消息被交付
		p.Flush(15 * 1000)
		time.Sleep(time.Second)
	}
}

// 消费者
func consume() {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092",
		"group.id":          "myGroup",
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		log.Fatalf("无法创建消费者: %v", err)
	}
	defer c.Close()

	err = c.SubscribeTopics([]string{"myTopic"}, nil)
	if err != nil {
		log.Fatalf("无法订阅主题: %v", err)
	}

	run := true
	for run {
		select {
		case sig := <-signalChan():
			fmt.Printf("捕获到信号 %v, 退出...\n", sig)
			run = false
		default:
			ev := c.Poll(100)
			switch e := ev.(type) {
			case *kafka.Message:
				fmt.Printf("消费消息: %s: %s\n", e.TopicPartition, string(e.Value))
			case kafka.Error:
				fmt.Fprintf(os.Stderr, "错误: %v\n", e)
				run = false
			}
		}
	}
}

// 主题管理
func manageTopics() {
	adminClient, err := kafka.NewAdminClient(&kafka.ConfigMap{"bootstrap.servers": "localhost:9092"})
	if err != nil {
		log.Fatalf("无法创建管理员客户端: %v", err)
	}
	defer adminClient.Close()

	// 创建主题
	topic := kafka.TopicSpecification{
		Topic:             "myTopic",
		NumPartitions:     1,
		ReplicationFactor: 1,
	}
	results, err := adminClient.CreateTopics(context.TODO(), []kafka.TopicSpecification{topic})
	if err != nil {
		log.Fatalf("创建主题失败: %v", err)
	}
	for _, result := range results {
		fmt.Printf("创建主题 %s 结果: %v\n", result.Topic, result.Error)
	}

	// 列出主题
	metadata, err := adminClient.GetMetadata(nil, true, 5000)
	if err != nil {
		log.Fatalf("获取元数据失败: %v", err)
	}
	for _, t := range metadata.Topics {
		fmt.Printf("主题: %s, 分区: %d\n", t.Topic, len(t.Partitions))
	}

	// 删除主题
	results, err = adminClient.DeleteTopics(context.TODO(), []string{"myTopic"})
	if err != nil {
		log.Fatalf("删除主题失败: %v", err)
	}
	for _, result := range results {
		fmt.Printf("删除主题 %s 结果: %v\n", result.Topic, result.Error)
	}
}

func main() {
	go produce()
	go consume()

	// 管理主题
	manageTopics()

	// 稍等片刻以确保操作完成
	time.Sleep(time.Second * 10)
}

// 捕获系统信号
func signalChan() chan os.Signal {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	return sigChan
}

```



# Redis：

```
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

func main() {
	// 创建一个新的 Redis 客户端
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Redis 服务器地址
		Password: "",               // 没有密码，默认
		DB:       0,                // 默认数据库
	})

	// 检查是否成功连接
	pong, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("无法连接到 Redis: %v", err)
	}
	fmt.Printf("连接成功: %v\n", pong)

	// 设置一个键值对
	err = rdb.Set(ctx, "name", "Go-Redis", 0).Err()
	if err != nil {
		log.Fatalf("无法设置键值对: %v", err)
	}

	// 获取键的值
	val, err := rdb.Get(ctx, "name").Result()
	if err != nil {
		log.Fatalf("无法获取键值: %v", err)
	}
	fmt.Printf("name: %s\n", val)

	// 设置带过期时间的键值对
	err = rdb.Set(ctx, "session_token", "abcd1234", time.Minute).Err()
	if err != nil {
		log.Fatalf("无法设置带过期时间的键值对: %v", err)
	}

	// 检查键是否存在
	exists, err := rdb.Exists(ctx, "name").Result()
	if err != nil {
		log.Fatalf("无法检查键是否存在: %v", err)
	}
	fmt.Printf("name 存在: %d\n", exists)

	// 删除一个键
	err = rdb.Del(ctx, "name").Err()
	if err != nil {
		log.Fatalf("无法删除键: %v", err)
	}

	// 使用哈希（hash）数据类型
	err = rdb.HSet(ctx, "user:1000", map[string]interface{}{
		"name":    "John",
		"age":     30,
		"country": "USA",
	}).Err()
	if err != nil {
		log.Fatalf("无法设置哈希: %v", err)
	}

	// 获取哈希字段的值
	age, err := rdb.HGet(ctx, "user:1000", "age").Result()
	if err != nil {
		log.Fatalf("无法获取哈希字段的值: %v", err)
	}
	fmt.Printf("user:1000 age: %s\n", age)

	// 使用列表（list）数据类型
	err = rdb.RPush(ctx, "tasks", "task1", "task2", "task3").Err()
	if err != nil {
		log.Fatalf("无法添加列表元素: %v", err)
	}

	// 弹出列表的第一个元素
	task, err := rdb.LPop(ctx, "tasks").Result()
	if err != nil {
		log.Fatalf("无法弹出列表元素: %v", err)
	}
	fmt.Printf("弹出任务: %s\n", task)

	// 使用集合（set）数据类型
	err = rdb.SAdd(ctx, "tags", "redis", "database", "golang").Err()
	if err != nil {
		log.Fatalf("无法添加集合元素: %v", err)
	}

	// 获取集合的所有成员
	tags, err := rdb.SMembers(ctx, "tags").Result()
	if err != nil {
		log.Fatalf("无法获取集合成员: %v", err)
	}
	fmt.Printf("tags: %v\n", tags)

	// 使用有序集合（sorted set）数据类型
	err = rdb.ZAdd(ctx, "scores", &redis.Z{
		Score:  100,
		Member: "Alice",
	}, &redis.Z{
		Score:  90,
		Member: "Bob",
	}).Err()
	if err != nil {
		log.Fatalf("无法添加有序集合元素: %v", err)
	}

	// 获取有序集合的成员及其分数
	zscores, err := rdb.ZRangeWithScores(ctx, "scores", 0, -1).Result()
	if err != nil {
		log.Fatalf("无法获取有序集合元素: %v", err)
	}
	fmt.Println("scores:")
	for _, z := range zscores {
		fmt.Printf("  %s: %f\n", z.Member, z.Score)
	}

	// 发布/订阅（Pub/Sub）模式
	pubsub := rdb.Subscribe(ctx, "mychannel")

	// 订阅消息
	go func() {
		ch := pubsub.Channel()
		for msg := range ch {
			fmt.Printf("接收到消息: %s\n", msg.Payload)
		}
	}()

	// 发布消息
	err = rdb.Publish(ctx, "mychannel", "Hello, Redis!").Err()
	if err != nil {
		log.Fatalf("无法发布消息: %v", err)
	}

	// 稍等片刻以确保消息传递
	time.Sleep(time.Second)

	// 关闭 Redis 客户端
	err = rdb.Close()
	if err != nil {
		log.Fatalf("无法关闭 Redis 客户端: %v", err)
	}
}

```

# PostgreSQL

**新建数据库** “mydb”

**替换密码** password

```
package main

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// User 定义用户模型
type User struct {
	ID        uint   `gorm:"primaryKey"`
	Name      string `gorm:"size:255"`
	Age       int
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Orders    []Order        `gorm:"foreignKey:UserID"`
}

// Order 定义订单模型
type Order struct {
	ID        uint `gorm:"primaryKey"`
	UserID    uint
	Product   string
	Price     float64
	CreatedAt time.Time
}

func initDB() *gorm.DB {
	// 替换为你的PostgreSQL连接信息
	dsn := "host=localhost user=postgres password=password dbname=mydb port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("无法连接到数据库: %v", err)
	}
	return db
}

// 自动迁移模型
func autoMigrate(db *gorm.DB) {
	err := db.AutoMigrate(&User{}, &Order{})
	if err != nil {
		log.Fatalf("自动迁移失败: %v", err)
	}
	fmt.Println("模型自动迁移完成")
}

// 插入数据
func createUser(db *gorm.DB, name string, age int) {
	user := User{Name: name, Age: age}
	result := db.Create(&user)
	if result.Error != nil {
		log.Fatalf("无法插入数据: %v", result.Error)
	}
	fmt.Printf("成功插入用户: %s, 年龄: %d, ID: %d\n", user.Name, user.Age, user.ID)
}

// 查询单个用户
func findUser(db *gorm.DB, id uint) {
	var user User
	result := db.First(&user, id)
	if result.Error != nil {
		log.Fatalf("查询用户失败: %v", result.Error)
	}
	fmt.Printf("查询到的用户: ID: %d, Name: %s, Age: %d\n", user.ID, user.Name, user.Age)
}

// 查询多个用户
func findUsers(db *gorm.DB) {
	var users []User
	result := db.Find(&users)
	if result.Error != nil {
		log.Fatalf("查询用户列表失败: %v", result.Error)
	}
	fmt.Println("用户列表:")
	for _, user := range users {
		fmt.Printf("ID: %d, Name: %s, Age: %d\n", user.ID, user.Name, user.Age)
	}
}

// 更新用户
func updateUser(db *gorm.DB, id uint, newAge int) {
	result := db.Model(&User{}).Where("id = ?", id).Update("age", newAge)
	if result.Error != nil {
		log.Fatalf("更新用户失败: %v", result.Error)
	}
	fmt.Printf("成功更新用户 ID %d 的年龄为 %d\n", id, newAge)
}

// 删除用户
func deleteUser(db *gorm.DB, id uint) {
	result := db.Delete(&User{}, id)
	if result.Error != nil {
		log.Fatalf("删除用户失败: %v", result.Error)
	}
	fmt.Printf("成功删除用户 ID %d\n", id)
}

// 创建订单
func createOrder(db *gorm.DB, userID uint, product string, price float64) {
	order := Order{UserID: userID, Product: product, Price: price}
	result := db.Create(&order)
	if result.Error != nil {
		log.Fatalf("创建订单失败: %v", result.Error)
	}
	fmt.Printf("成功创建订单: 用户ID: %d, 产品: %s, 价格: %.2f\n", userID, product, price)
}

// 查询用户的订单
func findUserOrders(db *gorm.DB, userID uint) {
	var orders []Order
	result := db.Where("user_id = ?", userID).Find(&orders)
	if result.Error != nil {
		log.Fatalf("查询订单失败: %v", result.Error)
	}
	fmt.Printf("用户 ID %d 的订单:\n", userID)
	for _, order := range orders {
		fmt.Printf("订单ID: %d, 产品: %s, 价格: %.2f\n", order.ID, order.Product, order.Price)
	}
}

// 事务处理
func transactionalCreate(db *gorm.DB, userName string, userAge int, product string, price float64) {
	err := db.Transaction(func(tx *gorm.DB) error {
		// 创建用户
		user := User{Name: userName, Age: userAge}
		if err := tx.Create(&user).Error; err != nil {
			return err
		}

		// 创建订单
		order := Order{UserID: user.ID, Product: product, Price: price}
		if err := tx.Create(&order).Error; err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		log.Fatalf("事务操作失败: %v", err)
	}
	fmt.Println("事务操作成功")
}

// 使用 GORM 的高级查询方法
func advancedQueries(db *gorm.DB) {
	var users []User

	// 获取所有用户并按年龄降序排序
	db.Order("age desc").Find(&users)
	fmt.Println("按年龄降序排序的用户列表:")
	for _, user := range users {
		fmt.Printf("ID: %d, Name: %s, Age: %d\n", user.ID, user.Name, user.Age)
	}

	// 查询年龄大于 25 的用户
	db.Where("age > ?", 25).Find(&users)
	fmt.Println("年龄大于 25 的用户:")
	for _, user := range users {
		fmt.Printf("ID: %d, Name: %s, Age: %d\n", user.ID, user.Name, user.Age)
	}

	// 使用原生 SQL 查询
	db.Raw("SELECT * FROM users WHERE age > ?", 25).Scan(&users)
	fmt.Println("使用原生 SQL 查询年龄大于 25 的用户:")
	for _, user := range users {
		fmt.Printf("ID: %d, Name: %s, Age: %d\n", user.ID, user.Name, user.Age)
	}
}

func main() {
	db := initDB()

	// 自动迁移
	autoMigrate(db)

	// 插入数据
	createUser(db, "Alice", 30)
	createUser(db, "Bob", 25)

	// 查询数据
	findUser(db, 1)
	findUsers(db)

	// 更新数据
	updateUser(db, 1, 31)

	// 查询更新后的数据
	findUsers(db)

	// 删除数据
	deleteUser(db, 2)

	// 查询删除后的数据
	findUsers(db)

	// 创建订单
	createOrder(db, 1, "Laptop", 1200.50)
	createOrder(db, 1, "Smartphone", 800.00)

	// 查询用户的订单
	findUserOrders(db, 1)

	// 事务处理
	transactionalCreate(db, "Charlie", 35, "Tablet", 600.00)

	// 使用高级查询方法
	advancedQueries(db)
}

```

# doris

```
package main

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// 定义一个模型
type User struct {
	ID        uint      `gorm:"primaryKey"`
	Name      string    `gorm:"size:255"`
	Age       int
	CreatedAt time.Time
	UpdatedAt time.Time
}

func initDB() *gorm.DB {
	// 替换为你的Doris连接信息
	dsn := "username:password@tcp(127.0.0.1:9030)/your_db_name?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("无法连接到数据库: %v", err)
	}
	return db
}

// 自动迁移模型
func autoMigrate(db *gorm.DB) {
	err := db.AutoMigrate(&User{})
	if err != nil {
		log.Fatalf("自动迁移失败: %v", err)
	}
	fmt.Println("模型自动迁移完成")
}

// 插入数据
func createUser(db *gorm.DB, name string, age int) {
	user := User{Name: name, Age: age}
	result := db.Create(&user)
	if result.Error != nil {
		log.Fatalf("无法插入数据: %v", result.Error)
	}
	fmt.Printf("成功插入用户: %s, 年龄: %d, ID: %d\n", user.Name, user.Age, user.ID)
}

// 查询单个用户
func findUser(db *gorm.DB, id uint) {
	var user User
	result := db.First(&user, id)
	if result.Error != nil {
		log.Fatalf("查询用户失败: %v", result.Error)
	}
	fmt.Printf("查询到的用户: ID: %d, Name: %s, Age: %d\n", user.ID, user.Name, user.Age)
}

// 查询多个用户
func findUsers(db *gorm.DB) {
	var users []User
	result := db.Find(&users)
	if result.Error != nil {
		log.Fatalf("查询用户列表失败: %v", result.Error)
	}
	fmt.Println("用户列表:")
	for _, user := range users {
		fmt.Printf("ID: %d, Name: %s, Age: %d\n", user.ID, user.Name, user.Age)
	}
}

// 更新用户
func updateUser(db *gorm.DB, id uint, newAge int) {
	result := db.Model(&User{}).Where("id = ?", id).Update("age", newAge)
	if result.Error != nil {
		log.Fatalf("更新用户失败: %v", result.Error)
	}
	fmt.Printf("成功更新用户 ID %d 的年龄为 %d\n", id, newAge)
}

// 删除用户
func deleteUser(db *gorm.DB, id uint) {
	result := db.Delete(&User{}, id)
	if result.Error != nil {
		log.Fatalf("删除用户失败: %v", result.Error)
	}
	fmt.Printf("成功删除用户 ID %d\n", id)
}

// 使用 GORM 的高级查询方法
func advancedQueries(db *gorm.DB) {
	var users []User

	// 获取所有用户并按年龄降序排序
	db.Order("age desc").Find(&users)
	fmt.Println("按年龄降序排序的用户列表:")
	for _, user := range users {
		fmt.Printf("ID: %d, Name: %s, Age: %d\n", user.ID, user.Name, user.Age)
	}

	// 查询年龄大于 25 的用户
	db.Where("age > ?", 25).Find(&users)
	fmt.Println("年龄大于 25 的用户:")
	for _, user := range users {
		fmt.Printf("ID: %d, Name: %s, Age: %d\n", user.ID, user.Name, user.Age)
	}

	// 使用原生 SQL 查询
	db.Raw("SELECT * FROM users WHERE age > ?", 25).Scan(&users)
	fmt.Println("使用原生 SQL 查询年龄大于 25 的用户:")
	for _, user := range users {
		fmt.Printf("ID: %d, Name: %s, Age: %d\n", user.ID, user.Name, user.Age)
	}
}

func main() {
	db := initDB()

	// 自动迁移
	autoMigrate(db)

	// 插入数据
	createUser(db, "Alice", 30)
	createUser(db, "Bob", 25)

	// 查询数据
	findUser(db, 1)
	findUsers(db)

	// 更新数据
	updateUser(db, 1, 31)

	// 查询更新后的数据
	findUsers(db)

	// 删除数据
	deleteUser(db, 2)

	// 查询删除后的数据
	findUsers(db)

	// 使用高级查询方法
	advancedQueries(db)
}

```



# MongoDB

```
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const (
	mongoURI = "mongodb://localhost:27017"
	dbName   = "exampleDB"
	collName = "users"
)

type User struct {
	ID    string `bson:"_id,omitempty"`
	Name  string `bson:"name"`
	Age   int    `bson:"age"`
	Email string `bson:"email"`
}

func main() {
	// 创建客户端
	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatalf("无法创建MongoDB客户端: %v", err)
	}

	// 连接到MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatalf("无法连接到MongoDB: %v", err)
	}
	defer client.Disconnect(ctx)

	// 选择数据库和集合
	db := client.Database(dbName)
	collection := db.Collection(collName)

	// 插入单个文档
	insertOne(ctx, collection)

	// 插入多个文档
	insertMany(ctx, collection)

	// 查询单个文档
	findOne(ctx, collection)

	// 查询多个文档
	findMany(ctx, collection)

	// 更新单个文档
	updateOne(ctx, collection)

	// 更新多个文档
	updateMany(ctx, collection)

	// 删除单个文档
	deleteOne(ctx, collection)

	// 删除多个文档
	deleteMany(ctx, collection)

	// 创建索引
	createIndexes(ctx, collection)

	// 使用事务
	// runTransaction(ctx, client)

	// 聚合查询
	aggregate(ctx, collection)
}

// 插入单个文档
func insertOne(ctx context.Context, collection *mongo.Collection) {
	user := User{Name: "Alice", Age: 25, Email: "alice@example.com"}
	result, err := collection.InsertOne(ctx, user)
	if err != nil {
		log.Fatalf("插入单个文档失败: %v", err)
	}
	fmt.Printf("插入单个文档的ID: %v\n", result.InsertedID)
}

// 插入多个文档
func insertMany(ctx context.Context, collection *mongo.Collection) {
	users := []interface{}{
		User{Name: "Bob", Age: 30, Email: "bob@example.com"},
		User{Name: "Charlie", Age: 35, Email: "charlie@example.com"},
	}

	result, err := collection.InsertMany(ctx, users)
	if err != nil {
		log.Fatalf("插入多个文档失败: %v", err)
	}
	fmt.Printf("插入多个文档的IDs: %v\n", result.InsertedIDs)
}

// 查询单个文档
func findOne(ctx context.Context, collection *mongo.Collection) {
	var user User
	filter := bson.M{"name": "Alice"}
	err := collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		log.Fatalf("查询单个文档失败: %v", err)
	}
	fmt.Printf("查询到的单个文档: %+v\n", user)
}

// 查询多个文档
func findMany(ctx context.Context, collection *mongo.Collection) {
	filter := bson.M{"age": bson.M{"$gt": 25}}
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		log.Fatalf("查询多个文档失败: %v", err)
	}
	defer cursor.Close(ctx)

	fmt.Println("查询到的多个文档:")
	for cursor.Next(ctx) {
		var user User
		err := cursor.Decode(&user)
		if err != nil {
			log.Fatalf("解码文档失败: %v", err)
		}
		fmt.Printf("%+v\n", user)
	}
}

// 更新单个文档
func updateOne(ctx context.Context, collection *mongo.Collection) {
	filter := bson.M{"name": "Alice"}
	update := bson.M{"$set": bson.M{"age": 26}}
	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Fatalf("更新单个文档失败: %v", err)
	}
	fmt.Printf("更新单个文档的匹配数: %v, 修改数: %v\n", result.MatchedCount, result.ModifiedCount)
}

// 更新多个文档
func updateMany(ctx context.Context, collection *mongo.Collection) {
	filter := bson.M{"age": bson.M{"$gt": 30}}
	update := bson.M{"$set": bson.M{"email": "updated@example.com"}}
	result, err := collection.UpdateMany(ctx, filter, update)
	if err != nil {
		log.Fatalf("更新多个文档失败: %v", err)
	}
	fmt.Printf("更新多个文档的匹配数: %v, 修改数: %v\n", result.MatchedCount, result.ModifiedCount)
}

// 删除单个文档
func deleteOne(ctx context.Context, collection *mongo.Collection) {
	filter := bson.M{"name": "Alice"}
	result, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		log.Fatalf("删除单个文档失败: %v", err)
	}
	fmt.Printf("删除单个文档的删除数: %v\n", result.DeletedCount)
}

// 删除多个文档
func deleteMany(ctx context.Context, collection *mongo.Collection) {
	filter := bson.M{"age": bson.M{"$lt": 35}}
	result, err := collection.DeleteMany(ctx, filter)
	if err != nil {
		log.Fatalf("删除多个文档失败: %v", err)
	}
	fmt.Printf("删除多个文档的删除数: %v\n", result.DeletedCount)
}

// 创建索引
func createIndexes(ctx context.Context, collection *mongo.Collection) {
	models := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "email", Value: 1}}, // bson.D 表示有序的 key-value 对
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{
				{Key: "name", Value: 1},
				{Key: "age", Value: -1},
			},
		},
	}

	indexes, err := collection.Indexes().CreateMany(ctx, models)
	if err != nil {
		log.Fatalf("创建索引失败: %v", err)
	}
	fmt.Printf("创建的索引: %v\n", indexes)
}

// 运行事务
// 事务功能要求 MongoDB 部署在副本集（Replica Set）或分片集群（Sharded Cluster）中，
// 而不是单节点的独立 MongoDB 实例。
// func runTransaction(ctx context.Context, client *mongo.Client) {
// 	session, err := client.StartSession()
// 	if err != nil {
// 		log.Fatalf("启动会话失败: %v", err)
// 	}
// 	defer session.EndSession(ctx)

// 	callback := func(sessCtx mongo.SessionContext) (interface{}, error) {
// 		collection := client.Database(dbName).Collection(collName)

// 		// 在事务中执行插入操作
// 		user := User{Name: "Dave", Age: 40, Email: "dave@example.com"}
// 		_, err := collection.InsertOne(sessCtx, user)
// 		if err != nil {
// 			return nil, err
// 		}

// 		// 在事务中执行更新操作
// 		filter := bson.M{"name": "Bob"}
// 		update := bson.M{"$set": bson.M{"age": 31}}
// 		_, err = collection.UpdateOne(sessCtx, filter, update)
// 		if err != nil {
// 			return nil, err
// 		}

// 		return nil, nil
// 	}

// 	_, err = session.WithTransaction(ctx, callback, options.Transaction().
// 		SetWriteConcern(&writeconcern.WriteConcern{W: writeconcern.Majority()}).
// 		SetReadConcern(readconcern.Snapshot()).
// 		SetReadPreference(readpref.Primary()))
// 	if err != nil {
// 		log.Fatalf("事务执行失败: %v", err)
// 	}

// 	fmt.Println("事务执行成功")
// }

// 聚合查询
func aggregate(ctx context.Context, collection *mongo.Collection) {
	pipeline := mongo.Pipeline{
		{primitive.E{Key: "$match", Value: bson.M{"age": bson.M{"$gt": 25}}}},
		{primitive.E{Key: "$group", Value: bson.M{"_id": "$age", "count": bson.M{"$sum": 1}}}},
	}

	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		log.Fatalf("聚合查询失败: %v", err)
	}
	defer cursor.Close(ctx)

	fmt.Println("聚合查询结果:")
	for cursor.Next(ctx) {
		var result bson.M
		err := cursor.Decode(&result)
		if err != nil {
			log.Fatalf("解码聚合结果失败: %v", err)
		}
		fmt.Printf("%v\n", result)
	}
}

```



# RabbitMQ

```
package main

import (
	"fmt"
	"log"
	"time"

	"github.com/streadway/amqp"
)

const (
	rabbitMQURL = "amqp://guest:guest@localhost:5672/"
	exchange    = "example_exchange"
	queue       = "example_queue"
	routingKey  = "example_routing_key"
)

func main() {
	// 连接到 RabbitMQ 服务器
	conn, err := amqp.Dial(rabbitMQURL)
	if err != nil {
		log.Fatalf("无法连接到 RabbitMQ: %v", err)
	}
	defer conn.Close()

	// 创建一个通道
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("无法创建通道: %v", err)
	}
	defer ch.Close()

	// 声明交换器
	err = ch.ExchangeDeclare(
		exchange, // 交换器名称
		"direct", // 交换器类型 (direct, fanout, topic, headers)
		true,     // 是否持久化
		false,    // 是否自动删除
		false,    // 是否为内置交换器
		false,    // 是否等待确认
		nil,      // 额外参数
	)
	if err != nil {
		log.Fatalf("无法声明交换器: %v", err)
	}

	// 声明队列
	q, err := ch.QueueDeclare(
		queue, // 队列名称
		true,  // 是否持久化
		false, // 是否自动删除
		false, // 是否为独占队列
		false, // 是否等待确认
		nil,   // 额外参数
	)
	if err != nil {
		log.Fatalf("无法声明队列: %v", err)
	}

	// 绑定队列到交换器
	err = ch.QueueBind(
		q.Name,     // 队列名称
		routingKey, // 路由键
		exchange,   // 交换器名称
		false,      // 是否等待确认
		nil,        // 额外参数
	)
	if err != nil {
		log.Fatalf("无法绑定队列到交换器: %v", err)
	}

	// 发送消息
	body := "Hello, RabbitMQ!"
	err = ch.Publish(
		exchange,   // 交换器名称
		routingKey, // 路由键
		false,      // 是否等待确认
		false,      // 是否等待确认
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	if err != nil {
		log.Fatalf("无法发送消息: %v", err)
	}
	fmt.Printf("发送消息: %s\n", body)

	// 接收消息
	msgs, err := ch.Consume(
		q.Name, // 队列名称
		"",     // 消费者标签 (空字符串为默认值)
		true,   // 是否自动确认
		false,  // 是否独占
		false,  // 是否等待确认
		false,  // 是否等待服务器推送
		nil,    // 额外参数
	)
	if err != nil {
		log.Fatalf("无法注册消费者: %v", err)
	}

	// 使用 goroutine 处理接收的消息
	go func() {
		for d := range msgs {
			fmt.Printf("接收到消息: %s\n", d.Body)
		}
	}()

	// 模拟长时间运行的服务
	fmt.Println("等待接收消息...")
	time.Sleep(10 * time.Second)
}

```





