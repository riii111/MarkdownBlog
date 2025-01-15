package health

import "time"

type Response struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
}
