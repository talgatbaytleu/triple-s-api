package core

import (
	"flag"
	"fmt"
	"os"
)

func InitFlags() (string, string) {
	var help bool
	var dirPath string
	var port string
	flag.BoolVar(&help, "help", false, "triple-s usage information")
	flag.StringVar(
		&dirPath,
		"dir",
		"",
		"dirPath to the directory where the files will be stored as arguments",
	)

	flag.StringVar(&port, "port", "8080", "port number API gonna listening to")
	flag.Parse()

	if help {
		fmt.Println(`Simple Storage Service.

**Usage:**
    triple-s [-port <N>] [-dir <S>]  
    triple-s --help

**Options:**
- --help     Show this screen.
- --port N   Port number
- --dir S    Path to the directory`)
		os.Exit(0)
	}
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
