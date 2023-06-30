package main

import (
	"flag"
	"fmt"
	"strings"

	// "github.com/containerd/containerd/diff/apply"

	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
)

//take namespace as an argument

func main() {
	// Get the kubeconfig file path
	kubeconfigPath := flag.String("kubeconfig", "", "path to the kubeconfig file")
	// namespace := flag.String("namespace", "", "namespace to list validating admission policies")
	flag.Parse()

	// Load the kubeconfig file
	config, err := clientcmd.LoadFromFile(*kubeconfigPath)
	if err != nil {
		fmt.Printf("Failed to load kubeconfig file: %v\n", err)
		return
	}

	// Create the REST config from the loaded kubeconfig
	restConfig, err := clientcmd.NewDefaultClientConfig(*config, &clientcmd.ConfigOverrides{}).ClientConfig()
	if err != nil {
		fmt.Printf("Failed to create REST config: %v\n", err)
		return
	}

	// Create the Kubernetes clientset
	clientset, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		fmt.Printf("Failed to create Kubernetes clientset: %v\n", err)
		return
	}

	// kubeInformer := informers.NewSharedInformerFactory(setup.KubeClient, resyncPeriod)

	// vapLister := kubeInformer.Admissionregistration().V1alpha1().ValidatingAdmissionPolicies().Lister()

	// kubeInformer := informers.NewSharedInformerFactory(clientset, time.Second*30, informers.WithNamespace(*namespace))

	//Create the informer
	kubeInformer := informers.NewSharedInformerFactory(clientset, 0)

	vapLister := kubeInformer.Admissionregistration().V1alpha1().ValidatingAdmissionPolicies().Lister()
	vapInformer := kubeInformer.Admissionregistration().V1alpha1().ValidatingAdmissionPolicies().Informer()

	podLister := kubeInformer.Core().V1().Pods().Lister()
	podInformer := kubeInformer.Core().V1().Pods().Informer()

	//Start the informer
	stopCh := make(chan struct{})
	defer close(stopCh)

	go vapInformer.Run(stopCh)
	go podInformer.Run(stopCh)

	// kubeInformer.Start(stopCh)

	// Wait for the caches to sync
	if !cache.WaitForCacheSync(stopCh, podInformer.HasSynced, vapInformer.HasSynced) {
		fmt.Println("Failed to sync caches")
		return
	}

	//Get the ValidatingAdmissionPolicy Informer
	// vapLister := kubeInformer.Admissionregistration().V1alpha1().ValidatingAdmissionPolicies().Lister()

	// List all validating admission policies
	// policies, err := clientset.AdmissionregistrationV1alpha1().ValidatingAdmissionPolicies().List(context.TODO(), metav1.ListOptions{})
	// policies, err := vapLister.List(labels.Everything())
	// if err != nil {
	// 	fmt.Printf("Failed to list validating admission policies: %v\n", err)
	// 	return
	// }

	// // Print the policies
	// fmt.Println("Validating Admission Policies:")
	// for _, policy := range policies.Items {
	// 	fmt.Printf("- %s\n", policy.Name)
	// }

	vaplist, err := vapLister.List(labels.Everything())
	if err != nil {
		fmt.Printf("Failed to list validating admission policies: %v\n", err)
		return
	}

	fmt.Println("Validating Admission Policies:")
	for _, policy := range vaplist {
		fmt.Printf("- %s\n", policy.Name)
	}

	// Get the list of Pods
	podlist, err := podLister.List(labels.Everything())
	if err != nil {
		fmt.Printf("Failed to list Pods: %v\n", err)
		return
	}

	// Print the pods
	fmt.Println("Pods:")
	for _, pod := range podlist {
		fmt.Printf("- %s\n", pod.Name)
	}

	// var policyy v1alpha1.ValidatingAdmissionPolicy
	<-stopCh
	for _, policy := range vaplist {
		// Apply the ValidatingAdmissionPolicy to the Pod
		for _, pod := range podlist {
			// ...
			policyDecisions := applyPolicyToResource(policy, pod)
			for _, decision := range policyDecisions {
				if strings.Compare(string(decision.Action), "deny") == 0 {
					fmt.Println(decision.Message)
				} else {
					fmt.Println(decision.Action)
				}
			}

			// _, err := clientset.AdmissionregistrationV1alpha1().ValidatingAdmissionPolicies().Apply(context.TODO(), applyConfig, metav1.ApplyOptions{})
			if err != nil {
				fmt.Printf("Failed to apply ValidatingAdmissionPolicy %s to Pod %s: %v\n", policy.Name, pod.Name, err)
			} else {
				fmt.Printf("Applied ValidatingAdmissionPolicy %s to Pod %s\n", policy.Name, pod.Name)
			}
		}
	}
}
