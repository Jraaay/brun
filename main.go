package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: brun <command> [args...]")
		fmt.Fprintln(os.Stderr, "env:   BARK_URL=https://api.day.app/<your-key>")
		os.Exit(2)
	}

	barkURL := strings.TrimRight(os.Getenv("BARK_URL"), "/")
	if barkURL == "" {
		fmt.Fprintln(os.Stderr, "brun: BARK_URL is not set")
		os.Exit(2)
	}

	cmdStr := strings.Join(os.Args[1:], " ")

	shell := os.Getenv("SHELL")
	if shell == "" {
		shell = "sh"
	}
	cmd := exec.Command(shell, "-ic", cmdStr)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(sigCh)

	start := time.Now()
	if err := cmd.Start(); err != nil {
		fmt.Fprintln(os.Stderr, "brun: failed to start:", err)
		pushBark(barkURL, "❌ 启动失败", fmt.Sprintf("%s\n\n%v", cmdStr, err))
		os.Exit(127)
	}

	go func() {
		for sig := range sigCh {
			if cmd.Process != nil {
				_ = cmd.Process.Signal(sig)
			}
		}
	}()

	err := cmd.Wait()
	duration := time.Since(start).Round(time.Millisecond)

	exitCode := 0
	switch e := err.(type) {
	case nil:
	case *exec.ExitError:
		exitCode = e.ExitCode()
	default:
		exitCode = 1
		fmt.Fprintln(os.Stderr, "brun: wait error:", err)
	}

	title := "命令运行成功"
	if exitCode != 0 {
		title = "命令运行失败"
	}
	body := fmt.Sprintf("%s\n耗时: %s", cmdStr, duration)
	if exitCode != 0 {
		body += fmt.Sprintf("\n退出码: %d", exitCode)
	}


	if err := pushBark(barkURL, title, body); err != nil {
		fmt.Fprintln(os.Stderr, "brun: bark push failed:", err)
	}

	os.Exit(exitCode)
}

func pushBark(url, title, body string) error {
	payload, _ := json.Marshal(map[string]any{
		"title": title,
		"body":  body,
	})
	req, err := http.NewRequest("POST", url, bytes.NewReader(payload))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return fmt.Errorf("bark returned status %s", resp.Status)
	}
	return nil
}
