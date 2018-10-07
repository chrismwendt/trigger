package main

import (
	"flag"
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

var clear = flag.Bool("clear", false, "clear the terminal each time")
var cmdParts []string

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

	cmdParts := flag.Args()
	if *clear {
		clear := exec.Command("clear")
		clear.Stdout = os.Stdout
		clear.Start()
		clear.Wait()
	}
	cmd = exec.Command(cmdParts[0], cmdParts[1:]...)
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
	flag.Parse()
	cmdParts = flag.Args()

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
	err := http.ListenAndServe(":7416", nil)
	if err != nil {
		killIfRunning(cmd)
		log.Fatal("ListenAndServe: ", err)
	}
}
