package cmds

import (
	"fmt"
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
MIICyjCCAbKgAwIBAgIBADANBgkqhkiG9w0BAQsFADAVMRMwEQYDVQQDEwprdWJl
cm5ldGVzMCAXDTE5MTAzMTA3MDMxMVoYDzIxMTgxMDA3MDcwMzExWjAVMRMwEQYD
VQQDEwprdWJlcm5ldGVzMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA
tiRXUHv/L1jM8JsvLzr0jywnwzr+iPqhDKT4YBytzDyEEEbo5bkUK2mErP35W96u
a1juMoiY6sz5VJllg4KHCR9IYd8MuaarHrdTfW5jDVpBXmhQwPhFY8NfSYl4rzLS
97beQoINVVrRZBb95iaCqWeVHgMJRSH28dtIrQDvs0aOq45k1WbTkvw0ZtyfdEO5
twtYxU0ur0w+tNdSAvNx+IOTsLxpe2+DdJ+PcvxSFF9q0bq7rpIeVRwyAas+0dan
+VU1eddT3xYI27mXe6CacrHRfzo6YPr8iqJkl1Q1SJwdTDpTetfttcbN/TzaLDMY
FpkLD62njuCBCdgP22YUZQIDAQABoyMwITAOBgNVHQ8BAf8EBAMCAqQwDwYDVR0T
AQH/BAUwAwEB/zANBgkqhkiG9w0BAQsFAAOCAQEAMcuT1f1+0sSb/SSzpRRs2byx
QkfUE4d9a+ER67uiv394nEnNWB0a4rTv0ZoZsuaN5D3lMIPriOcdezpIVizdpW0z
IWKQFe2/zkv9bYt5kVQbnMuE1lOOshUQcsSB5QYWK61dBkW3vHev7zCZ7iiCqU5m
aFhC/bWTKIvG0bFuVxlNIHeKd68J4Iy5LrhgIMKxwDoQuybevAV/m6S40BHQm1bs
Kw7osP2ErwDAL/l858TRbELMldYookR97XG1wR3qtVkAwb+J5VWeL9hcuId52/ke
AAJURwhvfixRr4TVkV7QWVEbSb6ffhYKRoa2dTnBq0KrbSSZtbQA85P5CNWVlQ==
-----END CERTIFICATE-----`

	// K8SAPIServer 表示 Kubernetes API Server 地址
	K8SAPIServer = `https://192.168.1.166:6443`

	// K8SAPIToken 表示 ServiceAccount shiyanlou-admin 的 Secret 对应的 Token
	K8SAPIToken = `eyJhbGciOiJSUzI1NiIsImtpZCI6Ikx5dDk3NTFEWExlUDU4TkVSZG0xY2lobDBUYU9OYWpqVWZOY3NLaW5VaU0ifQ.eyJpc3MiOiJrdWJlcm5ldGVzL3NlcnZpY2VhY2NvdW50Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9uYW1lc3BhY2UiOiJkZWZhdWx0Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9zZWNyZXQubmFtZSI6InRlc3QtYWRtaW4tdG9rZW4tazVsZjgiLCJrdWJlcm5ldGVzLmlvL3NlcnZpY2VhY2NvdW50L3NlcnZpY2UtYWNjb3VudC5uYW1lIjoidGVzdC1hZG1pbiIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VydmljZS1hY2NvdW50LnVpZCI6IjFjMTQ5NTdlLWUzMmQtNDc4Zi1iZmQ3LTEwZDczYWZkYjNjZiIsInN1YiI6InN5c3RlbTpzZXJ2aWNlYWNjb3VudDpkZWZhdWx0OnRlc3QtYWRtaW4ifQ.FcnmdqEZYlp3dwofT1yNc5Ago7Ru6X1ABqingNy5PecNozh0Q13I3SR_BEXNq1_QTU9pSvaAuFrtOvYKer4pNrw6H1szTE5EGPEL_V6AvJkbQKcRwdGFcJ2qbnZ56NRzrrFBnobHS_zuYOZgHPpVqj3ASMBzT1K7VMlR9Q4PQnWkzqlAOIbG4VmE2ohzztBxft0GVHbd7SmtN8PzTfW65Nnm77GSDPcKVvaWD0dQdyzNopEaZ6MO0TwTqVw8sszxW8qn6gz3oiKvi1lPMFGz-7DIIuI3oMzZHn4W8krDDOZqpluMFk0ouJQZ68RT8_XR8DKhKC5l65bX0xkvDuSXAA`

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
