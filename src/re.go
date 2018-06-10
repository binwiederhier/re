package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"io/ioutil"
	"strings"
	"regexp"
)

var needle string
var replacement string
var dirs []string
var replaced int
var excludes []regexp.Regexp
var includes []regexp.Regexp
var flagApply *bool
var flagExcludes *string
var flagIncludes *string

func walk(path string, info os.FileInfo, err error) error {
	// Skip if in exclude list
	basename := filepath.Base(path)

	for _, exclude := range excludes {
		if exclude.Match([]byte(basename)) {
			if info.IsDir() {
				fmt.Println("Skipping excluded directory", path)
				return filepath.SkipDir
			} else {
				return nil
			}
		}
	}


	// Skip if not matching the includes list
	if len(includes) > 0 {
		matches := false
		
		for _, include := range includes {
			if include.Match([]byte(basename)) {
				matches = true
				break
			}
		}
		
		if !matches {
			return nil
		}
	}

	// Skip folders
	if info.IsDir() {
		return nil
	}

	// Read contents
	contents, err := ioutil.ReadFile(path)

	if err != nil {
		fmt.Errorf("E Error reading %s. Skipping.\n", path)
		return nil
	}

	// Compare contents
	oldContents := string(contents)
	newContents := strings.Replace(oldContents, needle, replacement, -1)

	if oldContents != newContents {
		replaced++
		fmt.Println("+", path)

		if *flagApply {
			ioutil.WriteFile(path, []byte(newContents), info.Mode())
		}
	}

	return nil
}

func parseRegexList(s string) []regexp.Regexp {
	stringList := strings.Split(s, ",")
	result := make([]regexp.Regexp, 0)

	for i := 0; i < len(stringList); i++ {
		if stringList[i] == "" {
			continue
		}

		pattern := regexp.QuoteMeta(stringList[i])
		pattern = strings.Replace(pattern, "\\*", ".*", -1)
		pattern = strings.Replace(pattern, "\\?", ".", -1)
		pattern = "^" + pattern + "$"

		regex, err := regexp.Compile(pattern)

		if err != nil {
			panic(err)
		}

		result = append(result, *regex)
	}

	return result
}

func main() {
	flagApply = flag.Bool("f", false, "Apply changes")
	flagExcludes = flag.String("e", ".bzr,CVS,.git,.hg,.svn", "Comma-separated list of excluded files, wildcards supported")
	flagIncludes = flag.String("i", "", "Comma-separated list of included files, wildcards supported, e.g \"*.js,*.html,*index.*\"")

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

	if !*flagApply {
		fmt.Println("No changes will be applied unless -f is given.")
	}

	excludes = parseRegexList(*flagExcludes)
	includes = parseRegexList(*flagIncludes)

	for _, dir := range dirs {
		err := filepath.Walk(dir, walk)

		if err != nil {
			fmt.Errorf("Cannot scan directory: %s; Skipping.\n", err.Error())
		}
	}

	if !*flagApply {
		fmt.Printf("%d file(s) WOULD have been updated.\n", replaced)
	} else {
		fmt.Printf("%d file(s) were updated.\n", replaced)
	}
}
