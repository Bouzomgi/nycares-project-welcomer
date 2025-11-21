package domain

import "time"

type Project struct {
	Name string
	Date time.Time
	Id   string
}
