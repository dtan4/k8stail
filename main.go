package main

import (
	"flag"
	"fmt"
	"os"
	"time"

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

	for {
		pods, err := clientset.Core().Pods(namespace).List(v1.ListOptions{})
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		fmt.Printf("There are %d pods in the cluster\n", len(pods.Items))
		time.Sleep(10 * time.Second)
	}
}
