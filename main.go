package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"github.com/alessio/shellescape"
	"github.com/fatih/color"
)

var clear = flag.Bool("clear", false, "clear the terminal each time")
var cmdParts []string

var cmd *exec.Cmd
var ready = make(chan struct{}, 1)

func kill() {
	if cmd != nil && cmd.ProcessState == nil && cmd.Process != nil {
		syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
	}
	<-ready
}

func spawn() {
	kill()

	cmdParts := flag.Args()
	if *clear {
		clear := exec.Command("tput", "reset")
		clear.Stdout = os.Stdout
		clear.Start()
		clear.Wait()
	}
	cmd = exec.Command(cmdParts[0], cmdParts[1:]...)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	shellCommand := shellescape.QuoteCommand(cmd.Args)

	fmt.Println(color.BlueString("trigger") + " " + color.CyanString(shellCommand))
	if err := cmd.Start(); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	go func() {
		err := cmd.Wait()
		fmt.Print(color.BlueString("trigger"))
		if err != nil {
			fmt.Print(" " + color.RedString(err.Error()))
		} else {
			fmt.Print(" " + color.BlueString("exit code 0"))
		}
		fmt.Println()
		ready <- struct{}{}
	}()
}

func main() {
	flag.Parse()
	cmdParts = flag.Args()

	ready <- struct{}{}

	// Clean up on Ctrl+C
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		kill()
		os.Exit(0)
	}()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		spawn()
	})
	spawn()
	err := http.ListenAndServe(":7416", nil)
	if err != nil {
		kill()
		log.Fatal("ListenAndServe: ", err)
	}
}
