package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"os/user"

	"github.com/spf13/pflag"
)

var (
	hostname *string
	Version  string
	version  *bool
	// show     *bool
	help *bool
)

func main() {
	flags()
	fmt.Printf("hn version %s\n\n", Version)
	if *version {
		os.Exit(0)
	}

	if *hostname == "" {
		fmt.Println("Current values")
		show()
		os.Exit(0)
	}

	checkAdmin()

	if len(*hostname) > 1 {
		do("hostname", []string{*hostname})
		do("scutil", []string{"--set", "LocalHostName", *hostname})
		do("scutil", []string{"--set", "HostName", *hostname})
		do("scutil", []string{"--set", "ComputerName", *hostname})
		fmt.Println("Set new hostname")
		show()
	}
}

func flags() {

	version = pflag.BoolP("version", "v", false, "show version")
	// show = pflag.BoolP("show", "s", false, "show current values")
	hostname = pflag.StringP("hostname", "n", "", "set hostname")
	help = pflag.BoolP("help", "h", false, "show help")

	pflag.Parse()
}

func show() {
	fmt.Printf("hostname: %s\n", do("hostname", nil))
	fmt.Printf("scutil LocalHostName: %s\n", do("scutil", []string{"--get", "LocalHostName"}))
	fmt.Printf("scutil HostName: %s\n", do("scutil", []string{"--get", "HostName"}))
	fmt.Printf("scutil ComputerName: %s\n", do("scutil", []string{"--get", "ComputerName"}))
}

func do(command string, arg []string) string {
	var err error
	var cmd *exec.Cmd
	var out bytes.Buffer
	var stderr bytes.Buffer

	cmd = exec.Command(command, arg...)
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err = cmd.Run()
	if err != nil {
		fmt.Println(err)
		fmt.Println(stderr.String())
		if stderr.String() != "HostName: not set" {
			// fmt.Printf("Error occured. %s not set?\n\n", arg)
			fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
			// fmt.Printf("cmd.Output %s, %s\n\n", arg, err.Error())

			return ""
		}
	}

	// fmt.Printf("command output %s", string(out))
	return out.String()
}

func checkAdmin() {
	user, err := user.Current()
	if err != nil {
		fmt.Printf("user.Current error %s", err.Error())
		os.Exit(1)
	}

	if user.Uid != "0" {
		fmt.Printf("Hi %s. Run this app with root privileges.", user.Username)
		os.Exit(1)
	}
}
