package main

import (
	"context"
	"flag"
	"fmt"
	"path/filepath"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func main() {
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "(optional) absolute path to the kubeconfig file")
	}

	flag.Parse()

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		utilruntime.HandleError(err)
		return
	}

	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		utilruntime.HandleError(err)
		return
	}

	podList, err := client.
		CoreV1().
		Pods(corev1.NamespaceDefault).
		List(
			context.TODO(),
			metav1.ListOptions{
				LabelSelector: "app=echoserver",
			})
	if err != nil {
		utilruntime.HandleError(err)
		return
	}

	for _, pod := range podList.Items {
		fmt.Println(pod.Annotations)
		if _, ok := pod.Annotations["hello"]; !ok {
			pod.Annotations["hello"] = "world"
			pod, err := client.
				CoreV1().
				Pods(corev1.NamespaceDefault).
				Update(context.TODO(), &pod, metav1.UpdateOptions{})
			if err != nil {
				utilruntime.HandleError(err)
				return
			}
			fmt.Println("updating: ....")
			fmt.Println(pod.Annotations)
		}
	}
}
