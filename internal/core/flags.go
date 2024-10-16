package core

import (
	"flag"
	"fmt"
	"os"
)

func InitFlags() (string, string) {
	var dirPath string
	var port string
	flag.StringVar(
		&dirPath,
		"dir",
		"",
		"dirPath to the directory where the files will be stored as arguments",
	)

	flag.StringVar(&port, "port", "8080", "port number API gonna listening to")
	flag.Parse()

	if dirPath == "" {
		fmt.Fprintf(
			os.Stderr,
			"Please determine path by -dir flag\n Use -help flag for more information\n",
		)
		os.Exit(1)
	}

	dirPath += "/"

	return dirPath, port
}
