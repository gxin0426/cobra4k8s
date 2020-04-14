package cmds

import (
	"fmt"
    "github.com/spf13/cobra"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strings"
    "time"
)

func init(){
	serviceCreateCmd.Flags().StringVar(&serviceCreateName,"name","","service name")
	//更新Service的选项参数
	serviceUpdateCmd.Flags().StringVar(&serviceUpdateName,"name","","service name")
	//获取service的选项参数
	serviceGetCmd.Flags().StringVar(&serviceGetName,"name","","service name")
	//删除service 的选项参数
	serviceDeleteCmd.Flags().StringVar(&serviceDeleteName,"name","","service name")

}

// Service 命令
var serviceRootCmd = cobra.Command{
	Use:     "service",
	Aliases: []string{"svc"},
	Short:   "service is used to manage kubernetes Services",
}

// Service Create 命令
var serviceCreateName string
var serviceCreateCmd = cobra.Command{
	Use:   "create",
	Short: "create a new service",
	Run: func(cmd *cobra.Command, args []string) {
		if serviceCreateName == "" || namespace == ""{
			cmd.Help()
			return
		}
		//创建 kubernetes 客户端对象
		k8sClient, err := createK8SClient()
		if err != nil{
			fmt.Println("Err:", err)
			return
		}

		//新的 Service 的定义
		var newService v1.Service
		var newServiceSpec v1.ServiceSpec
		//设置标签选择器
		newServiceSpec.Selector = map[string]string{
			"app": "echo-go",
		}

		//设置Service 端口
		newServiceSpec.Ports = []v1.ServicePort{
			v1.ServicePort{
				Name: fmt.Sprintf("6666-%s",serviceCreateName),
				Port: 9090,
				TargetPort: intstr.FromInt(9090),
				Protocol: v1.ProtocolTCP,
			},
		}

			//设置 ServiceType 为 NodePort
			newServiceSpec.Type = v1.ServiceTypeNodePort
			//设置Service的各个参数
			newService.Spec = newServiceSpec 
			newService.Name = serviceCreateName
			newService.Namespace = namespace

			//调用Create 接口创建

			_, err = k8sClient.CoreV1().Services(namespace).Create(&newService)
			if err != nil{
				fmt.Println(err)
				return
			}

			fmt.Println("success !!!")
		
	},
}

// Service Update 命令
var serviceUpdateName string

var serviceUpdateCmd = cobra.Command{
	Use:   "update",
	Short: "update a service",
	Run: func(cmd *cobra.Command, args []string) {
		if serviceUpdateName == "" || namespace == ""{
			cmd.Help()
			return
		}

		k8sClient, err := createK8SClient()
		if err != nil {
			fmt.Println("err :", err)
			
		}

		//获取指定的Service 对象
		getOptions := metav1.GetOptions{}
		service, err := k8sClient.CoreV1().Services(namespace).Get(serviceUpdateName, getOptions)
		if err != nil{
			fmt.Println("err:",err)
			return
		}

		//设置原有Service暴露的端口 增加 9091
		service.Spec.Ports = append(service.Spec.Ports, v1.ServicePort{
			Name: fmt.Sprintf("tcp-9091-%s",serviceUpdateName),
			Port: 9091,
			TargetPort: intstr.FromInt(9091),
			Protocol: v1.ProtocolTCP,
		})
		// 调用Update接口创建
		_, err = k8sClient.CoreV1().Services(namespace).Update(service)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("success !!!!")
	},
}

// Service Get 命令
var serviceGetName string
var serviceGetCmd = cobra.Command{
	Use:   "get",
	Short: "get service or service list",
	Run: func(cmd *cobra.Command, args []string) {
		if namespace == ""{
			cmd.Help()
			return
		}
		//创建 k8s客户端对象
		k8sClient, err := createK8SClient()
		if err != nil{
			fmt.Println(err)
			return
		}

		//根据 serviceGetName 参数是否为空来决定显示单个Service信息还是所有Service信息
		listOption := metav1.ListOptions{}
		//如果指定了Service Name 那么只获取单个 Service 的信息
		if serviceGetName != ""{
			listOption.FieldSelector = fmt.Sprintf("metadata.name=%s",serviceGetName)
		}

		//调用List 接口获取Service列表
		serviceList, err := k8sClient.CoreV1().Services(namespace).List(listOption)
		if err != nil{
			fmt.Println(err)
			return
		}

		//遍历Service List 显示 Service 信息

		printFmt := "%-10s\t%-10s\t%-10s\t%-10s\t%-10s\t%-10s\n"
		fmt.Printf(printFmt,"NAME","TYPE","CLUSTER-IP","EXTERNAL-IP","PORT","AGE")
		for _, service := range serviceList.Items{
			//格式化 ServicePort
			servicePorts := make([]string, 0, len(service.Spec.Ports))
			for _, p := range service.Spec.Ports{
				servicePorts = append(servicePorts, fmt.Sprintf("%d:%d/%s", p.Port, p.NodePort, p.Protocol))

			}
			//格式化External IP
			externalIPs := make([]string, 0, len(service.Spec.ExternalIPs))
			for _, ip := range service.Spec.ExternalIPs{
				externalIPs = append(externalIPs, ip)
			}
			var externalIPsStr = "<none>"
			if len(externalIPs) > 0 {
				externalIPsStr = strings.Join(externalIPs, ",")
			}

			//打印输出
			fmt.Printf(printFmt, service.Name, service.Spec.Type, service.Spec.ClusterIP, externalIPsStr,
			strings.Join(servicePorts, ","), time.Now().Sub(service.GetCreationTimestamp().Time).String())

		}


	},
}

// Service Delete 命令
var serviceDeleteName string
var serviceDeleteCmd = cobra.Command{
	Use:     "delete",
	Aliases: []string{"del", "rm"},
	Short:   "delete a service",
	Run: func(cmd *cobra.Command, args []string) {
		if serviceDeleteName == "" || namespace == ""{
			 cmd.Help()
			 return
		}

		k8sClient, err := createK8SClient()
		if err != nil{
			fmt.Println("err", err)
			 return
		}

		//可选的删除选项参数
		deleteOption := metav1.DeleteOptions{}
		//delete service
		err = k8sClient.CoreV1().Services(namespace).Delete(serviceDeleteName, &deleteOption)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("success delete !!!")
	},
}
