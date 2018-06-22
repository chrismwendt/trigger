package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"strings"

	"github.com/fatih/color"
)

var cmd *exec.Cmd

func runCommand() {
	blue := color.New(color.FgBlue)
	reset := color.New(color.Reset)

	if cmd != nil {
		cmd.Process.Kill()
	}

	name := os.Args[1]
	args := os.Args[2:]
	cmd = exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	blue.Println("Running " + reset.Sprint(strings.Join(cmd.Args, " ")))
	if err := cmd.Start(); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	go func() {
		if err := cmd.Wait(); err == nil {
			blue.Println("Done")
		}
	}()
}

func main() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			if cmd != nil {
				cmd.Process.Kill()
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
		log.Fatal("ListenAndServe: ", err)
	}
}
