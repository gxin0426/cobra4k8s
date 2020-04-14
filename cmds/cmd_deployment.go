package cmds

import (
	"github.com/spf13/cobra"
    "fmt"
   appsv1 "k8s.io/api/apps/v1"
    "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strconv"
	"time"
	"strings"
)

func init(){
	//create deploy option param
	deploymentCreateCmd.Flags().StringVar(&deploymentCreateName, "name", "","deployment name")
	deploymentCreateCmd.Flags().StringVar(&deploymentCreateImage, "image", "","image name")
	deploymentCreateCmd.Flags().Int32Var(&deploymentCreateReplicas, "replicas", 0,"replicas count")
	// 更新 Deployment 的选项参数
	deploymentUpdateCmd.Flags().StringVar(&deploymentUpdateName, "name", "", "deployment name")
	deploymentUpdateCmd.Flags().Int32Var(&deploymentUpdateReplicas, "replicas", 0, "replica count")
    // 获取 Deployment 的选项参数
	deploymentGetCmd.Flags().StringVar(&deploymentGetName, "name", "", "deployment name")
	//删除 deployment 的选项参数
	deploymentDeleteCmd.Flags().StringVar(&deploymentDeleteName, "name", "", "deployment name")
	// 升级应用的选项参数
	deploymentUpgradeCmd.Flags().StringVar(&deploymentUpgradeOrRollbackName, "name", "", "deployment name")
	deploymentUpgradeCmd.Flags().StringVar(&deploymentUpgradeOrRollbackImageName, "image", "", "image name")

	// 回滚应用的选项参数
	deploymentRollbackCmd.Flags().StringVar(&deploymentUpgradeOrRollbackName, "name", "", "deployment name")
	deploymentRollbackCmd.Flags().StringVar(&deploymentUpgradeOrRollbackImageName, "image", "", "image name")
}

// Deployment 命令
var deploymentRootCmd = cobra.Command{
	Use:     "deployment",
	Aliases: []string{"deploy"},
	Short:   "deployment is used to manage kubernetes Deployments",
}

// Deployment Create 命令

