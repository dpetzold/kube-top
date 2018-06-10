package main

import (
	"flag"
	"os"
	"path/filepath"

	"github.com/Sirupsen/logrus"

	"github.com/dpetzold/kube-top/pkg/global"
	"github.com/dpetzold/kube-top/pkg/kube"
	"github.com/dpetzold/kube-top/pkg/ui"

	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/kubernetes/pkg/kubectl/metricsutil"
)

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}

var log = logrus.New()

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

	global.Namespace = *namespace

	/*
		filenameHook := filename.NewHook()
		filenameHook.Field = "source"
		log.AddHook(filenameHook)
	*/

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
	global.KubeClient = kube.NewKubeClient(clientset, heapsterClient)

	ui.KubeTop()
}
