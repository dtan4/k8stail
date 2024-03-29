package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"sync"

	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
)

type Tail struct {
	Finished     bool
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
		Finished:     false,
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
	t.logger.PrintPodDetected(t.pod, t.container)

	go func() {
		rs, err := clientset.CoreV1().Pods(t.namespace).GetLogs(t.pod, &v1.PodLogOptions{
			Container:    t.container,
			Follow:       true,
			SinceSeconds: &t.sinceSeconds,
			Timestamps:   t.timestamps,
		}).Stream(ctx)
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

// Finish finishes Pod log streaming with Pod completion
func (t *Tail) Finish() {
	t.logger.PrintPodFinished(t.pod, t.container)
	t.Finished = true
}

// Delete finishes Pod log streaming with Pod deletion
func (t *Tail) Delete() {
	t.logger.PrintPodDeleted(t.pod, t.container)
	close(t.closed)
}

type TailMap struct {
	mu sync.Mutex

	data map[string]*Tail
}

func NewTailMap() *TailMap {
	return &TailMap{
		data: make(map[string]*Tail),
	}
}

func (m *TailMap) Get(k string) (*Tail, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	d, ok := m.data[k]

	return d, ok
}

func (m *TailMap) Set(k string, v *Tail) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.data[k] = v
}

func (m *TailMap) Delete(k string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.data, k)
}
