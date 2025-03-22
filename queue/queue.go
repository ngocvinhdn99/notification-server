package queue

import (
	"fmt"
	"time"
)

type Job struct {
	ID      int
	Message string
}

type Queue struct {
	jobs []Job
}

func (q *Queue) Enqueue(job Job) {
	q.jobs = append(q.jobs, job)
}

func (q *Queue) Dequeue() Job {
	if len(q.jobs) == 0 {
		return Job{}
	}
	job := q.jobs[0]
	q.jobs = q.jobs[1:]
	return job
}

func (q *Queue) IsEmpty() bool {
	return len(q.jobs) == 0
}

func sendEmail(job Job) {
	fmt.Printf("Sending email: %s\n", job.Message)
	time.Sleep(2 * time.Second) // Giả lập thời gian gửi email
	fmt.Printf("Email sent: %s\n", job.Message)
}

func RunQueue() {
	queue := &Queue{}

	// Đưa công việc vào hàng đợi
	queue.Enqueue(Job{ID: 1, Message: "Welcome to our service!"})
	queue.Enqueue(Job{ID: 2, Message: "Your account has been activated."})
	queue.Enqueue(Job{ID: 3, Message: "Your password reset link."})

	// Xử lý công việc trong hàng đợi
	for !queue.IsEmpty() {
		job := queue.Dequeue()
		go sendEmail(job) // Xử lý công việc gửi email trong một goroutine
	}

	// Đợi cho tất cả email được gửi
	time.Sleep(3 * time.Second)
}
