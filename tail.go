package main

import (
	"bufio"
	"context"
	"fmt"
	"os"

	"github.com/fatih/color"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/pkg/api/v1"
)

var (
	greenBold = color.New(color.FgGreen, color.Bold)
	redBold   = color.New(color.FgRed, color.Bold)
)

type Tail struct {
	closed       chan struct{}
	logger       *Logger
	namespace    string
	pod          string
	container    string
	sinceSeconds int64
	timestamps   bool
}

// NewTail creates new Tail object
func NewTail(namespace, pod, container string, logger *Logger, sinceSeconds int64, timestamps bool) *Tail {
	return &Tail{
		closed:       make(chan struct{}),
		logger:       logger,
		namespace:    namespace,
		pod:          pod,
		container:    container,
		sinceSeconds: sinceSeconds,
		timestamps:   timestamps,
	}
}

// Start starts Pod log streaming
func (t *Tail) Start(ctx context.Context, clientset *kubernetes.Clientset) {
	t.logger.PrintColorizedLog(greenBold, fmt.Sprintf("Pod:%s Container:%s has been detected", t.pod, t.container))

	go func() {
		rs, err := clientset.Core().Pods(t.namespace).GetLogs(t.pod, &v1.PodLogOptions{
			Container:    t.container,
			Follow:       true,
			SinceSeconds: &t.sinceSeconds,
			Timestamps:   t.timestamps,
		}).Stream()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}
		defer rs.Close()

		go func() {
			<-t.closed
			rs.Close()
		}()

		sc := bufio.NewScanner(rs)

		for sc.Scan() {
			t.logger.PrintPodLog(t.pod, t.container, sc.Text(), t.timestamps)
		}
	}()

	go func() {
		<-ctx.Done()
		close(t.closed)
	}()
}

// Stop finishes Pod log streaming
func (t *Tail) Stop() {
	t.logger.PrintColorizedLog(redBold, fmt.Sprintf("Pod:%s Container:%s has been deleted", t.pod, t.container))
	close(t.closed)
}
