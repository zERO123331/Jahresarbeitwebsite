package main

import "time"

type Update struct {
	Title       string
	Author      string
	Body        string
	ID          int
	Created     time.Time
	LastUpdated time.Time
}
