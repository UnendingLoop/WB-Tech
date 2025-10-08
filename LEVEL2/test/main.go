package main

import (
	"io"
	"os"
	"os/exec"
)

func main() {
	r, w := io.Pipe()
	out, _ := os.Create("out.txt")

	// echo hello world
	cmd1 := exec.Command("cmd", "/c", "echo hello world")
	cmd1.Stdout = w

	// tr a-z A-Z
	cmd2 := exec.Command("tr", "a-z", "A-Z")
	cmd2.Stdin = r
	cmd2.Stdout = out

	cmd1.Start()
	cmd2.Start()
	cmd1.Wait()
	w.Close() // Закрыть после завершения echo!
	cmd2.Wait()
	out.Close()
}
