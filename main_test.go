package main

import (
	"log"
	"testing"
	"time"
)

func Test_Count(t *testing.T) {
	now := time.Now()

	time.Sleep(2 * time.Second)

	log.Printf("%.2fs", time.Since(now).Seconds())
}
