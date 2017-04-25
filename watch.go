package main

import (
	"context"

	"k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/pkg/watch"
)

// Watch starts and listens Kubernetes Pod events
func Watch(ctx context.Context, watcher watch.Interface) (chan *v1.Pod, chan *v1.Pod) {
	added := make(chan *v1.Pod)
	deleted := make(chan *v1.Pod)

	go func() {
		for {
			select {
			case e := <-watcher.ResultChan():
				switch e.Type {
				case watch.Added:
					pod := e.Object.(*v1.Pod)

					if pod.Status.Phase != v1.PodRunning {
						continue
					}

					added <- pod
				case watch.Modified:
					pod := e.Object.(*v1.Pod)

					if pod.Status.Phase != v1.PodRunning {
						continue
					}

					added <- pod
				case watch.Deleted:
					pod := e.Object.(*v1.Pod)

					deleted <- pod
				}

			case <-ctx.Done():
				watcher.Stop()
				return
			}
		}
	}()

	return added, deleted
}
