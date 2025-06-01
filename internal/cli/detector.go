package cli

import (
	"fmt"
	"os/exec"
)

func DetectCLI(kubeconfig string) (string, CLI, error) {
	// Check for oc first
	if path, err := exec.LookPath("oc"); err == nil {
		return "oc", &OCClient{BaseCLI: BaseCLI{command: path, kubeconfig: kubeconfig}}, nil
	}

	// Fall back to kubectl
	if path, err := exec.LookPath("kubectl"); err == nil {
		return "kubectl", &KubectlClient{BaseCLI: BaseCLI{command: path, kubeconfig: kubeconfig}}, nil
	}

	return "", nil, fmt.Errorf("neither oc nor kubectl found in PATH")
}

type OCClient struct {
	BaseCLI
}

func (c *OCClient) Supports(feature string) bool {
	// OpenShift specific features
	return feature == "projects" || feature == "deploymentconfig"
}

type KubectlClient struct {
	BaseCLI
}

func (c *KubectlClient) Supports(feature string) bool {
	// Kubernetes specific features
	return feature == "namespace" || feature == "deployment"
}
