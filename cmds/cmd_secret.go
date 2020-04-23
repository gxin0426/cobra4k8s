package cmds

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"time"

	"github.com/spf13/cobra"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func init() {
	//create secret option parameter
	secretCreateCmd.Flags().StringVar(&secretCreateName, "name", "", "secret name")
	secretCreateCmd.Flags().StringVar(&secretCreateCertFile, "cert-file", "cert", "cert file")
	secretCreateCmd.Flags().StringVar(&secretCreatePrivateKeyFile, "key-file", "key", "secret name")
	//update secret option parameter

	secretUpdateCmd.Flags().StringVar(&secretUpdateName, "name", "", "secret name")
	secretUpdateCmd.Flags().StringVar(&secretUpdateCertFile, "cert-file", "", "cert file")
	secretUpdateCmd.Flags().StringVar(&secretUpdatePrivateKeyFile, "key-file", "", "private key file")
	//obtian secret option param
	secretGetCmd.Flags().StringVar(&secretGetName, "name", "", "secret name")
	//delete secret option parameter
	secretDeleteCmd.Flags().StringVar(&secretDeleteName, "name", "", "sercret name")
}

// secret 命令
var secretRootCmd = cobra.Command{
	Use:   "secret",
	Short: "secret is used to manage kubernetes secrets",
}

// secret Create 命令

var secretCreateName string
var secretCreateCertFile string
var secretCreatePrivateKeyFile string

var secretCreateCmd = cobra.Command{
	Use:   "create",
	Short: "create a new secret",
	Run: func(cmd *cobra.Command, args []string) {
		if secretCreateName == "" || secretCreatePrivateKeyFile == "" || namespace == "" {
			cmd.Help()
			return
		}
		// create k8s client
		k8sClient, err := CreateK8sClient()
		if err != nil {
			fmt.Println(err)
			return
		}

		//读取Https的证书和秘钥
		certData, err := readFileContent(secretCreateCertFile)
		if err != nil {
			fmt.Println("Err :", err)
			return
		}

		keyData, err := readFileContent(secretCreatePrivateKeyFile)
		if err != nil {
			fmt.Println(err)
			return
		}

		//create a secret

		newSecret := v1.Secret{
			Type: v1.SecretTypeTLS,
			StringData: map[string]string{
				v1.TLSCertKey:       string(certData),
				v1.TLSPrivateKeyKey: string(keyData),
			},
		}

		newSecret.Name = secretCreateName
		newSecret.Namespace = namespace

		//调用Create 的接口方法
		_, err = k8sClient.CoreV1().Secrets(namespace).Create(&newSecret)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println(" success create secret !!!")
	},
}

func readFileContent(filePath string) (fileContent []byte, err error) {
	fp, openErr := os.Open(filePath)
	if openErr != nil {
		fmt.Println(openErr)
		return
	}
	defer fp.Close()
	fileContent, err = ioutil.ReadAll(fp)
	if err != nil {
		return
	}
	return
}

// secret Update 命令
var secretUpdateName string
var secretUpdateCertFile string
var secretUpdatePrivateKeyFile string
var secretUpdateCmd = cobra.Command{
	Use:   "update",
	Short: "update a secret",
	Run: func(cmd *cobra.Command, args []string) {
		if secretUpdateName == "" || secretUpdateCertFile == "" || namespace == "" {
			cmd.Help()
			return
		}
		//create k8s client
		k8sClient, err := CreateK8sClient()
		if err != nil {
			fmt.Println(err)
			return
		}

		//read https crt key
		certData, err := readFileContent(secretUpdateCertFile)
		if err != nil {
			fmt.Println(err)
			return
		}

		keyData, err := readFileContent(secretUpdatePrivateKeyFile)
		if err != nil {
			fmt.Println(err)
			return
		}

		// 创建新的Secret
		updateSecret := v1.Secret{
			Type: v1.SecretTypeTLS,
			StringData: map[string]string{
				v1.TLSCertKey:       string(certData),
				v1.TLSPrivateKeyKey: string(keyData),
			},
		}

		updateSecret.Name = secretUpdateName
		updateSecret.Namespace = namespace

		//调用update接口方法
		_, err = k8sClient.CoreV1().Secrets(namespace).Update(&updateSecret)

		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("success update !!!")
	},
}

// secret Get 命令
var secretGetName string
var secretGetCmd = cobra.Command{
	Use:   "get",
	Short: "get secret or secret list",
	Run: func(cmd *cobra.Command, args []string) {
		if namespace == "" {
			cmd.Help()
			return
		}

		//create k8s client
		k8sClient, err := CreateK8sClient()
		if err != nil {
			fmt.Println(err)
			return
		}

		//根据secretGetName参数是否为空来决定显示单个secret信息还是所有secret信息
		listOption := metav1.ListOptions{}
		//如果指定了Secret Name 那么只获取单个secret的信息
		if secretGetName != "" {
			listOption.FieldSelector = fmt.Sprintf("metadata.name = %s", secretGetName)
		}
		//调用list接口获取secret 列表
		secretList, err := k8sClient.CoreV1().Secrets(namespace).List(listOption)
		if err != nil {
			fmt.Println(err)
			return
		}

		//遍历Secret List 显示 secret 信息
		printFmt := "%-40s\t%-40s\t%-20s\t%-20s\n"
		fmt.Printf(printFmt, "name", "type", "data", "age")
		for _, secret := range secretList.Items {
			//打印输出
			fmt.Printf(printFmt, secret.Name, secret.Type, strconv.Itoa(len(secret.Data)),
				time.Now().Sub(secret.GetCreationTimestamp().Time))
		}

	},
}

// secret Delete 命令
var secretDeleteName string
var secretDeleteCmd = cobra.Command{
	Use:     "delete",
	Aliases: []string{"del", "rm"},
	Short:   "delete a secret",
	Run: func(cmd *cobra.Command, args []string) {
		if secretDeleteName == "" || namespace == "" {
			cmd.Help()
			return
		}

		k8sClient, err := CreateK8sClient()
		if err != nil {
			fmt.Println(err)
			return
		}

		deleteOption := metav1.DeleteOptions{}
		err = k8sClient.CoreV1().Secrets(namespace).Delete(secretDeleteName, &deleteOption)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("delete success !!!")
	},
}
