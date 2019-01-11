package main

import (
	"flag"
	"io"
	"log"
	"os"
	"time"
)

var logToFile = flag.Bool("log", false, "log to redshift.log")

func init() {
	if !flag.Parsed() {
		flag.Parse()
	}

	if *logToFile {
		if f, err := os.OpenFile("redshift.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644); err != nil {
			log.Println("error opening log file (redshift.log):", err)
		} else {
			if stat, err := f.Stat(); err == nil && stat.Size() != 0 {
				f.WriteString("\n\n")
			}

			f.WriteString("Server Started " + time.Now().Format(time.RFC3339) + "\n\n")

			log.SetOutput(io.MultiWriter(os.Stdout, f))
		}
	}
}
