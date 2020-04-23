package cmds

import (
	"fmt"
	"strconv"
	"time"
	k8s "k8s.io/client-go/kubernetes"
	"github.com/spf13/cobra"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)


	var k8sclient *k8s.Clientset
// Pod 命令
var podRootCmd = cobra.Command{
	Use:   "pod",
	Short: "pod is used to manage kubernetes Pods",
}

func init() {
	// 创建 Pod 的选项参数
	k8sclient, _ = CreateK8sClient()
	podCreateCmd.Flags().StringVar(&podCreateName, "name", "", "pod name")
	podCreateCmd.Flags().StringVar(&podCreateImage, "image", "", "image name")
	podGetCmd.Flags().StringVar(&podGetName, "name", "", "pod name")
	podUpdateCmd.Flags().StringVar(&podUpdateName, "name", "", "pod name")
	podUpdateCmd.Flags().StringVar(&podUpdateImage, "image", "", "image name")
	podDeleteCmd.Flags().StringVar(&podDeleteName, "name", "", "pod name")
	podCheckCmd.Flags().StringVar(&podGetName, "name", "", "pod name")


}

//pod check 命令

var podCheckName string
var podCheckCmd = cobra.Command{
	Use:   "check",
	Short: "check pod",
	Run: func(cmd *cobra.Command, args []string) {
		if namespace == "" {
			cmd.Help()
			return
		}
		// 创建 Kubernetes 客户端对象

		// 根据 podGetName 参数是否为空来决定显示单个 Pod 信息还是所有 Pod 信息
		listOption := metav1.ListOptions{}
		// 如果指定了 Pod Name，那么只获取单个 Pod 的信息
		if podGetName != "" {
			listOption.FieldSelector = fmt.Sprintf("metadata.name=%s", podGetName)
		}

		// 调用 List 接口获取 Pod 列表
		if namespace == "all" {
			namespace = ""
		}
		podList, err := k8sclient.CoreV1().Pods(namespace).List(listOption)
		if err != nil {
			fmt.Println("Err:", err)
			return
		}
		if podCheckName != "" {
			listOption.FieldSelector = fmt.Sprintf("metadata.name=%s", podCheckName)
		}

		formatPrint := "%-50s\t%s\n"
		fmt.Printf(formatPrint, "NAME", "STATUS")
		for _, pod := range podList.Items {

			if pod.Status.Phase == "Running" {
				//fmt.Println(pod.Status.Phase)

				fmt.Printf(formatPrint, pod.Name, pod.Status.Phase)
			}
		}

	},
}

//pod get 命令
var podGetName string
var podGetCmd = cobra.Command{
	Use:   "get",
	Short: "get pod or pod list",
	Run: func(cmd *cobra.Command, args []string) {
		if namespace == "" {
			cmd.Help()
			return
		}

		// 创建 Kubernetes 客户端对象
		//k8sClient, err := CreateK8sClient()
		//if err != nil {
		//	fmt.Println("Err:", err)
		//	return
		//}

		// 根据 podGetName 参数是否为空来决定显示单个 Pod 信息还是所有 Pod 信息
		listOption := metav1.ListOptions{}
		// 如果指定了 Pod Name，那么只获取单个 Pod 的信息
		if podGetName != "" {
			listOption.FieldSelector = fmt.Sprintf("metadata.name=%s", podGetName)
		}

		// 调用 List 接口获取 Pod 列表
		if namespace == "all" {
			namespace = ""
		}
		podList, err := k8sclient.CoreV1().Pods(namespace).List(listOption)
		if err != nil {
			fmt.Println("Err:", err)
			return
		}

		/********************test************************/

		//formatPrint := "%-50s\t%s\n"
		//fmt.Printf(formatPrint,"NAME", "STATUS")
		//for _, pod := range podList.Items {
		//
		//	if pod.Status.Phase == "Running" {
		//		//fmt.Println(pod.Status.Phase)
		//
		//		fmt.Printf(formatPrint, pod.Name, pod.Status.Phase)
		//	}
		//}

		/********************test************************/

		//遍历 Pod List，显示 Pod 信息
		printFmt := "%-30s\t%-10s\t%-10s\t%-10s\t%s\n"
		fmt.Printf(printFmt, "NAME", "READY", "STATUS", "RESTARTS", "AGE")
		for _, pod := range podList.Items {
			// 计算 Container Ready 的数量
			containerAllCount := len(pod.Status.ContainerStatuses)
			containerReadyCount := 0
			for _, cs := range pod.Status.ContainerStatuses {
				if cs.State.Running != nil {
					containerReadyCount++
				}
			}

			//打印输出
			fmt.Printf(printFmt,
				pod.Name,
				strconv.Itoa(containerReadyCount)+"/"+strconv.Itoa(containerAllCount),
				pod.Status.Phase,
				strconv.Itoa(int(pod.Status.ContainerStatuses[0].RestartCount)),
				time.Now().Sub(pod.Status.StartTime.Time).String())
		}
	},
}

