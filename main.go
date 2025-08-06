package main

import (
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
	fmt.Printf("hostname: %s", do("hostname", nil))
	fmt.Printf("scutil LocalHostName: %s", do("scutil", []string{"--get", "LocalHostName"}))
	fmt.Printf("scutil HostName: %s", do("scutil", []string{"--get", "HostName"}))
	fmt.Printf("scutil ComputerName: %s", do("scutil", []string{"--get", "ComputerName"}))
}

func do(command string, arg []string) string {
	var cmd *exec.Cmd
	var out []byte
	var err error

	cmd = exec.Command(command, arg...)
	out, err = cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("cmd.CombinedOutput %s, %s", arg, err.Error())

		return ""
	}

	// fmt.Printf("command output %s", string(out))
	return string(out)
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
