package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

const groupNameEnvVar = "AZ_GROUP"

func main() {
	groupName := os.Getenv(groupNameEnvVar)
	if groupName == "" {
		log.Fatalf(redF("[FATAL] Environment variable %s not set", groupNameEnvVar))
	}

	log.Printf("[INFO] Resetting Janitor for group %s", groupName)

	ts := time.Now().Format("2006-01-02T15:04:05Z")
	cmd := exec.Command(
		"az", "group", "update",
		"--name", groupName,
		"--tags",
		"CleanupFrequency=Weekly",
		"Created="+ts,
		"--only-show-errors",
		"--output", "none",
	)

	log.Printf("Running command: %v\n", cmd.Args)

	done := make(chan struct{})
	go tick(done)

	output, err := cmd.CombinedOutput()
	if err != nil {
		close(done)
		fmt.Println()
		trimmed := strings.Trim(string(output), " \n")
		log.Fatalf(redF("[FATAL] Received the following output (%v):\n\n%s\n", err, indented(trimmed, 4)))
	}

	close(done)
	fmt.Println()
	log.Println(green("[INFO] 😊  Finished resetting Janitor!"))
}

func tick(done chan struct{}) {
	ticker := time.NewTicker(80 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			fmt.Print(".")
		}
	}
}

func green(s string) string {
	return fmt.Sprintf("\033[32m%s\033[0m", s)
}

func greenF(s string, a ...interface{}) string {
	return green(fmt.Sprintf(s, a...))
}

func red(s string) string {
	return fmt.Sprintf("\033[31m%s\033[0m", s)
}

func redF(s string, a ...interface{}) string {
	return red(fmt.Sprintf(s, a...))
}

func indented(s string, n int) string {
	return strings.Repeat(" ", n) + strings.ReplaceAll(s, "\n", "\n"+strings.Repeat(" ", n))
}
