package main

import (
	"flag"
	"io"
	"os"

	log "github.com/sirupsen/logrus"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

var debug bool

func init() {
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.JSONFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(io.MultiWriter(os.Stdout, &lumberjack.Logger{
		Filename:  "log/log.txt",
		MaxSize:   10, // 10MB
		LocalTime: true,
		Compress:  true,
	}))

	flag.BoolVar(&debug, "debug", false, "Run in debug mode")

	// Info 로그부터 출력함
	log.SetLevel(log.InfoLevel)
}

func main() {
	flag.Parse()

	if debug {
		// Debug 로그부터 출력함
		log.SetLevel(log.DebugLevel)
	}

	log.Info("Starting logging")

	for i := 0; i < 10000000; i++ {
		log.Debug("Useful debugging information.")
		log.Info("Something noteworthy happened!")
		log.Warn("You should probably take a look at this.")
	}

	log.Info("Finished logging")
}
