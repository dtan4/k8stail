package main

import (
	"context"
	"fmt"
	"log"
	"math"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"strings"
	"time"

	flag "github.com/spf13/pflag"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	debugAddress     = ":6060"
	logSecondsOffset = 10
)

var (
	sinceSeconds = int64(math.Ceil(float64(logSecondsOffset) / float64(time.Second)))
)

type containerNames []string

func (c *containerNames) String() string {
	return "[" + strings.Join(*c, ", ") + "]"
}

func (c *containerNames) Type() string {
	return "containerNames"
}

func (c *containerNames) Set(value string) error {
	*c = append(*c, value)
	return nil
}

func (c *containerNames) Contains(value string) bool {
	for _, n := range *c {
		if value == n {
			return true
		}
	}
	return false
}


func main() {
	var (
		debug       bool
		kubeContext string
		kubeconfig  string
		labels      string
		namespace   string
		timestamps  bool
		version     bool
		excludeContainer []string
	)

	flags := flag.NewFlagSet("k8stail", flag.ExitOnError)
	flags.Usage = func() {
		flags.PrintDefaults()
	}

	flags.StringVar(&kubeContext, "context", "", "Kubernetes context")
	flags.BoolVar(&debug, "debug", false, "Debug mode using pprof (http://localhost:6060)")
	flags.StringVar(&kubeconfig, "kubeconfig", "", "Path of kubeconfig")
	flags.StringVarP(&labels, "labels", "l", "", "Label filter query")
	flags.StringVarP(&namespace, "namespace", "n", "", "Kubernetes namespace")
	flags.BoolVarP(&timestamps, "timestamps", "t", false, "Include timestamps on each line")
	flags.BoolVarP(&version, "version", "v", false, "Print version")
	flags.StringSliceVarP(
		&excludeContainer, "exclude-container", "e", []string{},
		"Exclude containers name (can specify multiple or separate values with commas: name1,name2)")

	if err := flags.Parse(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if kubeconfig == "" {
		if os.Getenv("KUBECONFIG") != "" {
			kubeconfig = os.Getenv("KUBECONFIG")
		} else {
			kubeconfig = clientcmd.RecommendedHomeFile
		}
	}

	if version {
		printVersion()
		os.Exit(0)
	}

	if debug {
		go func() {
			log.Println(http.ListenAndServe(debugAddress, nil))
		}()
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

	rawConfig, err := clientConfig.RawConfig()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	var currentContext string

	if kubeContext == "" {
		currentContext = rawConfig.CurrentContext
	} else {
		currentContext = kubeContext
	}

	if namespace == "" {
		if rawConfig.Contexts[currentContext].Namespace == "" {
			namespace = metav1.NamespaceDefault
		} else {
			namespace = rawConfig.Contexts[currentContext].Namespace
		}
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)

	logger := NewLogger()
	logger.PrintHeader(currentContext, namespace, labels)

	watcher, err := clientset.CoreV1().Pods(namespace).Watch(metav1.ListOptions{
		LabelSelector: labels,
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	added, finished, deleted := Watch(ctx, watcher, excludeContainer)

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

	<-sigCh
	cancel()
}
