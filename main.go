package main

import (
	"flag"
	"fmt"
	"os"
	"sync"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	defaultNamespace = "default"
)

func main() {
	var (
		kubeconfig string
		labels     string
		namespace  string
	)

	flags := flag.NewFlagSet("k8stail", flag.ExitOnError)
	flags.Usage = func() {
		flags.PrintDefaults()
	}

	flags.StringVar(&kubeconfig, "kubeconfig", clientcmd.RecommendedHomeFile, fmt.Sprintf("Path of kubeconfig (Default: %s)", clientcmd.RecommendedHomeFile))
	flags.StringVar(&namespace, "namespace", v1.NamespaceDefault, fmt.Sprintf("Kubernetes namespace (Default: %s)", v1.NamespaceDefault))

	if err := flags.Parse(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	for 0 < flags.NArg() {
		labels = flags.Args()[0]
		flags.Parse(flags.Args()[1:])
	}

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Printf("Namespace: %s\n", namespace)
	fmt.Printf("Labels:    %s\n", labels)
	fmt.Println("======")

	var m sync.Mutex
	var wg sync.WaitGroup

	for {
		pods, err := clientset.Core().Pods(namespace).List(v1.ListOptions{
			LabelSelector: labels,
		})
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		for _, p := range pods.Items {
			wg.Add(1)
			go func(pod v1.Pod) {
				defer wg.Done()

				m.Lock()
				defer m.Unlock()
				fmt.Println(pod.Name)
			}(p)
		}

		wg.Wait()
	}
}
