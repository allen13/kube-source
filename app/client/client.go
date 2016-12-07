package kubernetes

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/rest"
	"github.com/satori/go.uuid"
	"github.com/allen13/kube-source/app/config"
	"fmt"
)

type Client struct {
	clientset *kubernetes.Clientset
	namespace string
}

type ContainerRequest struct{
	dockerImage string
	ports []v1.ServicePort
}

type ContainerResponse struct {
	name string
	ip    string
	ports []int
}

func NewClient(namespace string)(client *Client, err error){
	config, err := rest.InClusterConfig()
	if err != nil {
		return
	}

	if err != nil {
		return
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return
	}

	client = &Client{clientset: clientset, namespace: namespace}

	return
}

func NewClientWithToken(namespace string, token string)(client *Client, err error){
	config, err := rest.InClusterConfig()
	if err != nil {
		return
	}
	config.BearerToken = token

	if err != nil {
		return
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return
	}

	client = &Client{clientset: clientset, namespace: namespace}

	return
}

func (c *Client) CreateService(name string, ports []v1.ServicePort)(*v1.Service, error){
	service := &v1.Service{
		ObjectMeta: v1.ObjectMeta{
			Name: name,
		},
		Spec: v1.ServiceSpec{
			Selector: map[string]string{"name": name},
			Type: "NodePort",
			Ports: ports,
		},
	}
	return c.clientset.Services(c.namespace).Create(service)
}

func (c *Client) DeleteService(name string)(error){
	return c.clientset.Services(c.namespace).Delete(name, &v1.DeleteOptions{})
}

func (c *Client) CreatePod(name string, dockerImage string)(*v1.Pod, error){
	pod := &v1.Pod{
		ObjectMeta: v1.ObjectMeta{
			Name: name,
			Labels: map[string]string{
				"name": name,
			},
		},
		Spec: v1.PodSpec{
			Containers: []v1.Container{
				{
					Name: name,
					Image: dockerImage,
				},
			},
		},
	}
	return c.clientset.Pods(c.namespace).Create(pod)
}

func (c *Client) ListPods()(*v1.PodList, error){
	return c.clientset.Pods(c.namespace).List(v1.ListOptions{})
}

func (c *Client) DeletePod(name string)(error){
	return c.clientset.Pods(c.namespace).Delete(name, &v1.DeleteOptions{})
}



func (c *Client) CreateContainerResource(containerRequest ContainerRequest, bearerToken string)(containerResponse ContainerResponse,err error){
	if err != nil {
		return
	}

	u1 := uuid.NewV4()
	name := "container-" + fmt.Sprintf("%s",u1)[0:7]

	svc, err := c.CreateService(name, containerRequest.ports)
	if err != nil {
		return
	}

	_, err = c.CreatePod(name, containerRequest.dockerImage)
	if err != nil {
		return
	}

	ports := []int{}

	for _ , svcPort := range svc.Spec.Ports {
		ports = append(ports, int(svcPort.NodePort))
	}

	containerResponse.name = name
	containerResponse.ports = ports
	containerResponse.ip = config.Get("container_ip")

	return
}

func (c *Client) DeleteContainerResource(name string, bearerToken string)(err error){
	if err != nil {
		return
	}

	err = c.DeletePod(name)
	if err != nil {
		return
	}

	err = c.DeleteService(name)

	return
}