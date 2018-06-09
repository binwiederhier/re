package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"io/ioutil"
	"strings"
)

var needle string
var replacement string
var dirs []string
var verbose *bool
var apply *bool
var replaced int

func walk(path string, info os.FileInfo, err error) error {
	if info.IsDir() {
		return nil
	}

	contents, err := ioutil.ReadFile(path)

	if err != nil {
		fmt.Errorf("E Error reading %s. Skipping.\n", path)
		return nil
	}

	oldContents := string(contents)
	newContents := strings.Replace(oldContents, needle, replacement, -1)

	if oldContents != newContents {
		replaced++
		fmt.Println("+", path)

		if *apply {
			ioutil.WriteFile(path, []byte(newContents), info.Mode())
		}
	} else {
		if *verbose {
			fmt.Println("-", path)
		}
	}

	return nil
}

func main() {
	verbose = flag.Bool("v", false, "Verbose output")
	apply = flag.Bool("f", false, "Apply changes")

	flag.Usage = func() {
		fmt.Println("Syntax: re [options] SEARCH REPLACEMENT [DIR ...]")
		fmt.Println()
		fmt.Println("Options:")
		flag.PrintDefaults()
	}

	flag.Parse()

	if len(flag.Args()) < 2 {
		flag.Usage()
		os.Exit(1)
	}

	needle = flag.Arg(0)
	replacement = flag.Arg(1)

	if len(flag.Args()) > 2 {
		dirs = flag.Args()[2:]
	} else {
		dirs = append(dirs, ".")
	}

	if !*apply {
		fmt.Println("No changes will be applied unless -f is given.")
	}

	for _, dir := range dirs {
		if *verbose {
			fmt.Println("Scanning", dir)
		}

		err := filepath.Walk(dir, walk)

		if err != nil {
			fmt.Errorf("Cannot scan directory: %s; Skipping.\n", err.Error())
		}
	}

	if !*apply {
		fmt.Printf("%d file(s) WOULD have been updated.\n", replaced)
	} else {
		fmt.Printf("%d file(s) were updated.\n", replaced)
	}
}
