package main

import (
	"flag"
	"fmt"
	k8s "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8sshell/cmds"
	"os"
	"path/filepath"
)


func main() {
	// 获取根命令
	rootCmd := cmds.GetRootCommand()
	// 执行命令
	err := rootCmd.Execute()
	if err != nil {
		fmt.Println("Err:", err)
		return
	}

	//k8sClient, err := CreateK8sClient()
	//if err != nil {
	//	fmt.Println("Err :", err)
	//	return
	//}
	//listOption := metav1.ListOptions{}
	//l, err :=k8sClient.CoreV1().Pods("kube-system").List(listOption)
	//if err != nil{
	//	fmt.Println(err)
	//}
	//for _, pod := range l.Items{
	//	//fmt.Println(pod.Name)
	//	if pod.Status.Phase == "Running"{
	//		fmt.Printf("%s的状态是：%s\n", pod.Name, pod.Status.Phase)
	//	}
	//
	//}
	//fmt.Println(l.Items)

}

//func createK8SClient() (k8sClient *k8s.Clientset, err error) {
//	cfg := restclient.Config{}
//	cfg.Host = K8SAPIServer
//	cfg.CAData = []byte(K8SCertificateData)
//	cfg.BearerToken = K8SAPIToken
//	cfg.Timeout = time.Second * time.Duration(K8SAPITimeout)
//	k8sClient, err = k8s.NewForConfig(&cfg)
//	return
//}

func CreateK8sClient() (client *k8s.Clientset, err error){
	var kubeconfig *string
	h := homeDir()
	fmt.Println(os.Getenv("HOME"))
	fmt.Println(os.Getenv("USERPROFILE"))

	h = filepath.Join(h, ".kube", "config")
	_, erro := os.Stat(h)

	if erro != nil {
		kubeconfig = flag.String("kubeconfig", "F:\\go_project\\cobra4k8s\\config", "absolute path to the kubeconfig file")
	}else {
		kubeconfig = flag.String("kubeconfig", h, "(optional) absolute path to the kubeconfig file")
	}

	//if  h != ""  {
	//	kubeconfig = flag.String("kubeconfig", h, "(optional) absolute path to the kubeconfig file")
	//} else {
	//	kubeconfig = flag.String("kubeconfig", "F:\\go_project\\cobra4k8s\\config", "absolute path to the kubeconfig file")
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

