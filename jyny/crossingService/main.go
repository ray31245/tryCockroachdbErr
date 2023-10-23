package main

import (
	"context"
	"fmt"
	"log"

	"github.com/cockroachdb/errors"
	"github.com/cockroachdb/errors/errorspb"
)

// CustomError is a custom error type
type CustomError struct {
	Message string
}

func (e *CustomError) Error() string {
	return e.Message
}

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)

	// Create an instance of the custom error
	originalErr := &CustomError{Message: "This is a custom error"}

	// Encode the error
	encodedErr := errors.EncodeError(context.Background(), originalErr)

	encodedByte, err := encodedErr.Marshal()
	if err != nil {
		log.Println(err)
	}
	log.Println(string(encodedByte))

	errorspb := errorspb.EncodedError{}
	errorspb.Unmarshal(encodedByte)

	// For demonstration purposes, let's decode the error back
	decodedErr := errors.DecodeError(context.Background(), errorspb)

	// Check if the decoded error is of type CustomError using errors.Is()
	if errors.Is(decodedErr, originalErr) {
		fmt.Printf("Is decoded error a CustomError? %v\n", decodedErr)
	}
}
