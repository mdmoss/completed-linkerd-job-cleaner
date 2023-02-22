package main

import (
	"context"
	"flag"
	"log"
	"net/http"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

const LinkerdContainerName = "linkerd-proxy"

var verbose = false

func main() {
	log.Println("completed-linkerd-job-cleaner is starting...")

	kubeconfig := flag.String("kubeconfig", "", "Path to kubeconfig file, for running out-of-cluster")
	verboseFlag := flag.Bool("verbose", false, "Provide verbose output")
	shutdownSelf := flag.Bool("shutdown-self", false, "Post to http://localhost:4191/shutdown to shutdown our own proxy when finished")
	flag.Parse()

	if *verboseFlag {
		verbose = *verboseFlag
	}

	var config *rest.Config
	var err error

	if *kubeconfig != "" {
		log.Printf("Using kubeconfig at %s\n", *kubeconfig)
		config, err = clientcmd.BuildConfigFromFlags("", *kubeconfig)
	} else {
		config, err = rest.InClusterConfig()
	}

	if err != nil {
		panic(err.Error())
	}

	clientset, err := kubernetes.NewForConfig(config)

	if err != nil {
		panic(err.Error())
	}

	pods, err := clientset.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	deleted := 0

	for _, pod := range pods.Items {
		if onlyLinkerdProxyRemaining(pod) {
			owner := getSingleOwningJob(pod)

			if owner != nil {
				deleteJobByNameAndNamespace(clientset, owner.Name, pod.Namespace)
				deleted += 1
			}
		}
	}

	log.Printf("cleaned up %d pod(s) with lingering containers (out of %d total)", deleted, len(pods.Items))

	if *shutdownSelf {
		log.Printf("shutting down local proxy")
		http.Post("http://localhost:4191/shutdown", "", nil)
	}
}

func onlyLinkerdProxyRemaining(pod v1.Pod) bool {
	linkerdProxyRunning := false

	if verbose {
		log.Printf("checking Pod %s/%s\n", pod.Namespace, pod.Name)
	}

	for _, container := range pod.Status.ContainerStatuses {
		if container.State.Waiting != nil {
			if verbose {
				log.Printf("skipping: container %s is waiting to start\n", container.Name)
			}
			return false
		}

		if container.State.Running != nil {
			if container.Name == "linkerd-proxy" {
				linkerdProxyRunning = true
			} else {
				if verbose {
					log.Printf("skipping: container %s is still running\n", container.Name)
				}
				return false
			}
		}
	}

	if !linkerdProxyRunning {
		if verbose {
			log.Printf("skipping: no linkerd-proxy container found\n")
		}
	}

	return linkerdProxyRunning
}

func getSingleOwningJob(pod v1.Pod) *metav1.OwnerReference {
	if len(pod.OwnerReferences) == 1 && pod.OwnerReferences[0].Kind == "Job" {
		return &pod.OwnerReferences[0]
	}
	return nil
}

func deleteJobByNameAndNamespace(clientset *kubernetes.Clientset, name, namespace string) {
	log.Printf("deleting batchv1/Job:%s/%s\n", namespace, name)

	foreground := metav1.DeletePropagationForeground
	clientset.BatchV1().Jobs(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{PropagationPolicy: &foreground})
}