var deploymentCreateName string
var deploymentCreateImage string
var deploymentCreateReplicas int32
var deploymentCreateCmd = cobra.Command{
	Use:   "create",
	Short: "create a new deployment",
	Run: func(cmd *cobra.Command, args []string) {
		if deploymentCreateName == "" || deploymentCreateImage == "" || namespace == "" || deploymentCreateReplicas == 0 {
            cmd.Help()
            return
        }
        fmt.Println("Creating deployment", deploymentCreateName, "with image", deploymentCreateImage, ", replicas", deploymentCreateReplicas)
        // 组装 DeploymentSpec
        var newDeployment appsv1.Deployment
        var newDeploymentSpec appsv1.DeploymentSpec
        var newDeploymentPodTemplateSpec v1.PodTemplateSpec
        var newDeploymentPodTemplateSpecPodSpec v1.PodSpec

        // 设置 PodTemplateSpec 中的容器相关参数
        newDeploymentPodTemplateSpecPodSpec.Containers = []v1.Container{
            v1.Container{
                // 容器名称
                Name: "echo-go-test",
                // 镜像名称
                Image: deploymentCreateImage,
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
        // 设置 PodTemplateSpec 的 PodSpec 和标签
        newDeploymentPodTemplateSpec.Spec = newDeploymentPodTemplateSpecPodSpec
        newDeploymentPodTemplateLabels := map[string]string{
            "app": "echo-go-test",
        }
        newDeploymentPodTemplateSpec.Labels = newDeploymentPodTemplateLabels

        // 设置 Deployment Spec 的 Template
        newDeploymentSpec.Template = newDeploymentPodTemplateSpec

        // 设置 Deployment Spec 的 Replicas
        newDeploymentSpec.Replicas = &deploymentCreateReplicas

        // 设置 Deployment Spec 的 Selectors，匹配 Templates 的标签
        newDeploymentSpec.Selector = &metav1.LabelSelector{
            MatchLabels: newDeploymentPodTemplateLabels,
        }

        // 设置 Deployment 名称
        newDeployment.Name = deploymentCreateName

        // 设置 Deployment 的命名空间
        newDeployment.Namespace = namespace

        //  设置 Deployment 的 Spec
        newDeployment.Spec = newDeploymentSpec

        // 创建 Kubernetes 的客户端
        k8sClient, err := createK8SClient()
        if err != nil {
            fmt.Println("Err:", err)
            return
        }

        // 调用 Create 的接口方法
        _, err = k8sClient.AppsV1().Deployments(namespace).Create(&newDeployment)
        if err != nil {
            fmt.Println("Err:", err)
            return
        }
        fmt.Println("Create new deployment success!")

	},
}

// Deployment Update 命令
var deploymentUpdateName string
var deploymentUpdateReplicas int32
var deploymentUpdateCmd = cobra.Command{
	Use:   "update",
	Short: "update a deployment",
	Run: func(cmd *cobra.Command, args []string) {
		if deploymentUpdateName == "" || namespace == "" || deploymentUpdateReplicas == 0 {
			cmd.Help()
			return
		}
		fmt.Println("Updating deployment", deploymentUpdateName, "to replicas", deploymentUpdateReplicas)

		// 创建 Kubernetes 的客户端
		k8sClient, err := createK8SClient()
		if err != nil {
			fmt.Println("Err:", err)
			return
		}

		// 根据 Deployment 获取 Deployment
		getOption := metav1.GetOptions{}
		deployment, err := k8sClient.AppsV1().Deployments(namespace).Get(deploymentUpdateName, getOption)
		if err != nil {
			fmt.Println("Err:", err)
			return
		}

		// 设置 Deployment 的 Replica 数量
		deployment.Spec.Replicas = &deploymentUpdateReplicas

		// 调用 Update 接口更新 Deployment
		_, err = k8sClient.AppsV1().Deployments(namespace).Update(deployment)
		if err != nil {
			fmt.Println("Err:", err)
			return
		}
		fmt.Println("Update deployment success!")
	},
}

// Deployment Get 命令
var deploymentGetName string
var deploymentGetCmd = cobra.Command{
	Use:   "get",
	Short: "get deployment or deployment list",
	Run: func(cmd *cobra.Command, args []string) {
        if namespace == "" {
            cmd.Help()
            return
        }

        // 创建 Kubernetes 客户端对象
        k8sClient, err := createK8SClient()
        if err != nil {
            fmt.Println("Err:", err)
            return
        }

        // 根据 deploymentGetName 参数是否为空来决定显示单个 Deployment 信息还是所有 Deployment 信息
        listOption := metav1.ListOptions{}
        // 如果指定了 Deployment Name，那么只获取单个 Deployment 的信息
        if deploymentGetName != "" {
            listOption.FieldSelector = fmt.Sprintf("metadata.name=%s", deploymentGetName)
        }

        // 调用 List 接口获取 Deployment 列表
        deploymentList, err := k8sClient.AppsV1().Deployments(namespace).List(listOption)
        if err != nil {
            fmt.Println("Err:", err)
            return
        }

        // 遍历 Deployment List，显示 Deployment 信息
        printFmt := "%-10s\t%-10s\t%-10s\t%-10s\t%-10s\t%-10s\n"
        fmt.Printf(printFmt, "NAME", "DESIRED", "CURRENT", "UP-TO-DATE", "AVAILABLE", "AGE")
        for _, deployment := range deploymentList.Items {
            fmt.Printf(printFmt, deployment.Name, strconv.Itoa(int(*deployment.Spec.Replicas)), strconv.Itoa(int(deployment.Status.Replicas)),
                strconv.Itoa(int(deployment.Status.UpdatedReplicas)), strconv.Itoa(int(deployment.Status.AvailableReplicas)),
				time.Now().Sub(deployment.GetCreationTimestamp().Time))
		}
	},
}

// Deployment Delete 命令
var deploymentDeleteName string
var deploymentDeleteCmd = cobra.Command{
	Use:     "delete",
	Aliases: []string{"del", "rm"},
	Short:   "delete a deployment",
	Run: func(cmd *cobra.Command, args []string) {
		if deploymentDeleteName == "" || namespace == "" {
			cmd.Help()
			return
		}

		//创建客户端
		k8sClient, err := createK8SClient()
		if err != nil {
			 fmt.Println(err)
			 return
		}

		//可选的删除选项参数
		//这个删除策略很重要 设置为background 可以将Deployment 线面的replicasets也删除 同事也是删了相关的pod

		var deletePolicy = metav1.DeletePropagationBackground
		deleteOption := metav1.DeleteOptions{
			PropagationPolicy: &deletePolicy,
		}
		//删除
		err = k8sClient.AppsV1().Deployments(namespace).Delete(deploymentDeleteName, &deleteOption)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("delete success !!!")
	},
}


// Deployment Upgrade 命令
var deploymentUpgradeCmd = createDeploymentUpgradeOrRollbackCommand("upgrade")

// Deployment Rollback 命令
var deploymentRollbackCmd = createDeploymentUpgradeOrRollbackCommand("rollback")

var deploymentUpgradeOrRollbackName string
var deploymentUpgradeOrRollbackImageName string

func createDeploymentUpgradeOrRollbackCommand(cmdName string) cobra.Command{
	return cobra.Command{
		Use: cmdName,
		Short: fmt.Sprintf("%s a deployment",cmdName),
		Run: func(cmd *cobra.Command, args []string){
			if deploymentUpgradeOrRollbackName == "" || deploymentUpgradeOrRollbackImageName == "" || namespace == ""{
				cmd.Help()
				return
			}

			fmt.Println(strings.ToTitle(cmdName), "deployment", deploymentUpgradeOrRollbackName, "to image", deploymentUpgradeOrRollbackImageName)
			//create k8s client
			k8sClient, err := createK8SClient()
			if err != nil {
				 fmt.Println(err)
				  return
			}
			//根据deployment获取deployment
			getOption := metav1.GetOptions{}
			deployment, err := k8sClient.AppsV1().Deployments(namespace).Get(deploymentUpgradeOrRollbackName, getOption)
			if err != nil{
				fmt.Println(err)
				return
			}

			//设置新的镜像
			deployment.Spec.Template.Spec.Containers = []v1.Container{
				v1.Container{
					Name: "echo-go",
					Image: deploymentUpgradeOrRollbackImageName,
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
	_, err = k8sClient.AppsV1().Deployments(namespace).Update(deployment)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(strings.ToTitle(cmdName),"deployment success")
},
	}
}