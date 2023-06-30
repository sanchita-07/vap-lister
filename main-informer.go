// package main

// import (
// 	"fmt"
// 	"path/filepath"

// 	"k8s.io/api/admissionregistration/v1alpha1"
// 	"k8s.io/client-go/informers"
// 	"k8s.io/client-go/kubernetes"
// 	"k8s.io/client-go/tools/cache"
// 	"k8s.io/client-go/tools/clientcmd"
// 	"k8s.io/client-go/util/homedir"
// )

// func getClientSet() (*kubernetes.Clientset, error) {
// 	// Use the current context in kubeconfig
// 	configPath := filepath.Join(homedir.HomeDir(), ".kube", "config")
// 	config, err := clientcmd.BuildConfigFromFlags("", configPath)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Create the clientset
// 	clientset, err := kubernetes.NewForConfig(config)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return clientset, nil
// }

// func createInformer(clientset *kubernetes.Clientset) cache.SharedIndexInformer {
// 	factory := informers.NewSharedInformerFactory(clientset, 0)
// 	informer := factory.Admissionregistration().V1alpha1().ValidatingAdmissionPolicies().Informer()

// 	// Add event handlers for the informer
// 	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
// 		AddFunc:    handleAdd,
// 		UpdateFunc: handleUpdate,
// 		DeleteFunc: handleDelete,
// 	})

// 	return informer
// }

// func handleAdd(obj interface{}) {
// 	config := obj.(*v1alpha1.ValidatingAdmissionPolicy)
// 	fmt.Printf("Added: %s\n", config.Name)
// }

// func handleUpdate(oldObj, newObj interface{}) {
// 	newConfig := newObj.(*v1alpha1.ValidatingAdmissionPolicy)
// 	fmt.Printf("Updated: %s\n", newConfig.Name)
// }

// func handleDelete(obj interface{}) {
// 	config := obj.(*v1alpha1.ValidatingAdmissionPolicy)
// 	fmt.Printf("Deleted: %s\n", config.Name)
// }

// func main() {
// 	clientset, err := getClientSet()
// 	if err != nil {
// 		panic(err.Error())
// 	}

// 	informer := createInformer(clientset)
// 	stopCh := make(chan struct{})
// 	defer close(stopCh)

// 	// Start the Informer
// 	go informer.Run(stopCh)

// 	// Wait until the Informer is synced
// 	if !cache.WaitForCacheSync(stopCh, informer.HasSynced) {
// 		panic("Failed to sync")
// 	}

// 	// Block until a termination signal is received
// 	<-stopCh
// }
