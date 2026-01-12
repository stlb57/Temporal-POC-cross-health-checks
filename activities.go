package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"time"
)

func CheckHealth(ctx context.Context, url string) (bool, error) {
	client := http.Client{Timeout: time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return false, nil
	}
	defer resp.Body.Close()
	return resp.StatusCode == 200, nil
}

func RestartTwin(ctx context.Context, selfID string, port2 int) error {
	var twin string
	var port1 int
	var peer int

	if selfID == "A" {
		twin = "B"
		port1 = port2
		peer = 8001
	} else {
		twin = "A"
		port1 = port2
		peer = 8002
	}

	args := []string{
		fmt.Sprintf("-twin=%s", twin),
		fmt.Sprintf("-port1=%d", port1),
		fmt.Sprintf("-port2=%d", peer),
	}

	cmd := exec.Command("./twin.exe", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Start()
}
