package main

import (
	"bufio"
	"flag"
	"fmt"
	"math"
	"os"
	"sync"
	"time"

	"github.com/fatih/color"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	defaultNamespace = "default"
	logSecondsOffset = 10
)

func main() {
	var (
		kubeconfig string
		labels     string
		namespace  string
		timestamps bool
	)

	flags := flag.NewFlagSet("k8stail", flag.ExitOnError)
	flags.Usage = func() {
		flags.PrintDefaults()
	}

	flags.StringVar(&kubeconfig, "kubeconfig", clientcmd.RecommendedHomeFile, fmt.Sprintf("Path of kubeconfig (Default: %s)", clientcmd.RecommendedHomeFile))
	flags.StringVar(&namespace, "namespace", v1.NamespaceDefault, fmt.Sprintf("Kubernetes namespace (Default: %s)", v1.NamespaceDefault))
	flags.BoolVar(&timestamps, "timestamps", false, "Include timestamps on each line (default: false)")

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

	c := color.New(color.Bold)

	c.Printf("Namespace: %s\n", namespace)
	c.Printf("Labels:    %s\n", labels)
	c.Println("======")

	var wg sync.WaitGroup

	runningPods := NewPodList()
	greenBold := color.New(color.FgGreen, color.Bold)
	redBold := color.New(color.FgRed, color.Bold)
	logger := NewLogger()

	for {
		pods, err := clientset.Core().Pods(namespace).List(v1.ListOptions{
			LabelSelector: labels,
		})
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		for _, pod := range pods.Items {
			if pod.Status.Phase != v1.PodRunning {
				continue
			}

			if runningPods.Exists(pod.Name) {
				continue
			}

			runningPods.Add(pod.Name)
			logger.PrintColorizedLog(greenBold, fmt.Sprintf("Pod %s has detected", pod.Name))
			sinceSeconds := int64(math.Ceil(float64(logSecondsOffset) / float64(time.Second)))

			wg.Add(1)
			go func(p v1.Pod) {
				defer wg.Done()

				rs, err := clientset.Core().Pods(namespace).GetLogs(p.Name, &v1.PodLogOptions{
					Follow:       true,
					SinceSeconds: &sinceSeconds,
					Timestamps:   timestamps,
				}).Stream()
				if err != nil {
					fmt.Fprintln(os.Stderr, err)
					os.Exit(1)
				}

				sc := bufio.NewScanner(rs)

				for sc.Scan() {
					logger.PrintPodLog(p.Name, sc.Text(), timestamps)
				}

				logger.PrintColorizedLog(redBold, fmt.Sprintf("Pod %s has been deleted", p.Name))
			}(pod)
		}

		if runningPods.Length() == 0 {
			break
		}
	}

	wg.Wait()
}
