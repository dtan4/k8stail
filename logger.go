package main

import (
	"fmt"
	"strings"
	"sync"

	"github.com/fatih/color"
)

var (
	boldFunc   = color.New(color.Bold).SprintFunc()
	yellowFunc = color.New(color.FgYellow).SprintFunc()
)

// Logger represents logger
type Logger struct {
	m sync.Mutex
}

// NewLogger returns new Logger object
func NewLogger() *Logger {
	return &Logger{
		m: sync.Mutex{},
	}
}

// PrintColorizedLog prints log with the given color
func (l *Logger) PrintColorizedLog(c *color.Color, line string) {
	l.m.Lock()
	defer l.m.Unlock()

	c.Println(line)
}

// PrintPlainLog prints log with no cosmetics
func (l *Logger) PrintPlainLog(line string) {
	l.m.Lock()
	defer l.m.Unlock()

	fmt.Println(line)
}

// PrintPodLog prints Pod log
func (l *Logger) PrintPodLog(podName, line string, timestamps bool) {
	l.m.Lock()
	defer l.m.Unlock()

	if timestamps {
		ss := strings.SplitN(line, " ", 2)
		fmt.Printf("[%s] %s  %s %s \n", boldFunc(podName), yellowFunc(ss[0]), boldFunc("|"), ss[1])
	} else {
		fmt.Printf("[%s]  %s %s\n", boldFunc(podName), boldFunc("|"), line)
	}
}
