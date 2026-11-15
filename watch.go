package main

import (
	"context"
	"fmt"
	"regexp"

	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/watch"
)

type Target struct {
	Namespace string
	Pod       string
	Container string
}

// NewTarget creates new Target object
func NewTarget(namespace, pod, container string) *Target {
	return &Target{
		Namespace: namespace,
		Pod:       pod,
		Container: container,
	}
}

// GetID returns target ID
func (t *Target) GetID() string {
	return t.Namespace + "_" + t.Pod + "_" + t.Container
}

// NewPodExcluder compiles a Pod name exclusion pattern. It returns nil if no
// pattern is given, in which case no Pod is excluded.
func NewPodExcluder(pattern string) (*regexp.Regexp, error) {
	if pattern == "" {
		return nil, nil
	}

	excludePods, err := regexp.Compile(pattern)
	if err != nil {
		return nil, fmt.Errorf("invalid --exclude-pods pattern %q: %w", pattern, err)
	}

	return excludePods, nil
}

// isPodExcluded tells whether logs of the given Pod should be skipped
func isPodExcluded(pod *v1.Pod, excludePods *regexp.Regexp) bool {
	return excludePods != nil && excludePods.MatchString(pod.Name)
}

// Watch starts and listens Kubernetes Pod events
func Watch(ctx context.Context, watcher watch.Interface, excludePods *regexp.Regexp) (chan *Target, chan *Target, chan *Target) {
	added := make(chan *Target)
	finished := make(chan *Target)
	deleted := make(chan *Target)

	go func() {
		for {
			select {
			case e := <-watcher.ResultChan():
				if e.Object == nil {
					return
				}

				pod := e.Object.(*v1.Pod)

				if isPodExcluded(pod, excludePods) {
					continue
				}

				switch e.Type {
				case watch.Added:
					if pod.Status.Phase != v1.PodRunning {
						continue
					}

					for _, container := range pod.Spec.Containers {
						added <- NewTarget(pod.Namespace, pod.Name, container.Name)
					}
				case watch.Modified:
					switch pod.Status.Phase {
					case v1.PodRunning:
						for _, container := range pod.Spec.Containers {
							added <- NewTarget(pod.Namespace, pod.Name, container.Name)
						}
					case v1.PodSucceeded, v1.PodFailed:
						for _, container := range pod.Spec.Containers {
							finished <- NewTarget(pod.Namespace, pod.Name, container.Name)
						}
					}
				case watch.Deleted:
					for _, container := range pod.Spec.Containers {
						deleted <- NewTarget(pod.Namespace, pod.Name, container.Name)
					}
				}

			case <-ctx.Done():
				watcher.Stop()
				return
			}
		}
	}()

	return added, finished, deleted
}
