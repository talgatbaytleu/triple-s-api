package main

import triples "triple-s/cmd/triple-s"

type struct_1 struct {
	fieldOne int
	fieldTwo string
}

func main() {
	triples.Run()
	// var obj struct_1
	//
	// obj.fieldOne = 1
	// fmt.Println(obj.fieldOne)
	// fmt.Printf("Type : %T\nValue : %+v\n", obj, obj)
	// fmt.Println(obj)
}
