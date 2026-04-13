package main

import "time"

var defaultValues = map[string]any{
	"app.name":    "my.ms-todo",
	"app.runmode": "prod",

	"log.min_level": "info",
	"log.format":    "auto",

	"entry.http.write_timeout":   20 * time.Second,
	"entry.http.idle_timeout":    90 * time.Second,
	"entry.http.read_timeout":    10 * time.Second,
	"entry.http.request_timeout": 15 * time.Second,

	"db.max_conns": 10,
	"db.min_conns": 5,

	"db.max_conn_lifetime":  5 * time.Minute,
	"db.max_conn_idle_time": 30 * time.Second,
}
