package main

import "os"

func main() {
	os.Exit(0) // want "call of os.Exit in func main of package main"
}
