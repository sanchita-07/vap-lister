// package main

// import (
// 	"context"
// 	"flag"
// 	"fmt"
// 	"time"

// 	// "time"
// 	// kyvernov1 "kyverno.io/kyverno/api/kyverno/v1"
// 	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
// 	// "k8s.io/apimachinery/pkg/labels"
// 	// "k8s.io/client-go/applyconfigurations/admissionregistration/v1alpha1"
// 	"k8s.io/client-go/informers"
// 	"k8s.io/client-go/kubernetes"
// 	"k8s.io/client-go/tools/clientcmd"
// )

// //take namespace as an argument

// func main() {
// 	// Get the kubeconfig file path
// 	kubeconfigPath := flag.String("kubeconfig", "", "path to the kubeconfig file")
// 	// namespace := flag.String("namespace", "", "namespace to list validating admission policies")
// 	flag.Parse()

// 	// Load the kubeconfig file
// 	config, err := clientcmd.LoadFromFile(*kubeconfigPath)
// 	if err != nil {
// 		fmt.Printf("Failed to load kubeconfig file: %v\n", err)
// 		return
// 	}

// 	// Create the REST config from the loaded kubeconfig
// 	restConfig, err := clientcmd.NewDefaultClientConfig(*config, &clientcmd.ConfigOverrides{}).ClientConfig()
// 	if err != nil {
// 		fmt.Printf("Failed to create REST config: %v\n", err)
// 		return
// 	}

// 	// Create the Kubernetes clientset
// 	clientset, err := kubernetes.NewForConfig(restConfig)
// 	if err != nil {
// 		fmt.Printf("Failed to create Kubernetes clientset: %v\n", err)
// 		return
// 	}

// 	// kubeInformer := informers.NewSharedInformerFactory(setup.KubeClient, resyncPeriod)

// 	// vapLister := kubeInformer.Admissionregistration().V1alpha1().ValidatingAdmissionPolicies().Lister()
	
// 	// kubeInformer := informers.NewSharedInformerFactory(clientset, time.Second*30, informers.WithNamespace(*namespace))
	
// 	// //Get the ValidatingAdmissionPolicy Informer
// 	// vapLister := kubeInformer.Admissionregistration().V1alpha1().ValidatingAdmissionPolicies().Lister()

// 	//Create the informer 
// 	kubeInformer := informers.NewSharedInformerFactory(clientset, time.Second*5)
	
// 	//Start the informer
// 	stopCh := make(chan struct{})
// 	defer close(stopCh)


// 	kubeInformer.Start(stopCh)
// 	kubeInformer.WaitForCacheSync(stopCh)

// 	//Wait for the caches to sync
// 	// if kubeInformer.WaitForCacheSync(stopCh) {
// 	// 	fmt.Println("Failed to sync caches")
// 	// 	return
// 	// }

// 	// //Get the Pod Informer
// 	// podInformer := kubeInformer.Core().V1().Pods()

// 	//Get the ValidatingAdmissionPolicy Informer
// 	// vapLister := kubeInformer.Admissionregistration().V1alpha1().ValidatingAdmissionPolicies().Lister()

// 	// List all validating admission policies
// 	policies, err := clientset.AdmissionregistrationV1alpha1().ValidatingAdmissionPolicies().List(context.TODO(), metav1.ListOptions{})
// 	// policies, err := vapLister.List(labels.Everything())
// 	if err != nil {
// 		fmt.Printf("Failed to list validating admission policies: %v\n", err)
// 		return
// 	}

// 	// Print the policies
// 	fmt.Println("Validating Admission Policies:")
// 	for _, policy := range policies.Items {
// 		fmt.Printf("- %s\n", policy.Name)
// 	}

// 	// vaplist, err := vapLister.List(labels.Everything())
// 	// if err != nil {
// 	// 	fmt.Printf("Failed to list validating admission policies: %v\n", err)
// 	// 	return
// 	// }


// 	// // Get the list of Pods
// 	// podList, err := podInformer.Lister().List(labels.Everything())
// 	// if err != nil {
// 	// 	fmt.Printf("Failed to list Pods: %v\n", err)
// 	// 	return
// 	// }

// 	// Print the pods
// 	// fmt.Println("Pods:")
// 	// for _, pod := range podList {
// 	// 	fmt.Printf("- %s\n", pod.Name)
// 	// }

// 	// for _, pod := range podList {
// 	// 		// Apply the ValidatingAdmissionPolicy to the Pod
// 	// 		for _, policy := range policies.Items {
// 	// 			applyConfig := .ApplyPolicyConfig{
// 	// 				Policy                    kyvernov1.PolicyInterface
// 	// 				ValidatingAdmissionPolicy v1alpha1.ValidatingAdmissionPolicy
// 	// 				Resource                  *unstructured.Unstructured
// 	// 				Variables                 map[string]interface{}
// 	// 				PolicyReport              bool
// 	// 				NamespaceSelectorMap      map[string]map[string]string
// 	// 				Stdin                     bool
// 	// 				PrintPatchResource        bool
// 	// 				RuleToCloneSourceResource map[string]string
// 	// 				Client                    kyverno.Interface
// 	// 			}

// 	// 			}
// 	// 			_, err := clientset.AdmissionregistrationV1alpha1().ValidatingAdmissionPolicies().Apply(context.TODO(), applyConfig, metav1.ApplyOptions{})
// 	// 			if err != nil {
// 	// 				fmt.Printf("Failed to apply ValidatingAdmissionPolicy %s to Pod %s: %v\n", policy.Name, pod.Name, err)
// 	// 			} else {
// 	// 				fmt.Printf("Applied ValidatingAdmissionPolicy %s to Pod %s\n", policy.Name, pod.Name)
// 	// 		}
// 		// }
// 	// }
// }