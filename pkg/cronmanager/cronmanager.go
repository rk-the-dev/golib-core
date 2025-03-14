package cronmanager

import (
	"log"
	"reflect"
	"sync"

	"github.com/robfig/cron/v3"
)

type Job struct {
	ID       cron.EntryID
	Schedule string
	Task     interface{} // Generic function
	Args     []interface{}
}

type CronManager struct {
	cron  *cron.Cron
	jobs  map[string]Job
	mutex sync.Mutex
}

// New creates a new CronManager instance
func New() *CronManager {
	return &CronManager{
		cron: cron.New(cron.WithSeconds()),
		jobs: make(map[string]Job),
	}
}

// Start begins the cron scheduler
func (c *CronManager) Start() {
	c.cron.Start()
	log.Println("CronManager started")
}

// Stop halts the cron scheduler
func (c *CronManager) Stop() {
	c.cron.Stop()
	log.Println("CronManager stopped")
}

// AddJob schedules a new job with parameters
func (c *CronManager) AddJob(name, schedule string, task interface{}, args ...interface{}) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	wrappedTask := func() {
		invokeFunction(task, args...)
	}

	id, err := c.cron.AddFunc(schedule, wrappedTask)
	if err != nil {
		return err
	}

	c.jobs[name] = Job{ID: id, Schedule: schedule, Task: task, Args: args}
	log.Printf("Added job: %s, schedule: %s", name, schedule)
	return nil
}

// RemoveJob removes a job by name
func (c *CronManager) RemoveJob(name string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if job, exists := c.jobs[name]; exists {
		c.cron.Remove(job.ID)
		delete(c.jobs, name)
		log.Printf("Removed job: %s", name)
	}
}

// ListJobs returns all scheduled jobs
func (c *CronManager) ListJobs() map[string]Job {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	return c.jobs
}

// invokeFunction dynamically calls a function with parameters
func invokeFunction(fn interface{}, args ...interface{}) {
	funcValue := reflect.ValueOf(fn)
	if funcValue.Kind() != reflect.Func {
		log.Println("Error: Task is not a function")
		return
	}

	var inputArgs []reflect.Value
	for _, arg := range args {
		inputArgs = append(inputArgs, reflect.ValueOf(arg))
	}

	funcValue.Call(inputArgs)
}
