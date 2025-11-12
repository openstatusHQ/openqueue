package database

type Task struct {
	ID          int64  `json:"id" db:"id"`
	Method      string `json:"method" db:"method"`
	Headers     string `json:"headers" db:"headers"`
	Body        string `json:"body,omitempty" db:"body"`
	URL         string `json:"url,omitempty" db:"url"`
	CreatedAt   int64  `json:"created_at" db:"created_at"`
	ScheduledAt int64  `json:"scheduled_at" db:"scheduled_at"`
	UpdatedAt   int64  `json:"updated_at" db:"updated_at"`
}

type TaskStatus string

const (
	StatusPending   TaskStatus = "pending"
	StatusRunning   TaskStatus = "running"
	StatusCompleted TaskStatus = "completed"
	StatusFailed    TaskStatus = "failed"
)

type TaskExecution struct {
	Status    TaskStatus `json:"status" db:"status"`
	ID        int64      `json:"id" db:"id"`
	TaskId    int64      `json:"task_id" db:"task_id"`
	CreatedAt int64      `json:"created_at" db:"created_at"`
	UpdatedAt int64      `json:"updated_at" db:"updated_at"`
}
