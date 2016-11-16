package main

import (
	"sync"
)

// ContainerList represents the goroutine-safe map
type ContainerList struct {
	containers map[string]map[string]bool
	m          sync.Mutex
}

// NewContainerList returns new ContainerList object
func NewContainerList() *ContainerList {
	return &ContainerList{
		containers: map[string]map[string]bool{},
		m:          sync.Mutex{},
	}
}

// Add adds the given container
func (p *ContainerList) Add(podName, containerName string) {
	p.m.Lock()
	defer p.m.Unlock()

	if _, ok := p.containers[podName]; !ok {
		p.containers[podName] = map[string]bool{}
	}

	p.containers[podName][containerName] = true
}

// Delete delete the given container
func (p *ContainerList) Delete(podName, containerName string) {
	p.m.Lock()
	defer p.m.Unlock()

	if _, ok := p.containers[podName]; !ok {
		return
	}

	delete(p.containers[podName], containerName)
}

// Exists returns whether the given container exists or not
func (p *ContainerList) Exists(podName, containerName string) bool {
	p.m.Lock()
	defer p.m.Unlock()

	if _, ok := p.containers[podName]; !ok {
		return false
	}

	_, ok := p.containers[podName][containerName]
	return ok
}

// Length returns the number of items
func (p *ContainerList) Length() int {
	length := 0

	p.m.Lock()
	defer p.m.Unlock()

	for _, cs := range p.containers {
		length += len(cs)
	}

	return length
}
