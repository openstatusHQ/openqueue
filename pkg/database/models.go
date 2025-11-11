package database

type Job struct {
    ID          int64  `json:"id" db:"id"`
    Method      string `json:"method" db:"method"`
    Headers     string `json:"headers" db:"headers"`
    Body        string `json:"body,omitempty" db:"body"`
    URL         string `json:"url,omitempty" db:"url"`
    CreatedAt   int64  `json:"created_at" db:"created_at"`
    ScheduledAt int64  `json:"scheduled_at" db:"scheduled_at"`
    UpdatedAt   int64  `json:"updated_at" db:"updated_at"`

}
