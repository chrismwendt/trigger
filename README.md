# `trigger` reruns a command on keystroke

Conveniently rerun any shell script on keystroke. It's useful for running build commands on demand when programming, especially when you don't have a file watcher or the file watcher is buggy or uses too much CPU.

`trigger` tries extra hard to force kill the subprocess group to prevent the accumulation of zombie processes.

## Installation

```
go get -u github.com/chrismwendt/trigger
```

Register a keybinding in your desktop automation tool of choice. For [Hammerspoon](https://www.hammerspoon.org/):

```lua
hs.hotkey.bind({"cmd"}, "'", function()
  hs.execute("/usr/local/bin/timeout 1s curl localhost:7416")
end)
```

## Usage

Then run:

```
$ trigger go run main.go
Running go run main.go
hello, world
Done
```

While it's is running, hit <kbd>Cmd+'</kbd> to run `go run main.go` again:

```
$ trigger go run main.go
Running go run main.go
hello, world
Done
Running go run main.go
hello, world
Done
```
