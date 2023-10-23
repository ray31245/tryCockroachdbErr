package main

import (
	"fmt"

	"github.com/cockroachdb/errors"
)

func main() {
	baseErr := fmt.Errorf("an unexpected error occurred")

	wrappedErr1 := errors.Wrap(baseErr, "fail to process data")
	errors.WithDetail(wrappedErr1, "this also is detail")

	hintErr := errors.WithHintf(wrappedErr1, "Make sure your data is in the correct format.")

	wrappedErr2 := errors.Wrap(hintErr, "top-level operation failed")
	errors.WithDetail(wrappedErr2, "this is a detail")

	fmt.Println(wrappedErr2)

	hint := errors.GetAllHints(wrappedErr2)
	fmt.Println("Hint:", hint)
	allDetail := errors.GetAllDetails(wrappedErr2)
	fmt.Println(wrappedErr2.Error())
	fmt.Printf("allDetail: %+v\n", allDetail)
	flattenDetail := errors.GetAllSafeDetails(wrappedErr2)
	fmt.Println("flattenDetail:", flattenDetail)
	safDetail := errors.GetAllSafeDetails(wrappedErr2)
	fmt.Println("safDetail:", safDetail)
}
