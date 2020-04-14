package main

import (
	"k8sshell/cmds"
	"fmt"
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
}
