package main

import (
	"sync"
)

// PodList represents the goroutine-safe map contains Pod names
type PodList struct {
	pods map[string]bool
	m    sync.Mutex
}

// NewPodList returns new PodList object
func NewPodList() *PodList {
	return &PodList{
		pods: map[string]bool{},
		m:    sync.Mutex{},
	}
}

// Add adds the given pod
func (p *PodList) Add(podName string) {
	p.m.Lock()
	defer p.m.Unlock()

	p.pods[podName] = true
}

// Delete delete the given pod
func (p *PodList) Delete(podName string) {
	p.m.Lock()
	defer p.m.Unlock()

	delete(p.pods, podName)
}

// Exists returns whether the given pod exists or not
func (p *PodList) Exists(podName string) bool {
	p.m.Lock()
	defer p.m.Unlock()

	_, ok := p.pods[podName]
	return ok
}

func (p *PodList) Length() int {
	return len(p.pods)
}
