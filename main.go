package main

import (
	"log"
	"os"
	"os/exec"
	"syscall"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("need at least 1 argument for duration/deadline")
	}

	deadlineStr := os.Args[1]
	deadlineBytes := []byte(deadlineStr)
	var deadline time.Time

	err := deadline.UnmarshalText(deadlineBytes)
	if err != nil {
		deadlineDur, err2 := time.ParseDuration(deadlineStr)
		if err2 != nil {
			log.Fatal("parsing os.Args[1] as deadline or duration failed: %v and %v", err, err2)
		}
		deadline = time.Now().Add(deadlineDur)

		deadlineBytes, err = deadline.MarshalText()
		if err2 != nil {
			log.Fatal("failed to marshal deadline %v: %v", deadline, err)
		}
		deadlineStr = string(deadlineBytes)
	}

	for {
		now := time.Now()
		if now.Equal(deadline) || now.After(deadline) {
			break
		}

		cmd := exec.Command(os.Args[0], deadlineStr)
		cmd.SysProcAttr = &syscall.SysProcAttr{
			Setpgid: true,
		}
		err = cmd.Start()

	}
}
