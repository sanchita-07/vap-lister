package main

import (
	"context"
	"flag"
	"fmt"
	"strings"

	v1 "k8s.io/api/admissionregistration/v1"
	"k8s.io/api/admissionregistration/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apiserver/pkg/admission"
	"k8s.io/apiserver/pkg/admission/plugin/cel"
	"k8s.io/apiserver/pkg/admission/plugin/validatingadmissionpolicy"
	"k8s.io/apiserver/pkg/admission/plugin/webhook/matchconditions"
	celconfig "k8s.io/apiserver/pkg/apis/cel"
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

	deploymentLister := kubeInformer.Apps().V1().Deployments().Lister()
	deploymentInformer := kubeInformer.Apps().V1().Deployments().Informer()
	
	//Start the informer
	stopCh := make(chan struct{})
	defer close(stopCh)

	go vapInformer.Run(stopCh)
	go podInformer.Run(stopCh)
	go deploymentInformer.Run(stopCh)
	// kubeInformer.Start(stopCh)

	// Wait for the caches to sync
	if !cache.WaitForCacheSync(stopCh, podInformer.HasSynced, vapInformer.HasSynced, deploymentInformer.HasSynced) {
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

	// Get the list of Deployments
	deploymentlist, err := deploymentLister.List(labels.Everything())
	if err != nil {
		fmt.Printf("Failed to list Deployments: %v\n", err)
		return
	}

	// Print the deployments
	fmt.Println("Deployments:")
	for _, deployment := range deploymentlist {
		fmt.Printf("- %s\n", deployment.Name)
	}

	<-stopCh
	// for _, policy := range vaplist {
	// 	// Apply the ValidatingAdmissionPolicy to the Pod
	// 	for _, pod := range podlist {
	// 		// ...
	// 		podun, err := runtime.DefaultUnstructuredConverter.ToUnstructured(&pod)
	// 		policyDecisions := applyPolicyToResource(policy, &unstructured.Unstructured{Object: podun})
	// 		denied := false
	// 		for _, decision := range policyDecisions {
	// 			if strings.Compare(string(decision.Action), "deny") == 0 {
	// 				denied = true
	// 				fmt.Println(decision.Message)
	// 				break
	// 			} else {
	// 				fmt.Println(decision.Action)
	// 			}
	// 		}
	// 		if !denied {
	// 			fmt.Printf("Failed to apply ValidatingAdmissionPolicy %s to Pod %s: %v\n", policy.Name, pod.Name, err)
	// 		}
	// 		if err != nil {
	// 			fmt.Printf("Failed to apply ValidatingAdmissionPolicy %s to Pod %s: %v\n", policy.Name, pod.Name, err)
	// 		} else {
	// 			fmt.Printf("Applied ValidatingAdmissionPolicy %s to Pod %s\n", policy.Name, pod.Name)
	// 		}
	// 	}
	// }


	for _, policy := range vaplist {
		// Apply the ValidatingAdmissionPolicy to the Pod
		for _, deployment := range deploymentlist {
			// ...
			deploymentun, err := runtime.DefaultUnstructuredConverter.ToUnstructured(&deployment)
			policyDecisions := applyPolicyToResource(policy, &unstructured.Unstructured{Object: deploymentun})
			denied := false
			for _, decision := range policyDecisions {
				if strings.Compare(string(decision.Action), "deny") == 0 {
					denied = true
					fmt.Println(decision.Message)
					break
				} else {
					fmt.Println(decision.Action)
				}
			}
			if !denied {
				fmt.Printf("Failed to apply ValidatingAdmissionPolicy %s to Pod %s: %v\n", policy.Name, deployment.Name, err)
			}
			if err != nil {
				fmt.Printf("Failed to apply ValidatingAdmissionPolicy %s to Pod %s: %v\n", policy.Name, deployment.Name, err)
			} else {
				fmt.Printf("Applied ValidatingAdmissionPolicy %s to Pod %s\n", policy.Name, deployment.Name)
			}
		}
	}
}

func applyPolicyToResource(policy *v1alpha1.ValidatingAdmissionPolicy, resource *unstructured.Unstructured) []validatingadmissionpolicy.PolicyDecision {
	forbiddenReason := metav1.StatusReasonForbidden
	matchPolicyType := v1alpha1.Exact

	var validations []v1alpha1.Validation = policy.Spec.Validations
	var expressions, messageExpressions []cel.ExpressionAccessor

	for _, expression := range validations {
		message := fmt.Sprintf("error: failed to create %s: %s \"%s\" is forbidden: ValidatingAdmissionPolicy '%s' denied request: failed expression: %s", resource.GetKind(), resource.GetAPIVersion(), resource.GetName(), policy.Name, expression.Expression)
		condition := &validatingadmissionpolicy.ValidationCondition{
			Expression: expression.Expression,
			Message:    message,
			Reason:     &forbiddenReason,
		}

		messageCondition := &validatingadmissionpolicy.MessageExpressionCondition{
			MessageExpression: expression.MessageExpression,
		}

		expressions = append(expressions, condition)
		messageExpressions = append(messageExpressions, messageCondition)
	}

	filterCompiler := cel.NewFilterCompiler()
	filter := filterCompiler.Compile(expressions, cel.OptionalVariableDeclarations{HasParams: false, HasAuthorizer: false}, celconfig.PerCallLimit)

	compileErrors := filter.CompilationErrors()

	if len(compileErrors) > 0 {
		for _, err := range compileErrors {
			fmt.Println(err.Error())
		}
		return nil
	}

	messageExpressionCompiler := cel.NewFilterCompiler()
	messageExpressionfilter := messageExpressionCompiler.Compile(messageExpressions, cel.OptionalVariableDeclarations{HasParams: false, HasAuthorizer: false}, celconfig.PerCallLimit)

	admissionAttributes := admission.NewAttributesRecord(resource.DeepCopyObject(), nil, resource.GroupVersionKind(), resource.GetNamespace(), resource.GetName(), schema.GroupVersionResource{}, "", admission.Create, nil, false, nil)
	versionedAttr, _ := admission.NewVersionedAttributes(admissionAttributes, admissionAttributes.GetKind(), nil)

	ctx := context.TODO()
	failPolicy := v1.FailurePolicyType(*policy.Spec.FailurePolicy)

	matchConditions := policy.Spec.MatchConditions
	var matchExpressions []cel.ExpressionAccessor

	for _, expression := range matchConditions {
		condition := &matchconditions.MatchCondition{
			Name:       expression.Name,
			Expression: expression.Expression,
		}
		matchExpressions = append(matchExpressions, condition)
	}

	matchFilterCompiler := cel.NewFilterCompiler()
	matchFilter := matchFilterCompiler.Compile(matchExpressions, cel.OptionalVariableDeclarations{HasParams: false, HasAuthorizer: false}, celconfig.PerCallLimit)

	newMatcher := matchconditions.NewMatcher(matchFilter, nil, &failPolicy, string(matchPolicyType), "test")

	auditAnnotations := policy.Spec.AuditAnnotations
	var auditExpressions []cel.ExpressionAccessor

	for _, expression := range auditAnnotations {
		condition := &validatingadmissionpolicy.AuditAnnotationCondition{
			Key:             expression.Key,
			ValueExpression: expression.ValueExpression,
		}
		auditExpressions = append(auditExpressions, condition)
	}

	auditAnnotationFilterCompiler := cel.NewFilterCompiler()
	auditAnnotationFilter := auditAnnotationFilterCompiler.Compile(auditExpressions, cel.OptionalVariableDeclarations{HasParams: false, HasAuthorizer: false}, celconfig.PerCallLimit)

	validator := validatingadmissionpolicy.NewValidator(filter, newMatcher, auditAnnotationFilter, messageExpressionfilter, &failPolicy, nil)
	validateResult := validator.Validate(ctx, versionedAttr, nil, celconfig.RuntimeCELCostBudget)

	//fmt.Println(validateResult.AuditAnnotations[0].Action)

	return validateResult.Decisions
}
