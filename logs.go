package main

import (
	"log"
	"os"
	"strings"
)

var (
	logdir = os.Getenv("HOME") + "/.sqlpeek/logs/"
	// logdir = "./"
	// basedir = "/Users/broaddus/Desktop/go/sink/"
	// basedir = "/home/colemanb/src/notes/go/sink/"
)

var (
	WarningLog *log.Logger
	InfoLog    *log.Logger
	ErrorLog   *log.Logger
)

func init() {
	err := os.MkdirAll(logdir, os.ModePerm)
	check(err)

	log.SetFlags(log.Lshortfile)
	// Get the executable name
	exeName := os.Args[0]

	var file *os.File
	// fmt.Println("ExeName = ", exeName)
	// fmt.Println("os.args = ", os.Args)
	if strings.Contains(exeName, "go-build") {
		// fmt.Println("Running via `go run`")
		file, err = os.OpenFile(logdir+"dev.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		check(err)
	} else {
		// fmt.Println("Running from a built binary")
		file, err = os.OpenFile(logdir+"prod.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		check(err)
	}

	InfoLog = log.New(file, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	WarningLog = log.New(file, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLog = log.New(file, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func check(err error) {
	if err != nil {
		log.Println(err)
		panic(err)
	}
}
