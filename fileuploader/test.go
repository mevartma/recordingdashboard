package main

import (
	"time"
	"fmt"
)

func main() {
	now := time.Now()
	newdate := fmt.Sprintf("%s",now.Format("2006-01-02"))
	fmt.Printf("Before: %v\r\nAfter: %s",now,newdate)
}
