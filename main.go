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
		version    bool
	)

	flags := flag.NewFlagSet("k8stail", flag.ExitOnError)
	flags.Usage = func() {
		flags.PrintDefaults()
	}

	flags.StringVar(&kubeconfig, "kubeconfig", "", fmt.Sprintf("Path of kubeconfig (Default: %s)", clientcmd.RecommendedHomeFile))
	flags.StringVar(&labels, "labels", "", "Label filter query (Default: \"\")")
	flags.StringVar(&namespace, "namespace", v1.NamespaceDefault, fmt.Sprintf("Kubernetes namespace (Default: %s)", v1.NamespaceDefault))
	flags.BoolVar(&timestamps, "timestamps", false, "Include timestamps on each line (default: false)")
	flags.BoolVar(&version, "version", false, "Print version")
	flags.BoolVar(&version, "v", false, "Print version")

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

	bold := color.New(color.Bold).SprintFunc()

	fmt.Printf("%s %s\n", bold("Namespace:"), namespace)
	fmt.Printf("%s %s\n", bold("Labels:   "), labels)
	color.New(color.Bold).Println("----------")

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

			for _, container := range pod.Spec.Containers {
				wg.Add(1)
				go func(p v1.Pod, c v1.Container) {
					defer wg.Done()

					rs, err := clientset.Core().Pods(namespace).GetLogs(p.Name, &v1.PodLogOptions{
						Container:    c.Name,
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
						logger.PrintPodLog(p.Name, c.Name, sc.Text(), timestamps)
					}

					logger.PrintColorizedLog(redBold, fmt.Sprintf("Pod:%s Container:%s has been deleted", p.Name, c.Name))
				}(pod, container)
			}
		}

		if runningPods.Length() == 0 {
			break
		}
	}

	wg.Wait()
}
