//go:build mage
// +build mage

package main

import (
	"fmt"
	"os"
	"os/exec"
)

func isPodman() bool {
	cmd := exec.Command("podman", "--version")
	err := cmd.Run()
	return err == nil
}

// Run tests
func Test() error {
	cmd := exec.Command("go", "test", "-v", "./tests")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Run linter
func Lint() error {
	cmd := exec.Command("docker", "run", "--rm", "--workdir=/data", "--volume=.:/data", "quay.io/helmpack/chart-testing:v3.7.1", "ct", "lint")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Update docs
func Docs() error {
	var cmd *exec.Cmd
	if isPodman() {
		cmd = exec.Command("podman", "run", "--rm", "--volume=.:/helm-docs:Z", "jnorwood/helm-docs:v1.14.2")
	} else {
		uid := fmt.Sprint(os.Getuid())
		cmd = exec.Command("docker", "run", "--rm", "--user", uid, "--volume=.:/helm-docs", "jnorwood/helm-docs:v1.14.2")
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
