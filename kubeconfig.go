package kubeconfig

import (
	"os"
	"path/filepath"

	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE")
}

func FindConfig() (string, bool) {
	if p, ok := os.LookupEnv("KUBECONFIG"); ok {
		return p, true
	}

	if home := homeDir(); home != "" {
		p := filepath.Join(home, ".kube", "config")
		_, err := os.Stat(p)
		if err != nil {
			return "", false
		}
		return p, true
	}
	return "", false
}

func Load(kubeconfig string) (*rest.Config, error) {
	if kubeconfig != "" {
		_, err := os.Stat(kubeconfig)
		if err == nil {
			// the first parameter of BuildConfigFromFlags is "masterUrl"
			return clientcmd.BuildConfigFromFlags("", kubeconfig)
		}
	}
	mykubeconfig, found := FindConfig()
	if found {
		return clientcmd.BuildConfigFromFlags("", mykubeconfig)
	}
	return rest.InClusterConfig()
}
