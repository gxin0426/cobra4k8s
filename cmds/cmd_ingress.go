package cmds

import (
    "fmt"
    "k8s.io/apimachinery/pkg/util/intstr"
	v1beta1 "k8s.io/api/extensions/v1beta1"
	"github.com/spf13/cobra"
	"strconv"
    "strings"
    "time"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func init(){
	//创建 ingress 的选项参数
	ingressCreateCmd.Flags().StringVar(&ingressCreateName,"name","","ingress name")
	ingressCreateCmd.Flags().StringVar(&ingressCreateServiceName,"service-name","","service name")
	//更新Ingress的选项参数
	ingressUpdateCmd.Flags().StringVar(&ingressUpdateName,"name","n","ingress name")
	ingressUpdateCmd.Flags().StringVar(&ingressUpdateServiceName,"service-name","","service name")
	//obtain ingress option parameter
	ingressGetCmd.Flags().StringVar(&ingressGetName, "name", "", "ingress name")
	//删除 ingress 的可选参数
	ingressDeleteCmd.Flags().StringVar(&ingressDeleteName, "name", "", "ingress name")
	// 创建 Ingress 的选项参数（支持https）
    ingressCreateHTTPSCmd.Flags().StringVar(&ingressCreateHTTPSName, "name", "", "ingress name")
    ingressCreateHTTPSCmd.Flags().StringVar(&ingressCreateHTTPSSecretName, "secret-name", "", "secret name")
    ingressCreateHTTPSCmd.Flags().StringVar(&ingressCreateHTTPSServiceName, "service-name", "", "service name")
}


// Ingress 命令
var ingressRootCmd = cobra.Command{
	Use:     "ingress",
	Aliases: []string{"ing"},
	Short:   "ingress is used to manage kubernetes Ingresses",
}

// Ingress Create 命令
var ingressCreateName string
var ingressCreateServiceName string
var ingressCreateCmd = cobra.Command{
	Use:   "create",
	Short: "create a new ingress",
	Run: func(cmd *cobra.Command, args []string) {
		if ingressCreateName == "" || ingressCreateServiceName == "" || namespace == ""{
			cmd.Help()
			return
		}
		fmt.Println("Creating ingress", ingressCreateName, "for service", ingressCreateServiceName)
		var newIngress v1beta1.Ingress

		//组装 ingressSpec
		var newIngressSpec v1beta1.IngressSpec
		newIngressSpec.Rules = []v1beta1.IngressRule{
			v1beta1.IngressRule{
				Host: "echo-go.gree.com",
				IngressRuleValue: v1beta1.IngressRuleValue{
					HTTP: &v1beta1.HTTPIngressRuleValue{
						Paths: []v1beta1.HTTPIngressPath{
							v1beta1.HTTPIngressPath{
								Path: "/",
								Backend: v1beta1.IngressBackend{
									ServiceName: ingressCreateServiceName,
									ServicePort: intstr.FromInt(9090),
								},
							},
						},
					},
				},
			},
		}
		//设置 ingress shuxing
		newIngress.Spec = newIngressSpec
		newIngress.Name = ingressCreateName

		//创建 k8s客户端
		k8sClient, err := createK8SClient()
		if err != nil{
			fmt.Println(err)
			return
		}

		//创建ingress
		_, err = k8sClient.ExtensionsV1beta1().Ingresses(namespace).Create(&newIngress)
		if err != nil{
			fmt.Println(err)
			return
		}
		fmt.Println("create success")
	},
}

// Ingress Update 命令
var ingressUpdateName string
var ingressUpdateServiceName string
var ingressUpdateCmd = cobra.Command{
	Use:   "update",
	Short: "update a ingress",
	Run: func(cmd *cobra.Command, args []string) {
		if ingressUpdateName == "" || ingressUpdateServiceName == "" || namespace == ""{
			cmd.Help()
			return
		}
		fmt.Println("Updating ingress", ingressUpdateName, "for service", ingressUpdateServiceName)
		var updateIngress v1beta1.Ingress
		var updateIngressSpec v1beta1.IngressSpec
		var ingressRuleValue =  v1beta1.IngressRuleValue{
			HTTP: &v1beta1.HTTPIngressRuleValue{
				Paths: []v1beta1.HTTPIngressPath{
					v1beta1.HTTPIngressPath{
						Path: "/",
						Backend: v1beta1.IngressBackend{
							ServiceName: ingressUpdateServiceName,
							ServicePort: intstr.FromInt(9090),
						},
					},
				},
			},
		}
		updateIngressSpec.Rules = []v1beta1.IngressRule{
			v1beta1.IngressRule{
				Host: "echo-go.gree.com",
				IngressRuleValue: ingressRuleValue,
			},
			//增加一个域名
			v1beta1.IngressRule{
				Host: "echo-go.gree.io",
				IngressRuleValue: ingressRuleValue,
			},

		}
		//设置 ingress 属性
		updateIngress.Spec = updateIngressSpec
		updateIngress.Name = ingressUpdateName

		//establish k8s client
		k8sClient, err := createK8SClient()
		if err != nil {
			fmt.Println(err)
			return
		}

		//establish ingress
		_, err = k8sClient.ExtensionsV1beta1().Ingresses(namespace).Update(&updateIngress)

		if err != nil{
			fmt.Println(err)
			return
		}
		fmt.Println("update success !!!")

	},
}

// Ingress Get 命令
var ingressGetName string
var ingressGetCmd = cobra.Command{
	Use:   "get",
	Short: "get ingress or ingress list",
	Run: func(cmd *cobra.Command, args []string) {
		if namespace == ""{
			cmd.Help()
			return
		}

		//establish k8s client
		k8sClient, err := createK8SClient()
		if err != nil {
			fmt.Println(err)
			return
		}

		//根据ingressgetname 参数是否为空来决定显示单个ingress 信息还是所有ingress信息
		listOption := metav1.ListOptions{}
		// 如果指定了ingress name n那么只获取单个ingress 信息
		if ingressGetName != ""{
			listOption.FieldSelector = fmt.Sprintf("metadata.name=%s", ingressGetName)
		}
		//调用List接口获取Ingress列表
		ingressList, err := k8sClient.ExtensionsV1beta1().Ingresses(namespace).List(listOption)
		if err != nil{
			fmt.Println(err)
			return
		}

		//遍历 ingress list 显示 ingress 信息

		printFmt := "%-10s\t%-50s\t%-10s\t%-10s\t%s\n"
		fmt.Printf(printFmt, "name", "hosts", "address", "ports", "age")
		for _, ingress := range ingressList.Items{
			//获取hosts
			ingressHosts := make([]string, 0, len(ingress.Spec.Rules))
			for _, rule := range ingress.Spec.Rules{
				ingressHosts = append(ingressHosts, rule.Host)
			}
			//获取address
			ingressAddress := make([]string, 0, len(ingress.Spec.Rules))
			for _, ingStatus := range ingress.Status.LoadBalancer.Ingress{
				ingressAddress = append(ingressAddress, ingStatus.IP)
			}
			//获取Service的端口
			servicePortsSet := make(map[int]struct{})
			for _, rule := range ingress.Spec.Rules{
				for _, path := range rule.IngressRuleValue.HTTP.Paths{
					servicePortsSet[path.Backend.ServicePort.IntValue()] = struct{}{}
				}
			}
			servicePorts := make([]string, 0, len(servicePortsSet))
			for port := range servicePortsSet{
				servicePorts = append(servicePorts, strconv.Itoa(port))
			}

			fmt.Printf(printFmt,ingress.Name, strings.Join(ingressHosts, ","),
			strings.Join(ingressAddress, ","), strings.Join(servicePorts, ","),
			time.Now().Sub(ingress.GetCreationTimestamp().Time))
		}
	},
}

// Ingress Delete 命令
var ingressDeleteName string
var ingressDeleteCmd = cobra.Command{
	Use:     "delete",
	Aliases: []string{"del", "rm"},
	Short:   "delete a ingress",
	Run: func(cmd *cobra.Command, args []string) {
		if ingressDeleteName == "" || namespace == ""{
			cmd.Help()
			return
		}
		// establish k8s client
		k8sClient, err := createK8SClient()
		if err != nil {
			fmt.Println(err)
			return
		} 

		//可选的删除选项参数
		deleteOption := metav1.DeleteOptions{}
		//delete po
		err = k8sClient.ExtensionsV1beta1().Ingresses(namespace).Delete(ingressDeleteName, &deleteOption)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println("success deleting  !!!")
	},
}


// Ingress Create-HTTPS 命令
var ingressCreateHTTPSName string
var ingressCreateHTTPSSecretName string
var ingressCreateHTTPSServiceName string
var ingressCreateHTTPSCmd = cobra.Command{
    Use:   "create-https",
    Short: "create a new ingress which supports https",
    Run: func(cmd *cobra.Command, args []string) {
        if ingressCreateHTTPSName == "" || ingressCreateHTTPSServiceName == "" || ingressCreateHTTPSSecretName == "" || namespace == "" {
            cmd.Help()
            return
        }
        fmt.Println("Creating ingress", ingressCreateHTTPSName, "for service", ingressCreateHTTPSServiceName, "with https support")

        var newIngress v1beta1.Ingress

        // 组装 IngressSpec
        var newIngressSpec v1beta1.IngressSpec
        newIngressSpec.Rules = []v1beta1.IngressRule{
            v1beta1.IngressRule{
                Host: "echo-go-https.shiyanlou.com",
                IngressRuleValue: v1beta1.IngressRuleValue{
                    HTTP: &v1beta1.HTTPIngressRuleValue{
                        Paths: []v1beta1.HTTPIngressPath{
                            v1beta1.HTTPIngressPath{
                                Path: "/",
                                Backend: v1beta1.IngressBackend{
                                    ServiceName: ingressCreateHTTPSServiceName,
                                    ServicePort: intstr.FromInt(9090),
                                },
                            },
                        },
                    },
                },
            },
        }

        // 设置 TLS 字段
        newIngressSpec.TLS = []v1beta1.IngressTLS{
            v1beta1.IngressTLS{
                Hosts:      []string{"echo-go-https.shiyanlou.com"},
                SecretName: ingressCreateHTTPSSecretName,
            },
        }

        // 设置 Ingress 属性
        newIngress.Spec = newIngressSpec
        newIngress.Name = ingressCreateHTTPSName

        // 创建 Kubernetes 的客户端
        k8sClient, err := createK8SClient()
        if err != nil {
            fmt.Println("Err:", err)
            return
        }

        // 创建 Ingress
        _, err = k8sClient.ExtensionsV1beta1().Ingresses(namespace).Create(&newIngress)
        if err != nil {
            fmt.Println("Err:", err)
            return
        }
        fmt.Println("Create new ingress success!")
    },
}