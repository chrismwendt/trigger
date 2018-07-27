package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"

	"github.com/fatih/color"
)

var cmd *exec.Cmd

func killIfRunning(cmd *exec.Cmd) {
	if cmd != nil {
		syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
	}
}

func runCommand() {
	blue := color.New(color.FgBlue)
	reset := color.New(color.Reset)

	killIfRunning(cmd)

	name := os.Args[1]
	args := os.Args[2:]
	cmd = exec.Command(name, args...)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	blue.Println("Running " + reset.Sprint(strings.Join(cmd.Args, " ")))
	if err := cmd.Start(); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	go func() {
		if err := cmd.Wait(); err == nil {
			blue.Println("Done")
		} else {
			blue.Println("Done: " + err.Error())
		}
	}()
}

func main() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			if cmd != nil {
				syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
			}
			os.Exit(0)
		}
	}()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			return
		}
		runCommand()
	})
	runCommand()
	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		killIfRunning(cmd)
		log.Fatal("ListenAndServe: ", err)
	}
}
