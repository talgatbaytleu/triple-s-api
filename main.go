package main

import (
	"fmt"
	"regexp"
)

type struct_1 struct {
	fieldOne int
	fieldTwo string
}

func main() {
	// triples.Run()

	re := regexp.MustCompile("([[:alpha:]]-[[:upper:]]){4}")
	string1 := "quicK brown fox jumped over the lazy dog"
	testBool := re.FindStringSubmatch(string1)
	fmt.Println(testBool)
}
