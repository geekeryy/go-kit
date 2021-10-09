package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var root = &cobra.Command{}

var newProject = &cobra.Command{
	Use:   "new",
	Short: "n",
	Long:  "new project name",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			log.Fatalln("invalid param", args)
		}
		Run(args[0])
	},
}

func main() {
	root.AddCommand(newProject)
	if err := root.Execute(); err != nil {
		log.Fatalln(err)
	}
}

func Run(serverName string) {
	if err := Clone("git@github.com:comeonjy/go-layout.git", serverName); err != nil {
		log.Fatalln("Clone", err)
	}

	if err := Make(serverName); err != nil {
		log.Fatalln("Make", err)
	}

	if err := Clear(serverName); err != nil {
		log.Fatalln("Clear", err)
	}

	fmt.Println(color.WhiteString("$ cd %s", serverName))
	fmt.Println(color.WhiteString("$ go generate ."))
	fmt.Println(color.WhiteString("$ go build "))
	fmt.Println(color.WhiteString("$ ./%s\n", serverName))
	fmt.Println("			🤝 Thanks for using go-kit")
	fmt.Println("	📚 Tutorial: https://github.com/comeonjy/go-layout")
}

// Clone 克隆框架模板
func Clone(url string, serverName string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, "git", "clone", url, serverName)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

// Make 执行初始化脚本
func Make(serverName string) error {
	cmd := exec.Command("make", "-C", serverName, "server_name="+serverName, "kit")
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

// Clear 清除初始化脚本
// sed -i '' -e '/kit/,$d' my-server/Makefile
func Clear(serverName string) error {
	cmd := exec.Command(`sed`, `-i`, `-e`, `/kit/,$d`, "./"+serverName+"/Makefile")
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}
