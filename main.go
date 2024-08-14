package main

import (
	"context"
	"flag"
	"fmt"
	"path/filepath"
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

// NoNodesInKubernetes, Kubernetes cluster'ında hiç node olmadığında döndürülür.
type NoNodesInKubernetes struct{}

func (err NoNodesInKubernetes) Error() string {
	return "Kubernetes cluster'ında hiç node yok"
}

// PersistentVolumeClaimNotInStatus, bir PersistentVolumeClaim beklenen durumda olmadığında döndürülür.
type PersistentVolumeClaimNotInStatus struct {
	pvc            *corev1.PersistentVolumeClaim
	pvcStatusPhase *corev1.PersistentVolumeClaimPhase
}

func (err PersistentVolumeClaimNotInStatus) Error() string {
	return fmt.Sprintf("PersistentVolumeClaim %s beklenen %v durumunda değil", err.pvc.Name, *err.pvcStatusPhase)
}

func main() {
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(isteğe bağlı) kubeconfig dosyasının mutlak yolu")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "kubeconfig dosyasının mutlak yolu")
	}
	flag.Parse()

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	for {
		fmt.Println("Cluster Durumu:")
		checkPods(clientset)
		checkNamespaces(clientset)
		checkNodes(clientset)
		checkEvents(clientset)
		checkPersistentVolumeClaims(clientset)

		fmt.Println("\nDetaylı pod kontrolü:")
		namespace := "default"
		pod := "alpine-deployment-548dbddc9b-dnq9r"
		checkSpecificPod(clientset, namespace, pod)

		fmt.Println("\n-----------------------------------")
		time.Sleep(10 * time.Second)
	}
}

func checkPods(clientset *kubernetes.Clientset) {
	pods, err := clientset.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		fmt.Printf("Pod'ları listelerken hata oluştu: %v\n", err)
		return
	}
	fmt.Printf("Cluster'da %d pod var\n", len(pods.Items))
}

func checkNamespaces(clientset *kubernetes.Clientset) {
	namespaces, err := clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		fmt.Printf("Namespace'leri listelerken hata oluştu: %v\n", err)
		return
	}
	fmt.Printf("Cluster'da %d namespace var\n", len(namespaces.Items))
}

func checkNodes(clientset *kubernetes.Clientset) {
	nodes, err := clientset.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		fmt.Printf("Node'ları listelerken hata oluştu: %v\n", err)
		return
	}
	if len(nodes.Items) == 0 {
		fmt.Println(NoNodesInKubernetes{}.Error())
	} else {
		fmt.Printf("Cluster'da %d node var\n", len(nodes.Items))
	}
}

func checkEvents(clientset *kubernetes.Clientset) {
	events, err := clientset.CoreV1().Events("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		fmt.Printf("Event'leri listelerken hata oluştu: %v\n", err)
		return
	}
	fmt.Printf("Son 1 saatte %d event var\n", len(events.Items))
}

func checkPersistentVolumeClaims(clientset *kubernetes.Clientset) {
	pvcs, err := clientset.CoreV1().PersistentVolumeClaims("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		fmt.Printf("PersistentVolumeClaim'leri listelerken hata oluştu: %v\n", err)
		return
	}
	fmt.Printf("Cluster'da %d PersistentVolumeClaim var\n", len(pvcs.Items))
	
	for _, pvc := range pvcs.Items {
		expectedPhase := corev1.ClaimBound
		if pvc.Status.Phase != expectedPhase {
			err := PersistentVolumeClaimNotInStatus{&pvc, &expectedPhase}
			fmt.Println(err.Error())
		}
	}
}

func checkSpecificPod(clientset *kubernetes.Clientset, namespace, podName string) {
	pod, err := clientset.CoreV1().Pods(namespace).Get(context.TODO(), podName, metav1.GetOptions{})
	if errors.IsNotFound(err) {
		fmt.Printf("Pod %s namespace %s içinde bulunamadı\n", podName, namespace)
	} else if statusError, isStatus := err.(*errors.StatusError); isStatus {
		fmt.Printf("Pod %s namespace %s içinde alınan hata: %v\n", podName, namespace, statusError.ErrStatus.Message)
	} else if err != nil {
		fmt.Printf("Pod bilgisi alınırken hata oluştu: %v\n", err)
	} else {
		fmt.Printf("Pod %s namespace %s içinde bulundu\n", podName, namespace)
		fmt.Printf("Pod durumu: %s\n", pod.Status.Phase)
		fmt.Printf("Pod IP: %s\n", pod.Status.PodIP)
		fmt.Printf("Node: %s\n", pod.Spec.NodeName)
	}
}