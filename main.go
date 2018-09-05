/*
Copyright 2018 tkggo

*/

// Note: the example only works with the code within the same release/branch.
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"

	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	// Uncomment the following line to load the gcp plugin (only required to authenticate against GKE clusters).
	// _ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

func init() {
	// Log as JSON instead of the default ASCII formatter.
	//log.SetFormatter(&log.JSONFormatter{})

	log.SetFormatter(&log.JSONFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	log.SetLevel(log.InfoLevel)
}

func main() {
	var kubeconfig *string
	if home := homeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	for {
		//pods, err := clientset.CoreV1().NetworkPolicy("").List(metav1.ListOptions{})
		netpol, err := clientset.NetworkingV1().NetworkPolicies("").List(metav1.ListOptions{})
		//clientset.CoreV1().
		//list := &networkingv1.NetworkPolicyList{ListMeta: obj.(*networkingv1.NetworkPolicyList).ListMeta}
		if err != nil {
			panic(err.Error())
		}
		fmt.Printf("There are %d netpol in the cluster\n", len(netpol.Items))

		for _, np := range netpol.Items {
			//log.Info(np.Name)

			/*
				{"metadata":{"name":"web-allow-all-namespaces","namespace":"np","selfLink":"/apis/networking.k8s.io/v1/namespaces/np/networkpolicies/web-allow-all-namespaces","uid":"718810d5-b14c-11e8-91fd-023c04a44d83","resourceVersion":"3675","generation":1,"creationTimestamp":"2018-09-05T20:44:11Z"},
				"spec":{"podSelector":{"matchLabels":{"app":"web"}},
				"ingress":[{"from":[{"namespaceSelector":{}}]}],"policyTypes":["Ingress"]}}
			*/

			log.WithFields(log.Fields{
				"Ingress":     np.Spec.Ingress,
				"PodSelector": np.Spec.PodSelector,
			}).Info(np.Name, np.Namespace)

			/*
				data, err := json.Marshal(np.Spec.Ingress)
				if err != nil {
					log.Fatal(err)
				}
				fmt.Printf("%s\n", data)
			*/

			// Examples for error handling:
			// - Use helper functions like e.g. errors.IsNotFound()
			// - And/or cast to StatusError and use its properties like e.g. ErrStatus.Message

		}

		time.Sleep(10 * time.Second)
	}
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}
