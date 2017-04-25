package main

import (
	"context"
	"fmt"
	"math"
	"os"
	"time"

	"github.com/fatih/color"
	flag "github.com/spf13/pflag"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	defaultNamespace = "default"
	logSecondsOffset = 10
)

var (
	sinceSeconds = int64(math.Ceil(float64(logSecondsOffset) / float64(time.Second)))
)

func main() {
	var (
		kubeContext string
		kubeconfig  string
		labels      string
		namespace   string
		timestamps  bool
		version     bool
	)

	flags := flag.NewFlagSet("k8stail", flag.ExitOnError)
	flags.Usage = func() {
		flags.PrintDefaults()
	}

	flags.StringVar(&kubeContext, "context", "", "Kubernetes context")
	flags.StringVar(&kubeconfig, "kubeconfig", "", "Path of kubeconfig")
	flags.StringVar(&labels, "labels", "", "Label filter query")
	flags.StringVar(&namespace, "namespace", "", "Kubernetes namespace")
	flags.BoolVar(&timestamps, "timestamps", false, "Include timestamps on each line")
	flags.BoolVarP(&version, "version", "v", false, "Print version")

	if kubeconfig == "" {
		if os.Getenv("KUBECONFIG") != "" {
			kubeconfig = os.Getenv("KUBECONFIG")
		} else {
			kubeconfig = clientcmd.RecommendedHomeFile
		}
	}

	if err := flags.Parse(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if version {
		printVersion()
		os.Exit(0)
	}

	clientConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		&clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeconfig},
		&clientcmd.ConfigOverrides{CurrentContext: kubeContext})

	config, err := clientConfig.ClientConfig()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	bold := color.New(color.Bold).SprintFunc()

	rawConfig, err := clientConfig.RawConfig()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if namespace == "" {
		if rawConfig.Contexts[rawConfig.CurrentContext].Namespace == "" {
			namespace = v1.NamespaceDefault
		} else {
			namespace = rawConfig.Contexts[rawConfig.CurrentContext].Namespace
		}
	}

	fmt.Printf("%s %s\n", bold("Context:  "), rawConfig.CurrentContext)
	fmt.Printf("%s %s\n", bold("Namespace:"), namespace)
	fmt.Printf("%s %s\n", bold("Labels:   "), labels)
	color.New(color.FgYellow).Println("Press Ctrl-C to exit.")
	color.New(color.Bold).Println("----------")

	logger := NewLogger()

	watcher, err := clientset.Core().Pods(namespace).Watch(v1.ListOptions{
		LabelSelector: labels,
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	added, finished, deleted := Watch(ctx, watcher)

	tails := map[string]*Tail{}

	go func() {
		for target := range added {
			tail := NewTail(target.Namespace, target.Pod, target.Container, logger, sinceSeconds, timestamps)
			tails[target.GetID()] = tail
			tail.Start(ctx, clientset)
		}
	}()

	go func() {
		for target := range finished {
			id := target.GetID()

			if tails[id] == nil {
				continue
			}

			if tails[id].Finished {
				continue
			}

			tails[id].Finish()
		}
	}()

	go func() {
		for target := range deleted {
			id := target.GetID()

			if tails[id] == nil {
				continue
			}

			tails[id].Delete()
			delete(tails, id)
		}
	}()

	<-ctx.Done()
}
