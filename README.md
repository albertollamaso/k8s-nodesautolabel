# k8s-nodesautolabel
Microservice that keeps Kubernetes Nodes labels up to date.

### Why?

Let's say you have a kubernetes cluster on Bare-Metal or in a Cloud Provider that don't keep track of labels after a Maintenance windows. It might sounds familiar for you.

This microservice intends to make those labels to survive after a maintance windows, Node failures (Add/Delete) etc.

When you spin up Kubernes cluster in any  Cloud provider. By default they setup default labels in the metadata fields.

This microservice uses these default labels as reference to keep custom labels up to date. You can use any default label as a reference and define them in the code.

### How it works?

- You define a set of rules for each label you want to keep in any given node. As an example:

**Define custom labels here**

Example:

**Rule1:** For each node with instance-type=t2.xlarge make sure it is labeled with app=web
**Rule2:** For each node with zone=us-east-1b make sure it is labeled with app=api

```
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
```

- Kubernetes Service Account and Deployment are in folder `kubernetes/`
- The go modules are set with versions to work with Kubernetes v1.11.10. You can modify them accordinly in file **vendor.conf**

### How to use it?

Within your kubernetes cluster execute:
`kubectl apply -f kubernetes/ -n default`

### Improvements

There are few things to improve as implement the microservice using watchers. Any question,advise don't hesitate to contact me or PR.




