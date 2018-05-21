package main

import (
	"flag"
	"os"
	"path/filepath"

	"github.com/onrik/logrus/filename"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/kubernetes/pkg/kubectl/metricsutil"
)

var GitCommit string

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}

func main() {

	var (
		kubeconfig *string
	)

	if home := homeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}

	namespace := flag.String("namespace", "default", "Specify the namespace")

	flag.Parse()

	Globals.Namespace = *namespace

	filenameHook := filename.NewHook()
	filenameHook.Field = "source"
	log.AddHook(filenameHook)

	file, err := os.OpenFile("/tmp/kube-top.log", os.O_CREATE|os.O_WRONLY, 0666)
	if err == nil {
		log.Out = file
	} else {
		log.Info("Failed to log to file, using default stderr")
	}

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	heapsterClient := metricsutil.DefaultHeapsterMetricsClient(clientset.Core())
	Globals.KubeClient = NewKubeClient(clientset, heapsterClient)

	KubeTop()
}
