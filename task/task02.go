package main

import (
	"fmt"
	"math"
	"sync"
	"sync/atomic"
)

func main2() {
	// 题目1：编写一个Go程序，定义一个函数，该函数接收一个整数指针作为参数，在函数内部将该指针指向的值增加10，然后在主函数中调用该函数并输出修改后的值。
	fmt.Println("------------------question 1------------------")
	num := 20
	pointTest(&num)
	fmt.Println("after pointTest num = ", num)

	// 题目2：编写一个程序，使用 go 关键字启动两个协程，一个协程打印从1到10的奇数，另一个协程打印从2到10的偶数。
	fmt.Println("------------------question 2------------------")
	goroutineTest01()
	fmt.Println("所有奇偶数打印完成")

	// 题目3 ：设计一个任务调度器，接收一组任务（可以用函数表示），并使用协程并发执行这些任务，同时统计每个任务的执行时间。
	// 详见 task02_q3.go

	// 题目4：定义一个 Shape 接口，包含 Area() 和 Perimeter() 两个方法。
	// 然后创建 Rectangle 和 Circle 结构体，实现 Shape 接口。在主函数中，创建这两个结构体的实例，并调用它们的 Area() 和 Perimeter() 方法。
	fmt.Println("------------------question 4------------------")
	rect := Rectangle{Width: 5, Height: 3.22}
	circle := Circle{Radius: 0.55}
	PrintInfo(rect)
	PrintInfo(circle)

	// 题目5：使用组合的方式创建一个 Person 结构体，包含 Name 和 Age 字段，
	// 再创建一个 Employee 结构体，组合 Person 结构体并添加 EmployeeID 字段。为 Employee 结构体实现一个 PrintInfo() 方法，输出员工的信息。
	fmt.Println("------------------question 5------------------")
	em := Employee{
		Person:     Person{Name: "张三", Age: 18},
		EmployeeID: "213123123123",
	}
	em.PrintInfo()
	em.PrintInfo()

	// 题目6 ：编写一个程序，使用通道实现两个协程之间的通信。一个协程生成从1到10的整数，并将这些整数发送到通道中，另一个协程从通道中接收这些整数并打印出来。
	fmt.Println("------------------question 6------------------")
	RunChannelCommunication()

	// 题目7：实现一个带有缓冲的通道，生产者协程向通道中发送100个整数，消费者协程从通道中接收这些整数并打印。
	// 省略，和6一致，主要区别为有缓存通道同步通信，一端未准备好另一端会阻塞
	// 无缓存通道异步通信，缓冲满时发送会阻塞，缓冲空时接收会阻塞

	// 题目8：编写一个程序，使用 sync.Mutex 来保护一个共享的计数器。启动10个协程，每个协程对计数器进行1000次递增操作，最后输出计数器的值。
	fmt.Println("------------------question 8------------------")
	MutexTest()

	// 题目9：使用原子操作（ sync/atomic 包）实现一个无锁的计数器。启动10个协程，每个协程对计数器进行1000次递增操作，最后输出计数器的值。
	atomicTest()
}

// 题目1：编写一个Go程序，定义一个函数，该函数接收一个整数指针作为参数，在函数内部将该指针指向的值增加10，然后在主函数中调用该函数并输出修改后的值。
func pointTest(digital *int) {
	*digital += 10
}

// 题目2：编写一个程序，使用 go 关键字启动两个协程，一个协程打印从1到10的奇数，另一个协程打印从2到10的偶数。
// 备注：普通函数中的协程独立于函数调用周期，
func goroutineTest01() {
	var wg sync.WaitGroup
	wg.Add(2)

	// 打印奇数的协程
	go func() {
		defer wg.Done()
		for i := 1; i <= 10; i += 2 {
			fmt.Print("奇数:", i, ";")
			//time.Sleep(1 * time.Second)
		}
	}()

	// 打印偶数的协程
	go func() {
		defer wg.Done()
		for i := 2; i <= 10; i += 2 {
			fmt.Print("偶数：", i, ";")
			//time.Sleep(1 * time.Second)
		}
	}()

	wg.Wait()
}

