package main

import (
	"log"
	"time"

	"github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type Node struct {
	Name          string
	DefaultLabels map[string]string // Default labels assigned by Cloud Provider
	MatchRules    []Rules
}

type Rules struct {
	DefaultLabel string
	DefaultValue string
	CustomLabel  string
	CustomValue  string
}

/*
Define custom labels here

Examples:

Rule1: For each node with instance-type=t2.xlarge make sure it is labeled with app=web
Rule2: For each node with zone=us-east-1b make sure it is labeled with app=api
*/

var RulesList = []Rules{
	Rules{
		DefaultLabel: "beta.kubernetes.io/instance-type",
		DefaultValue: "t2.xlarge",
		CustomLabel:  "app",
		CustomValue:  "web",
	},
	Rules{
		DefaultLabel: "failure-domain.beta.kubernetes.io/zone",
		DefaultValue: "us-east-1b",
		CustomLabel:  "app",
		CustomValue:  "api",
	},
}

var clientset *kubernetes.Clientset

func getClient(pathToCfg string) (*kubernetes.Clientset, error) {
	var config *rest.Config
	var err error
	if pathToCfg == "" {
		logrus.Info("Using in cluster config")
		config, err = rest.InClusterConfig()
		// in cluster access
	} else {
		logrus.Info("Using out of cluster config")
		config, err = clientcmd.BuildConfigFromFlags("", pathToCfg)
	}
	if err != nil {
		return nil, err
	}
	return kubernetes.NewForConfig(config)
}

func iterateNodes(nodes *v1.NodeList) {

	if len(nodes.Items) > 0 {
		for _, node := range nodes.Items {

			defaultlabels := make(map[string]string)
			var matchrules []Rules

			n := Node{
				Name:          node.Name,
				DefaultLabels: defaultlabels,
				MatchRules:    matchrules,
			}

			for label, value := range node.Labels {
				if label != "" {
					defaultlabels[label] = value
				}

			}

			applyLabels(&n, &node)
			logrus.Info("--------------------------")
		}
	}
}

func applyLabels(node *Node, n *v1.Node) {

	for _, rule := range RulesList {
		patch := `{"metadata":{"labels":{"` + rule.CustomLabel + `":"` + rule.CustomValue + `"}}}`
		_, err := clientset.CoreV1().Nodes().Patch(n.Name, types.StrategicMergePatchType, []byte(patch))
		if err != nil {
			logrus.Info("Failed to update the node: %v", err)
			continue
		}
	}
	logrus.Infof("Node [%s] updated", node.Name)
}

func main() {

	var err error
	clientset, err = getClient("")
	if err != nil {
		panic(err.Error())
	}

	api := clientset.CoreV1()

	// reconcile loop
	for {
		nodes, err := api.Nodes().List(metav1.ListOptions{})
		if err != nil {
			log.Fatal(err)
		}
		iterateNodes(nodes)

		time.Sleep(60 * time.Second)
	}

}
