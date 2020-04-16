package cmds

import (
	"flag"
	"fmt"
	"k8s.io/client-go/tools/clientcmd"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
	k8s "k8s.io/client-go/kubernetes"
	restclient "k8s.io/client-go/rest"
)

// 以下这些值在每次实验中因为环境重建问题可能发现变化，需要自行填写
// 获取这些值的方法参考第三节我们讲解创建 ServiceAccount 名称为
// shiyanlou-admin的地方。
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
	K8SAPIServer = `https://192.168.1.170:6443`

	// K8SAPIToken 表示 ServiceAccount shiyanlou-admin 的 Secret 对应的 Token
	K8SAPIToken = `ZXlKaGJHY2lPaUpTVXpJMU5pSXNJbXRwWkNJNklpSjkuZXlKcGMzTWlPaUpyZFdKbGNtNWxkR1Z6TDNObGNuWnBZMlZoWTJOdmRXNTBJaXdpYTNWaVpYSnVaWFJsY3k1cGJ5OXpaWEoyYVdObFlXTmpiM1Z1ZEM5dVlXMWxjM0JoWTJVaU9pSnJkV0psTFhONWMzUmxiU0lzSW10MVltVnlibVYwWlhNdWFXOHZjMlZ5ZG1salpXRmpZMjkxYm5RdmMyVmpjbVYwTG01aGJXVWlPaUprWldaaGRXeDBMWFJ2YTJWdUxXeG9iVFkwSWl3aWEzVmlaWEp1WlhSbGN5NXBieTl6WlhKMmFXTmxZV05qYjNWdWRDOXpaWEoyYVdObExXRmpZMjkxYm5RdWJtRnRaU0k2SW1SbFptRjFiSFFpTENKcmRXSmxjbTVsZEdWekxtbHZMM05sY25acFkyVmhZMk52ZFc1MEwzTmxjblpwWTJVdFlXTmpiM1Z1ZEM1MWFXUWlPaUl6TmpsbFpXRTBNaTAwTWpsaExUUTJNREl0T1dVM01TMDBZVEZpWm1ZMU0yRTNZak1pTENKemRXSWlPaUp6ZVhOMFpXMDZjMlZ5ZG1salpXRmpZMjkxYm5RNmEzVmlaUzF6ZVhOMFpXMDZaR1ZtWVhWc2RDSjkubmFKSUVFcTFYX2F6Z01JV3dTb2ROUWNXTVhaREJFRGpnUlJfTlB4TGRPTnBrRFltbDNBWWhIeDNYR3Q0WlV0R2MxWjRwb2E4aDVaZ0xFYnVoR3BvUXg0V0ZweTFTUnFWd2Y1UnRnNzQ4a1dsLUJHekpUbWpJZDFLMV80V1FscUYxaWdVRGw2djhjOXlhaU5zNnhfOFotcFVTckhFQzFWTlJyUHpCSGdxWUFTN3FoTnNIZXdzTzdteUlBdzgwdXIyQWZEWnRubWZmdkVZVTVEZjlZeE9YVTZZZEgtVlQxNmU4SV9ESnNhUzlueFJyQXVZZXEzM2hfMlc0ZTVpclZiVkRZeG9VZklrd1lUNGptMGJFZzVmN3lIM0J0VlE5LWhmM3pqRU94M3poN3o5UklqVmU5T1dvUkVtVWdEbXJKbXJSbFpGcXpWV0x6akJZb0pzd19CaDhR`

	// K8SAPITimeout 表示超时时间
	K8SAPITimeout = 30
)

var namespace string
var version bool

// GetRootCommand 返回组装好的根命令
func GetRootCommand() *cobra.Command {
	// 定义根命令
	rootCmd := cobra.Command{
		Use: "k8sshell",
		Run: func(cmd *cobra.Command, args []string) {
			if version {
				restclient, err := createK8SClient()
				if err != nil { 
					fmt.Println("Err:", err)
					return
				}
				// 通过 ServerVersion 方法来获取版本号
				versionInfo, err := restclient.ServerVersion()
				
				if err != nil {
					fmt.Println("Err:", err)
					return
				}
				fmt.Println("Kubernetes Version:", versionInfo.String())
			}
		},
	}
	// 添加全局选项参数
	rootCmd.PersistentFlags().StringVarP(&namespace, "namespace", "n", "", "namespace")

	// 添加显示版本的信息
	rootCmd.Flags().BoolVarP(&version, "version", "v", false, "kubernetes version")

	// 添加子命令
	addCommands(&rootCmd)
	return &rootCmd
}

// addCommands 将各个命令拼装在一起 
func addCommands(rootCmd *cobra.Command) { 
	// Pod 
	podRootCmd.AddCommand(&podCreateCmd) 
	podRootCmd.AddCommand(&podUpdateCmd)  
	podRootCmd.AddCommand(&podGetCmd)  
	podRootCmd.AddCommand(&podDeleteCmd)  

	// Service
	serviceRootCmd.AddCommand(&serviceCreateCmd)
	serviceRootCmd.AddCommand(&serviceUpdateCmd)
	serviceRootCmd.AddCommand(&serviceGetCmd)
	serviceRootCmd.AddCommand(&serviceDeleteCmd)

	// Ingress
	ingressRootCmd.AddCommand(&ingressCreateCmd)
	ingressRootCmd.AddCommand(&ingressUpdateCmd)
	ingressRootCmd.AddCommand(&ingressGetCmd)
	ingressRootCmd.AddCommand(&ingressDeleteCmd)
	ingressRootCmd.AddCommand(&ingressCreateHTTPSCmd)

	// Secret
	secretRootCmd.AddCommand(&secretCreateCmd)
	secretRootCmd.AddCommand(&secretUpdateCmd)
	secretRootCmd.AddCommand(&secretGetCmd)
	secretRootCmd.AddCommand(&secretDeleteCmd)

	// Deployment
	deploymentRootCmd.AddCommand(&deploymentCreateCmd)
	deploymentRootCmd.AddCommand(&deploymentUpdateCmd)
	deploymentRootCmd.AddCommand(&deploymentGetCmd)
	deploymentRootCmd.AddCommand(&deploymentDeleteCmd)
	deploymentRootCmd.AddCommand(&deploymentUpgradeCmd)
    deploymentRootCmd.AddCommand(&deploymentRollbackCmd)

	// 组装命令
	rootCmd.AddCommand(&podRootCmd)
	rootCmd.AddCommand(&serviceRootCmd)
	rootCmd.AddCommand(&ingressRootCmd)
	rootCmd.AddCommand(&secretRootCmd)
	rootCmd.AddCommand(&deploymentRootCmd)
}

// createK8SClient 根据鉴权信息创建 Kubernetes 的连接客户端
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
	if home := homeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "config", "absolute path to the kubeconfig file")
	}
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
