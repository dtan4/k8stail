package main

import (
	"fmt"
	"strings"
	"sync"

	"github.com/fatih/color"
)

var (
	greenBold  = color.New(color.FgGreen, color.Bold)
	yellowBold = color.New(color.FgYellow, color.Bold)
	redBold    = color.New(color.FgRed, color.Bold)
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

// PrintHeader prints header
func (l *Logger) PrintHeader(context, namespace, labels string) {
	fmt.Printf("%s %s\n", boldFunc("Context:  "), context)
	fmt.Printf("%s %s\n", boldFunc("Namespace:"), namespace)
	fmt.Printf("%s %s\n", boldFunc("Labels:   "), labels)
	color.New(color.FgYellow).Println("Press Ctrl-C to exit.")
	color.New(color.Bold).Println("----------")
}

// PrintPlainLog prints log with no cosmetics
func (l *Logger) PrintPlainLog(line string) {
	l.m.Lock()
	defer l.m.Unlock()

	fmt.Println(line)
}

// PrintPodDetected prints that Pod was detected
func (l *Logger) PrintPodDetected(pod, container string) {
	l.PrintColorizedLog(greenBold, fmt.Sprintf("Pod:%s Container:%s has been detected", pod, container))
}

// PrintPodDeleted prints that Pod was finished
func (l *Logger) PrintPodFinished(pod, container string) {
	l.PrintColorizedLog(yellowBold, fmt.Sprintf("Pod:%s Container:%s has been finished", pod, container))
}

// PrintPodDeleted prints that Pod was deleted
func (l *Logger) PrintPodDeleted(pod, container string) {
	l.PrintColorizedLog(redBold, fmt.Sprintf("Pod:%s Container:%s has been deleted", pod, container))
}

// PrintPodLog prints Pod log
func (l *Logger) PrintPodLog(pod, container, line string, timestamps bool) {
	l.m.Lock()
	defer l.m.Unlock()

	if timestamps {
		ss := strings.SplitN(line, " ", 2)
		fmt.Printf("[%s][%s] %s  %s %s \n", boldFunc(pod), boldFunc(container), yellowFunc(ss[0]), boldFunc("|"), ss[1])
	} else {
		fmt.Printf("[%s][%s]  %s %s\n", boldFunc(pod), boldFunc(container), boldFunc("|"), line)
	}
}
