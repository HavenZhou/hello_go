package main

import (
	"fmt"
	"sync"
	"time"
)

// 定义可执行任务
type Task func()

// TaskResult 记录任务执行结果
type TaskResult struct {
	TaskId    int
	StartTime time.Time
	EndTime   time.Time
	Duration  time.Duration
}

// Scheduler 任务调度器
type Scheduler struct {
	tasks     []Task
	results   []TaskResult
	wg        sync.WaitGroup
	resultMu  sync.Mutex
	startTime time.Time
}

// NewScheduler 创建新的调度器实例
func NewScheduler() *Scheduler {
	return &Scheduler{
		tasks:   make([]Task, 0),
		results: make([]TaskResult, 0),
	}
}

// AddTask 添加单个任务到调度器
func (s *Scheduler) AddTask(task Task) {
	s.tasks = append(s.tasks, task)
}

// AddTask 添加一组任务到调度器
func (s *Scheduler) AddTaskMore(taskm []Task) {
	s.tasks = append(s.tasks, taskm...)
}

// Run 并发执行所有任务
func (s *Scheduler) Run() {
	s.startTime = time.Now()
	s.results = make([]TaskResult, len(s.tasks))
	s.wg.Add(len(s.tasks))

	for i, task := range s.tasks {
		go s.executeTask(i, task)
	}
	s.wg.Wait()
}

// 执行单个任务并记录结果
func (s *Scheduler) executeTask(taskId int, task Task) {
	defer s.wg.Done()
	start := time.Now()
	reslut := TaskResult{
		TaskId:    taskId,
		StartTime: start,
	}

	defer func() {
		end := time.Now()
		reslut.EndTime = end
		reslut.Duration = end.Sub(start)

		s.resultMu.Lock()
		s.results[taskId] = reslut
		s.resultMu.Unlock()
	}()

	task()
}

// PrintStats 打印统计信息
func (s *Scheduler) PrintStats() {
	fmt.Println("\n任务执行统计:")
	fmt.Printf("总任务数:%d\n", len(s.tasks))
	fmt.Printf("总执行时间：%v", time.Since(s.startTime))

	for _, result := range s.results {
		fmt.Printf("%v\n", result)
	}
}

func main23() {
	scheduler := NewScheduler()

	scheduler.AddTask(func() {
		fmt.Println("任务1开始执行")
		time.Sleep(10 * time.Second)
		fmt.Println("任务1完成")
	})

	scheduler.AddTask(func() {
		fmt.Println("任务2开始执行")
		time.Sleep(10 * time.Second)
		fmt.Println("任务2完成")
	})

	scheduler.AddTask(func() {
		fmt.Println("任务3开始执行")
		time.Sleep(10 * time.Second)
		fmt.Println("任务3完成")
	})

	scheduler.Run()
	scheduler.PrintStats()

	fmt.Println("-------------添加组任务-----------------")
	tasks := []Task{
		func() {
			fmt.Println("任务1m")
			time.Sleep(5 * time.Second)
			fmt.Println("任务1m完成")
		},
		func() {
			fmt.Println("任务2m")
			time.Sleep(5 * time.Second)
			fmt.Println("任务2m完成")
		},
		func() {
			fmt.Println("任务3m")
			time.Sleep(5 * time.Second)
			fmt.Println("任务3m完成")
		},
	}
	scheduler.AddTaskMore(tasks)
	scheduler.Run()
	scheduler.PrintStats()

}
