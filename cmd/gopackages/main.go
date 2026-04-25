package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func main() {
	goBin := os.Getenv("GO")
	if goBin == "" {
		goBin = "go"
	}

	cmd := exec.Command(goBin, "list", "./...")
	cmd.Stderr = os.Stderr
	output, err := cmd.Output()
	if err != nil {
		fmt.Fprintf(os.Stderr, "list Go packages: %v\n", err)
		os.Exit(1)
	}

	scanner := bufio.NewScanner(bytes.NewReader(output))
	for scanner.Scan() {
		pkg := strings.TrimSpace(scanner.Text())
		if pkg == "" || !includePackage(pkg) {
			continue
		}
		fmt.Println(pkg)
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "read Go package list: %v\n", err)
		os.Exit(1)
	}
}

func includePackage(pkg string) bool {
	normalized := strings.ReplaceAll(pkg, "\\", "/")
	return !strings.Contains(normalized, "/node_modules/")
}
