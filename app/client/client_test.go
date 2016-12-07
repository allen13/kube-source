package kubernetes

import (
	"testing"
	"k8s.io/client-go/pkg/api/v1"
	"io/ioutil"
)

func TestKubeSourceClient(t *testing.T) {
	client, err := NewClient("integration-containers")
	if err != nil{
		t.Error(err)
	}

	ports := []v1.ServicePort{
		{
			Name: "6379-tcp",
			Protocol: "tcp",
			Port: int32(6379),
		},
	}

	request := ContainerRequest{
		dockerImage: "redis:alpine",
		ports: ports,
	}

	token, err := ioutil.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/token")
	if err != nil {
		t.Error(err)
	}

	containerResource, err := client.CreateContainerResource(request, string(token))

	if err != nil {
		t.Error(err)
	}

	nodePort := containerResource.ports[0]
	if !(nodePort > 30000 && nodePort < 32767){
		t.Errorf("Failed to get back valid not port. Got %d", nodePort)
	}

	err = client.DeleteContainerResource(containerResource.name, string(token))
	if err != nil {
		t.Error(err)
	}
}
