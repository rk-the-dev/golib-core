package main

import (
	"fmt"
	"time"

	"github.com/rk-the-dev/golib-core/pkg/cronmanager"
)

func main() {
	// Create a new CronManager instance
	cm := cronmanager.New()

	// 1️⃣ Simple Job: No Parameters
	cm.AddJob("job1", "*/5 * * * * *", func() {
		fmt.Println("🔹 Simple Job executed at:", time.Now())
	})

	// 2️⃣ Job with One Parameter
	printMessage := func(message string) {
		fmt.Println("📢 Message Job:", message)
	}
	cm.AddJob("job2", "*/10 * * * * *", printMessage, "Hello from cron!")

	// 3️⃣ Job with Multiple Parameters
	processData := func(id int, name string) {
		fmt.Printf("📦 Processing Data: ID=%d, Name=%s\n", id, name)
	}
	cm.AddJob("job3", "*/15 * * * * *", processData, 101, "Kishan")

	// 4️⃣ Job with Struct Parameter
	type TaskDetails struct {
		ID    int
		Owner string
	}
	processTask := func(task TaskDetails) {
		fmt.Printf("📝 Processing Task: ID=%d, Owner=%s\n", task.ID, task.Owner)
	}
	cm.AddJob("job4", "*/20 * * * * *", processTask, TaskDetails{ID: 1, Owner: "Admin"})

	// 5️⃣ Job with Variadic Parameters
	logData := func(params ...interface{}) {
		fmt.Println("📝 Log Data:", params)
	}
	cm.AddJob("job5", "*/25 * * * * *", logData, "Log1", 123, true)

	// Start the cron manager (All jobs will start executing automatically)
	cm.Start()

	// Run indefinitely
	select {}
}
