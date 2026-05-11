package utils

import (
	"context"
	"fmt"
	"net"
	"os"

	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// GetPodNodeInternalIP 获取指定 Pod 所在节点的内网 IP
func GetPodNodeInternalIP(clientSet *kubernetes.Clientset, namespace, podName string) (string, error) {
	pod, err := clientSet.CoreV1().Pods(namespace).Get(context.TODO(), podName, metav1.GetOptions{})
	if err != nil {
		return "", fmt.Errorf("failed to get pod: %v", err)
	}

	nodeName := pod.Spec.NodeName
	if nodeName == "" {
		return "", fmt.Errorf("pod is not scheduled to a node yet")
	}

	node, err := clientSet.CoreV1().Nodes().Get(context.TODO(), nodeName, metav1.GetOptions{})
	if err != nil {
		return "", fmt.Errorf("failed to get node: %v", err)
	}

	for _, address := range node.Status.Addresses {
		if address.Type == v1.NodeInternalIP {
			return address.Address, nil
		}
	}

	return "", fmt.Errorf("no internal IP found for node %s", nodeName)
}
func GetServiceIpK8s() (string, error) {
	config, err := clientcmd.BuildConfigFromFlags("", "~/.kube/config")
	if err != nil {
		panic(err.Error())
	}

	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	namespace := os.Getenv("NAMESPACE") // 替换为实际命名空间
	podName := os.Getenv("POD_NAME")    // 替换为实际 Pod 名称

	ip, err := GetPodNodeInternalIP(clientSet, namespace, podName)
	if err != nil {
		fmt.Println("Error:", err)
		return "", err
	}
	return ip, nil
}
func GetServiceIpLocal() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}
	for _, address := range addrs {
		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String(), nil
			}
		}
	}
	return "", fmt.Errorf("无法获取本地IP地址")
}
func GetServiceIp() (string, error) {
	//deploy := global.ServerConfig.Deploy
	return GetServiceIpLocal()
	//if deploy == "k8s" {
	//	return GetServiceIpK8s()
	//} else {
	//	return GetServiceIpLocal()
	//}

}