//  Pod Create 命令的选项参数
var podCreateName string
var podCreateImage string

// Pod Create 命令
var podCreateCmd = cobra.Command{
	Use:   "create",
	Short: "create a new pod",
	Run: func(cmd *cobra.Command, args []string) {
		if podCreateName == "" || podCreateImage == "" || namespace == "" {
			cmd.Help()
			return
		}
		fmt.Println("Creating pod", podCreateName, "with image", podCreateImage)
		// 组装 PodSpec
		var newPod v1.Pod
		var newPodSpec v1.PodSpec

		// 设置 Pod 中的容器相关参数
		newPodSpec.Containers = []v1.Container{
			v1.Container{
				// 容器名称
				Name: "echo-go",
				// 镜像名称
				Image: podCreateImage,
				// 工作目录
				WorkingDir: "/home/app/",
				// 容器启动命令
				Command: []string{"/home/app/echo-go"},
				// 容器启动参数
				Args: []string{"-port", "9090"},
				// 容器暴露端口
				Ports: []v1.ContainerPort{
					v1.ContainerPort{
						Name:          "http",
						ContainerPort: 9090,
						Protocol:      v1.ProtocolTCP,
					},
				},
			},
		}

		// 设置 Pod 名称
		newPod.Name = podCreateName
		// 设置 Pod 标签
		newPod.Labels = map[string]string{
			"app": "echo-go",
		}
		// 设置 Pod 的命名空间
		newPod.Namespace = namespace
		newPod.Spec = newPodSpec

		// 创建 Kubernetes 的客户端
		//k8sClient, err := CreateK8sClient()
		//if err != nil {
		//	fmt.Println("Err:", err)
		//	return
		//}

		// 调用 Create 的接口方法
		_, err := k8sclient.CoreV1().Pods(namespace).Create(&newPod)
		if err != nil {
			fmt.Println("Err:", err)
			return
		}
		fmt.Println("Create new pod success!")
	},
}

//pod update 命令
var podUpdateName string
var podUpdateImage string

// Pod Update 命令
var podUpdateCmd = cobra.Command{
	Use:   "update",
	Short: "update a pod",
	Run: func(cmd *cobra.Command, args []string) {
		if podUpdateName == "" || podUpdateImage == "" || namespace == "" {
			cmd.Help()
			return
		}
		fmt.Println("Updating pod", podUpdateName, "with image", podUpdateImage)
		//k8sClient, err := CreateK8sClient()
		//if err != nil {
		//	fmt.Println("Err :", err)
		//	return
		//}

		//删除旧的pod
		var deleteGracePeriodSeconds int64 = 0
		deleteOption := metav1.DeleteOptions{
			//设置宽限时间为0，立刻删除pod
			GracePeriodSeconds: &deleteGracePeriodSeconds,
		}
		err := k8sclient.CoreV1().Pods(namespace).Delete(podUpdateName, &deleteOption)
		if err != nil {
			fmt.Println("Err :", err)
			return
		}

		//创建新的pod
		var updatePod v1.Pod
		var updatePodSpec v1.PodSpec
		//设置pod中的容器相关参数
		updatePodSpec.Containers = []v1.Container{
			v1.Container{
				Name:       "echo-go",
				Image:      podUpdateImage,
				WorkingDir: "/home/app/",
				Command:    []string{"/home/app/echo-go"},
				Args:       []string{"-port", "9090"},
				Ports: []v1.ContainerPort{
					v1.ContainerPort{
						Name:          "http",
						ContainerPort: 9090,
						Protocol:      v1.ProtocolTCP,
					},
				},
			},
		}
		//设置pod名称
		updatePod.Name = podUpdateName
		updatePod.Labels = map[string]string{
			"app": "echo-go",
		}

		//设置pod的命名空间
		updatePod.Namespace = namespace
		updatePod.Spec = updatePodSpec

		//调用Create 的接口方法
		_, err = k8sclient.CoreV1().Pods(namespace).Create(&updatePod)
		if err != nil {
			fmt.Println("err :", err)
			return
		}

		fmt.Println("update success !!!")

	},
}

// Pod Delete 命令
var podDeleteName string
var podDeleteCmd = cobra.Command{
	Use:     "delete",
	Aliases: []string{"del", "rm"},
	Short:   "delete a pod",
	Run: func(cmd *cobra.Command, args []string) {

		if podDeleteName == "" || namespace == "" {
			cmd.Help()
			return
		}

		//创建k8s 客户端对象
		k8sClient, err := CreateK8sClient()
		if err != nil {
			fmt.Println("Err :", err)
			return
		}
		//可选的删除选项参数
		deleteOption := metav1.DeleteOptions{}
		//删除 pod
		err = k8sClient.CoreV1().Pods(namespace).Delete(podDeleteName, &deleteOption)
		if err != nil {
			fmt.Println("Err :", err)
			return
		}

		fmt.Println("delete pod success !!!")

	},
}
