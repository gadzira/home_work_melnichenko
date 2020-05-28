package main

import (
	"fmt"
	"log"
	"time"

	"github.com/beevik/ntp"
)

func main() {
	// Place your code here
	const (
		NtpServer = "0.beevik-ntp.pool.ntp.org"
		FatalfMsg = "Панека, беда, егор, отложенные вызовы не будут вызваны, буферы не отчистятся, потому что: %s"
	)

	et, err := ntp.Time(NtpServer)
	if err != nil {
		log.Fatalf(FatalfMsg, err)
	}
	ct := time.Now().UTC()
	fmt.Printf("current time: %s\n", ct.Format("2006-01-02 15:04:05 +0000 MST"))
	fmt.Printf("exact time: %s\n", et.Format("2006-01-02 15:04:05 +0000 MST"))
}
