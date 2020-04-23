package cmds

import (

	"fmt"
	"path/filepath"
	"os"
	"time"
	"k8s.io/client-go/tools/clientcmd"
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
	K8SAPIServer = `https://192.168.1.170:8443`

	// K8SAPIToken 表示 ServiceAccount shiyanlou-admin 的 Secret 对应的 Token
	K8SAPIToken = `eyJhbGciOiJSUzI1NiIsImtpZCI6IiJ9.eyJpc3MiOiJrdWJlcm5ldGVzL3NlcnZpY2VhY2NvdW50Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9uYW1lc3BhY2UiOiJrdWJlLXN5c3RlbSIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VjcmV0Lm5hbWUiOiJkZWZhdWx0LXRva2VuLWxobTY0Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9zZXJ2aWNlLWFjY291bnQubmFtZSI6ImRlZmF1bHQiLCJrdWJlcm5ldGVzLmlvL3NlcnZpY2VhY2NvdW50L3NlcnZpY2UtYWNjb3VudC51aWQiOiIzNjllZWE0Mi00MjlhLTQ2MDItOWU3MS00YTFiZmY1M2E3YjMiLCJzdWIiOiJzeXN0ZW06c2VydmljZWFjY291bnQ6a3ViZS1zeXN0ZW06ZGVmYXVsdCJ9.naJIEEq1X_azgMIWwSodNQcWMXZDBEDjgRR_NPxLdONpkDYml3AYhHx3XGt4ZUtGc1Z4poa8h5ZgLEbuhGpoQx4WFpy1SRqVwf5Rtg748kWl-BGzJTmjId1K1_4WQlqF1igUDl6v8c9yaiNs6x_8Z-pUSrHEC1VNRrPzBHgqYAS7qhNsHewsO7myIAw80ur2AfDZtnmffvEYU5Df9YxOXU6YdH-VT16e8I_DJsaS9nxRrAuYeq33h_2W4e5irVbVDYxoUfIkwYT4jm0bEg5f7yH3BtVQ9-hf3zjEOx3zh7z9RIjVe9OWoREmUgDmrJmrRlZFqzVWLzjBYoJsw_Bh8Q`

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
				restclient, err := CreateK8sClient()
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
	podRootCmd.AddCommand(&podCheckCmd)

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
func createK8sClient() (k8sClient *k8s.Clientset, err error) {
	cfg := restclient.Config{}
	cfg.Host = K8SAPIServer
	cfg.CAData = []byte(K8SCertificateData)
	cfg.BearerToken = K8SAPIToken
	cfg.Timeout = time.Second * time.Duration(K8SAPITimeout)
	k8sClient, err = k8s.NewForConfig(&cfg)
	return
}

func CreateK8sClient() (client *k8s.Clientset, err error) {
	//var b []byte
	h := homeDir()
	//fmt.Println(os.Getenv("HOME"))
	//fmt.Println(os.Getenv("USERPROFILE"))

	h = filepath.Join(h, ".kube", "config")
	//_, erro := os.Stat(h)

	//if erro != nil {
	//	b, _ = ioutil.ReadFile("F:\\go_project\\cobra4k8s\\config")
	//	//kubeconfig = flag.String("kubeconfig", "F:\\go_project\\cobra4k8s\\config", "absolute path to the kubeconfig file")
	//} else {
	//	b, _ = ioutil.ReadFile(h)
	//
	//	//kubeconfig = flag.String("kubeconfig", h, "(optional) absolute path to the kubeconfig file")
	//}

	//if  h != ""  {
	//	kubeconfig = flag.String("kubeconfig", h, "(optional) absolute path to the kubeconfig file")
	//} else {
	//	kubeconfig = flag.String("kubeconfig", "F:\\go_project\\cobra4k8s\\config", "absolute path to the kubeconfig file")
	//}
	//flag.Parse()

	//在 kubeconfig 中使用当前上下文环境，config 获取支持 url 和 path 方式

	//b, _ := ioutil.ReadFile("C:\\Users\\Administrator\\.kube\\config")
	config, err := clientcmd.BuildConfigFromFlags("", h)
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
