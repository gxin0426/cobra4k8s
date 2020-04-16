package main

import (
	"flag"
	"fmt"
	"k8s.io/client-go/tools/clientcmd"
	"os"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8s "k8s.io/client-go/kubernetes"
	restclient "k8s.io/client-go/rest"
)


var (
	// K8SCertificateData 表示 Kubernetes 服务端证书
	K8SCertificateData = `-----BEGIN CERTIFICATE-----
MIIC5zCCAc+gAwIBAgIBATANBgkqhkiG9w0BAQsFADAVMRMwEQYDVQQDEwptaW5p
a3ViZUNBMB4XDTIwMDQxNDA2NTExNFoXDTMwMDQxMzA2NTExNFowFTETMBEGA1UE
AxMKbWluaWt1YmVDQTCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBAMCY
FPQXjPzE1ntbc0lqLFWKbs5wFd/VYuexyOFeSSTj8oZx9ti8+/5MsCb6tvUZ6g5u
eaIPD556frvFwPAYOB0gQ9PZjbRHHGaFE0vgNDKhXnyamZNYpiT+Bl3Do4K2Bncc
Ryop3u38GMsQqQC1fqzdkOw9sq0mHgoHot1nVWDJtVJwDXGlhkxVhkjIkCs/4p88
tX+7GMzk4y+SD0CMbqKsDwX3wLgddLXsh2PiIkBCx7Oj4hiuz8Fd3X7EX+koZJHg
46NzvKnaPjNTAMc8EBoAclQThv7KW3O4zaaoYRLJlNiH+BrmQGYaJ8PL3hCfT+9I
UuWiN0SpeWFxah1BM/8CAwEAAaNCMEAwDgYDVR0PAQH/BAQDAgKkMB0GA1UdJQQW
MBQGCCsGAQUFBwMCBggrBgEFBQcDATAPBgNVHRMBAf8EBTADAQH/MA0GCSqGSIb3
DQEBCwUAA4IBAQBJXzDxATTRQZ4UUA4vn/9uA2vrA5q9EEgN9FAvF3/0HV1u2qTw
bb0c81ASrC7UBQf1pisfK2HSTGDuC54pc6ODxd1jjQxkCrj7MW7BZo6eKfAFu5Wl
KnoE/dsQNqc/vJvcdq5vXIhY+/sMnI3RrCRee7YNqvA4k0U/xv+kl46bQX6pleoA
6siAulRD5bivT+bNQ8R3J4UqHSqBFuYHQGXsQuHTeTT0I2US3sZC+/TDiAGKBd9d
PlgIuvgubWo9urcFN0W+5lrfWexcK2WSVPWQlJuTTG+OLADwbXLsgerXHnG3hxCI
sg5XPmLsDpH0DR6IRNBZalZfYLtOP6b1C1W9
-----END CERTIFICATE-----`

	// K8SAPIServer 表示 Kubernetes API Server 地址
	K8SAPIServer = `https://192.168.1.170:8443`

	// K8SAPIToken 表示 ServiceAccount shiyanlou-admin 的 Secret 对应的 Token
	K8SAPIToken = `eyJhbGciOiJSUzI1NiIsImtpZCI6IiJ9.eyJpc3MiOiJrdWJlcm5ldGVzL3NlcnZpY2VhY2NvdW50Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9uYW1lc3BhY2UiOiJrdWJlLXN5c3RlbSIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VjcmV0Lm5hbWUiOiJkZWZhdWx0LXRva2VuLWxobTY0Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9zZXJ2aWNlLWFjY291bnQubmFtZSI6ImRlZmF1bHQiLCJrdWJlcm5ldGVzLmlvL3NlcnZpY2VhY2NvdW50L3NlcnZpY2UtYWNjb3VudC51aWQiOiIzNjllZWE0Mi00MjlhLTQ2MDItOWU3MS00YTFiZmY1M2E3YjMiLCJzdWIiOiJzeXN0ZW06c2VydmljZWFjY291bnQ6a3ViZS1zeXN0ZW06ZGVmYXVsdCJ9.naJIEEq1X_azgMIWwSodNQcWMXZDBEDjgRR_NPxLdONpkDYml3AYhHx3XGt4ZUtGc1Z4poa8h5ZgLEbuhGpoQx4WFpy1SRqVwf5Rtg748kWl-BGzJTmjId1K1_4WQlqF1igUDl6v8c9yaiNs6x_8Z-pUSrHEC1VNRrPzBHgqYAS7qhNsHewsO7myIAw80ur2AfDZtnmffvEYU5Df9YxOXU6YdH-VT16e8I_DJsaS9nxRrAuYeq33h_2W4e5irVbVDYxoUfIkwYT4jm0bEg5f7yH3BtVQ9-hf3zjEOx3zh7z9RIjVe9OWoREmUgDmrJmrRlZFqzVWLzjBYoJsw_Bh8Q`

	// K8SAPITimeout 表示超时时间
	K8SAPITimeout = 30
)
func main() {
	//// 获取根命令
	//rootCmd := cmds.GetRootCommand()
	//// 执行命令
	//err := rootCmd.Execute()
	//if err != nil {
	//	fmt.Println("Err:", err)
	//	return
	//}

	k8sClient, err := CreateK8sClient()
	if err != nil {
		fmt.Println("Err :", err)
		return
	}
	listOption := metav1.ListOptions{}
	l, err :=k8sClient.CoreV1().Pods("kube-system").List(listOption)
	if err != nil{
		fmt.Println(err)
	}
	for _, pod := range l.Items{
		//fmt.Println(pod.Name)
		if pod.Status.Phase != "Running"{
			fmt.Printf("%s的状态是：%s\n", pod.Name, pod.Status.Phase)
		}



	}
	//fmt.Println(l.Items)

}

func createK8SClient() (k8sClient *k8s.Clientset, err error) {
	cfg := restclient.Config{}
	cfg.Host = K8SAPIServer
	cfg.CAData = []byte(K8SCertificateData)
	cfg.BearerToken = K8SAPIToken
	cfg.Timeout = time.Second * time.Duration(K8SAPITimeout)
	k8sClient, err = k8s.NewForConfig(&cfg)
	return
}

func CreateK8sClient() (client *k8s.Clientset, err error){
	var kubeconfig *string
	//if home := homeDir(); home != "" {
	//	kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	//} else {
		kubeconfig = flag.String("kubeconfig", "F:\\go_project\\cobra4k8s\\config", "absolute path to the kubeconfig file")
	//}
	flag.Parse()

	//在 kubeconfig 中使用当前上下文环境，config 获取支持 url 和 path 方式
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	// 根据指定的 config 创建一个新的 clientset
	client, err = k8s.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	return
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}

