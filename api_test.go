package main

import (
	"fmt"
	"log"
	"os/exec"
	"testing"
)

func TestAPIPort(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	// os.Exit(1)
	cmd := exec.Command("redis-server")
	err := cmd.Start()
	if err != nil {
		t.Fatalf("Can't run redis server, %v", err)
	}

	cmd.Process.Kill()
	cmd.Wait()
	log.Printf("Command finished with error: %v", err)
	fmt.Println(*redisServers)
}
