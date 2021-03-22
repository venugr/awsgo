package main

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
)

func main() {
	fmt.Println("\nWelcome AWS SDK V2 Go...awsgo v0-Beta-0.0.0\n\n")

	myContext := context.TODO()
	myConfig, lErr := config.LoadDefaultConfig(myContext)

	if lErr != nil {
		log.Fatalf("Error: config error, %v", lErr)
	}

	// DoUser(myConfig, myContext)
	DoPolicy(myConfig, myContext)
	fmt.Println()
}
