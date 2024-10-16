package main

import triples "triple-s/cmd/triple-s"

func main() {
	triples.Run()
	// _, err := os.Stat("cmd")
	// if os.IsNotExist(err) {
	// 	fmt.Println("not exist")
	// } else if os.IsExist(err) {
	// 	fmt.Println("exist")
	// }
}