// 题目4：定义一个 Shape 接口，包含 Area() 和 Perimeter() 两个方法。
// 然后创建 Rectangle 和 Circle 结构体，实现 Shape 接口。在主函数中，创建这两个结构体的实例，并调用它们的 Area() 和 Perimeter() 方法。
type Shape interface {
	Area() float64
	Perimeter() float64
}

type Rectangle struct {
	Width  float64
	Height float64
}

func (r Rectangle) Area() float64 {
	return r.Width * r.Height
}

func (r Rectangle) Perimeter() float64 {
	return 2 * (r.Width + r.Height)
}

type Circle struct {
	Radius float64
}

func (c Circle) Area() float64 {
	return math.Pi * c.Radius * c.Radius
}

func (c Circle) Perimeter() float64 {
	return 2 * math.Pi * c.Radius
}

func PrintInfo(s Shape) {
	fmt.Printf("面积: %.2f\n", s.Area())
	fmt.Printf("周长: %.2f\n", s.Perimeter())
	fmt.Println("------------------")
}

// 题目5：使用组合的方式创建一个 Person 结构体，包含 Name 和 Age 字段，
// 再创建一个 Employee 结构体，组合 Person 结构体并添加 EmployeeID 字段。为 Employee 结构体实现一个 PrintInfo() 方法，输出员工的信息。
type Person struct {
	Name string
	Age  int
}

type Employee struct {
	Person     // 匿名嵌入，可以直接访问字段
	EmployeeID string
}

func (e *Employee) PrintInfo() {
	fmt.Printf("员工信息:\n")
	fmt.Printf("姓名: %s\n", e.Name)       // 可以直接访问嵌入结构的字段
	fmt.Printf("年龄: %d\n", e.Age)        // 同上
	fmt.Printf("工号: %s\n", e.EmployeeID) // Employee 自身的字段
	e.Age = 22
}

// 题目6 ：编写一个程序，使用通道实现两个协程之间的通信。一个协程生成从1到10的整数，并将这些整数发送到通道中，另一个协程从通道中接收这些整数并打印出来。
func RunChannelCommunication() {
	ch := make(chan int, 5)

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done() // 表示一个 goroutine 已经完成，相当于将 WaitGroup 的计数器减 1,实际上是 wg.Add(-1) 的简化形式
		defer close(ch)

		for i := 1; i <= 10; i++ {
			ch <- i
			fmt.Printf("生产者发送：%d\n", i)
		}
	}()

	go func() {
		defer wg.Done()

		for num := range ch {
			//time.Sleep(1 * time.Second)
			fmt.Printf("消费者接收：%d\n", num)
		}
	}()

	wg.Wait() // 阻塞当前 goroutine，直到 WaitGroup 的计数器归零（即所有 goroutine 都执行完毕）
	fmt.Println("通信完成！")
}

// 题目8：编写一个程序，使用 sync.Mutex 来保护一个共享的计数器。启动10个协程，每个协程对计数器进行1000次递增操作，最后输出计数器的值。
func MutexTest() {
	var mutex sync.Mutex
	var wg sync.WaitGroup

	num := 0

	// 外层启动10个协程
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			// 内层开始计数
			for j := 0; j < 1000; j++ {
				mutex.Lock()
				num++
				mutex.Unlock()
			}
		}()
	}
	wg.Wait()
	fmt.Printf("mutex最终相加后num = %d\n", num)
}

// 题目9：使用原子操作（ sync/atomic 包）实现一个无锁的计数器。启动10个协程，每个协程对计数器进行1000次递增操作，最后输出计数器的值。
func atomicTest() {
	// var mutex sync.Mutex
	var wg sync.WaitGroup
	num := int64(0)

	// 外层启动10个协程
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			// 内层开始计数
			for j := 0; j < 1000; j++ {
				//mutex.Lock()
				//num++
				//mutex.Unlock()
				atomic.AddInt64(&num, 1)
			}
		}()
	}
	wg.Wait()
	fmt.Printf("atomic最终相加后num = %d\n", num)
}
