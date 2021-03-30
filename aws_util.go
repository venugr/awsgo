package main

import (
	"log"
	"time"
)

func DoSleep(sleepInSecs time.Duration, reason string) {

	log.Printf("Info: Sleep for %v secs...%v", sleepInSecs, reason)
	time.Sleep(sleepInSecs * time.Second)

}
