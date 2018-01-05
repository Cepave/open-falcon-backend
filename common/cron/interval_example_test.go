package cron

import (
	"fmt"
	"regexp"
	"time"
)

func ExampleIntervalService_GetWorkerProfile() {
	newService := NewIntervalService()
	newService.Add(
		"worker-33", &IntervalConfig{FixedDelay: 200 * time.Millisecond},
		ProcJob(func() error {
			return nil
		}),
	)

	newService.Start()
	defer newService.Stop()

	time.Sleep(time.Second)
	profile := newService.GetWorkerProfile("worker-33")

	match, _ := regexp.MatchString("Called: \\d+", fmt.Sprintf("Called: %d", profile.Called))

	fmt.Printf("Match /Called: \\d+/: %v", match)
	// Output:
	// Match /Called: \d+/: true
}

func ExampleIntervalService_Start() {
	var message = ""

	newService := NewIntervalService()
	newService.Add(
		"worker-1", &IntervalConfig{FixedDelay: 500 * time.Millisecond},
		ProcJob(func() error {
			message = "Get Called 1"
			return nil
		}),
	)

	newService.Start()
	defer newService.Stop()

	time.Sleep(time.Second)

	fmt.Println("Get Called 1")
	// Output:
	// Get Called 1
}

func ExampleIntervalService_StartWorker() {
	var message = ""

	newService := NewIntervalService()
	newService.Add(
		"worker-2", &IntervalConfig{FixedDelay: 500 * time.Millisecond},
		ProcJob(func() error {
			message = "Get Called 2"
			return nil
		}),
	)

	newService.StartWorker("worker-2")
	defer newService.StopWorker("worker-2")

	time.Sleep(time.Second)

	fmt.Println(message)
	// Output:
	// Get Called 2
}
