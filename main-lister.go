package main

import (
	"fmt"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"path/filepath"

	// admissionv1 "k8s.io/api/admissionregistration/v1"
	informerv1 "k8s.io/client-go/informers"
	listerv1 "k8s.io/client-go/listers/admissionregistration/v1alpha1"
)

func getClientSet() (*kubernetes.Clientset, error) {
	// Use the current context in kubeconfig
	configPath := filepath.Join(homedir.HomeDir(), ".kube", "config")
	config, err := clientcmd.BuildConfigFromFlags("", configPath)
	if err != nil {
		return nil, err
	}

	// Create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return clientset, nil
}

func createInformerLister(clientset *kubernetes.Clientset) (cache.SharedIndexInformer, listerv1.ValidatingAdmissionPolicyLister) {
	factory := informerv1.NewSharedInformerFactory(clientset, 0)
	informer := factory.Admissionregistration().V1alpha1().ValidatingAdmissionPolicies().Informer()
	lister := factory.Admissionregistration().V1alpha1().ValidatingAdmissionPolicies().Lister()

	// Add event handlers for the informer
	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{})

	return informer, lister
}

func listValidatingAdmissionPolicies(lister listerv1.ValidatingAdmissionPolicyLister) error {
	validatingAdmissionPolicies, err := lister.List(labels.Everything())
	if err != nil {
		return err
	}

	for _, policy := range validatingAdmissionPolicies {
		fmt.Printf("Validating Admission Policy: %s\n", policy.Name)
	}

	return nil
}

func main() {
	clientset, err := getClientSet()
	if err != nil {
		panic(err.Error())
	}

	informer, lister := createInformerLister(clientset)
	stopCh := make(chan struct{})
	defer close(stopCh)

	// Start the Informer
	go informer.Run(stopCh)

	// Wait until the Informer is synced
	if !cache.WaitForCacheSync(stopCh, informer.HasSynced) {
		panic("Failed to sync")
	}

	// List the validating admission policies
	err = listValidatingAdmissionPolicies(lister)
	if err != nil {
		panic(err.Error())
	}

	// Block until a termination signal is received
	<-stopCh
}