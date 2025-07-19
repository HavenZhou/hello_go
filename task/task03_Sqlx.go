package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

func main03() {
	db := getConnect()
	fmt.Printf("连接数据库db=%v,", db)
	createTables(db)

	initDataBase(db)

	// Sqlx入门
	// 题目1：使用SQL扩展库进行查询
	// 假设你已经使用Sqlx连接到一个数据库，并且有一个 employees 表，包含字段 id 、 name 、 department 、 salary 。
	// 要求 ：
	// 编写Go代码，使用Sqlx查询 employees 表中所有部门为 "技术部" 的员工信息，并将结果映射到一个自定义的 Employee 结构体切片中。
	// 编写Go代码，使用Sqlx查询 employees 表中工资最高的员工信息，并将结果映射到一个 Employee 结构体中。
	var employees []Employees
	err := db.Select(&employees, "SELECT * FROM employees where Department = ? ORDER BY salary DESC", "技术部")
	if err != nil {
		log.Fatal("报错=", err)
	}

	employeesJson, _ := json.MarshalIndent(employees, "", "  ")
	fmt.Println(string(employeesJson))

	var employee Employees
	err1 := db.Get(&employee, "SELECT * FROM employees ORDER BY salary DESC LIMIT 1")
	if err1 != nil {
		log.Fatal("报错=", err)
	}
	employeeJson, _ := json.MarshalIndent(employee, "", "  ")
	fmt.Println("工资最高的员工信息=", string(employeeJson))

	// 题目2：实现类型安全映射
	// 假设有一个 books 表，包含字段 id 、 title 、 author 、 price 。
	// 要求 ：
	// 定义一个 Book 结构体，包含与 books 表对应的字段。
	// 编写Go代码，使用Sqlx执行一个复杂的查询，例如查询价格大于 50 元的书籍，并将结果映射到 Book 结构体切片中，确保类型安全。
	var expensiveBooks []Book
	err = db.Select(&expensiveBooks, `
		SELECT id, title, author, price 
		FROM books 
		WHERE price > ? 
		ORDER BY price DESC
	`, 50)
	if err != nil {
		log.Fatal("查询高价书籍失败:", err)
	}
	booksJson, _ := json.MarshalIndent(expensiveBooks, "", "  ")
	fmt.Println("\n价格大于50元的书籍:")
	fmt.Println(string(booksJson))
}

// 数据库连接
func getConnect() *sqlx.DB {
	// 连接MYSQL
	db, err := sqlx.Connect("mysql", "root:147258@tcp(localhost:3306)/hello_sqlx?parseTime=true")
	if err != nil {
		log.Fatal(err)
	}
	//defer db.Close()
	return db
}

type Employees struct {
	ID         int       `db:"id"`
	Name       string    `db:"name"`
	Department string    `db:"department"`
	Salary     float64   `db:"salary"`
	Created_at time.Time `db:"created_at"`
	Cpdated_at time.Time `db:"updated_at"`
}

type Book struct {
	ID     int64   `db:"id"`
	Title  string  `db:"title"`
	Author string  `db:"author"`
	Price  float64 `db:"price"`
}

// 创建表
// createTables 创建所有需要的表
func createTables(db *sqlx.DB) error {
	// 使用事务确保所有表都能成功创建
	tx, err := db.Beginx()
	if err != nil {
		return fmt.Errorf("could not begin transaction: %v", err)
	}

	// 如果中途出错，回滚事务
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// 创建用户表
	if _, err := tx.Exec(`
		CREATE TABLE IF NOT EXISTS employees (
			id BIGINT PRIMARY KEY AUTO_INCREMENT,
			name VARCHAR(255) UNIQUE NOT NULL,
			department VARCHAR(255)  NOT NULL,
			salary DECIMAL(10,2) NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
	`); err != nil {
		return fmt.Errorf("could not create users table: %v", err)
	}

	fmt.Println("建表成功")

	// 创建书籍表（题目2）
	if _, err := tx.Exec(`
		CREATE TABLE IF NOT EXISTS books (
			id BIGINT PRIMARY KEY AUTO_INCREMENT,
			title VARCHAR(255) NOT NULL,
			author VARCHAR(255) NOT NULL,
			price DECIMAL(10,2) NOT NULL,
			INDEX idx_price (price)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
	`); err != nil {
		return fmt.Errorf("could not create books table: %v", err)
	}

	fmt.Println("建表成功")
	// 提交事务
	return tx.Commit()
}

// 数据初始化
func initDataBase(db *sqlx.DB) {

	// 使用事务保证原子性
	tx, err := db.Beginx()
	if err != nil {
		log.Fatal("开始事务失败:", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			log.Fatal("初始化数据失败，已回滚:", err)
		} else {
			tx.Commit()
			fmt.Println("数据初始化成功")
		}
	}()

	// 插入数据 - 每个VALUES中的?数量必须与列数匹配
	employees := []struct {
		name       string
		department string
		salary     float64
	}{
		{"Alice", "技术部", 15000.50},
		{"Bob", "技术部", 1000.50},
		{"Charlie", "项目部", 9000.50},
		{"David", "营业部", 8000.50},
	}

	for _, emp := range employees {
		_, err = tx.Exec(
			"INSERT INTO employees (name, department, salary) VALUES (?, ?, ?)",
			emp.name, emp.department, emp.salary,
		)
		if err != nil {
			return
		}
	}

	// 插入书籍数据（题目2）
	books := []struct {
		title  string
		author string
		price  float64
	}{
		{"Go语言编程", "Alan A. Donovan", 89.00},
		{"数据库系统概念", "Abraham Silberschatz", 99.00},
		{"算法导论", "Thomas H. Cormen", 129.00},
		{"计算机网络", "Andrew S. Tanenbaum", 79.00},
		{"代码大全", "Steve McConnell", 59.00},
		{"设计模式", "Erich Gamma", 49.00}, // 这个不会被查询到
	}

	for _, book := range books {
		_, err = tx.Exec(
			"INSERT INTO books (title, author, price) VALUES (?, ?, ?)",
			book.title, book.author, book.price,
		)
		if err != nil {
			return
		}
	}
}
