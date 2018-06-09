package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"io/ioutil"
	"strings"
)

func main() {
	verbose := flag.Bool("v", false, "Print files verbosely")
	force := flag.Bool("f", false, "Apply changes without asking")

	flag.Parse()

	if len(flag.Args()) < 3 {
		flag.PrintDefaults()
		os.Exit(1)
	}

	var dirs []string
	needle := flag.Arg(0)
	replacement := flag.Arg(1)

	if len(flag.Args()) >= 2 {
		dirs = flag.Args()[2:]
	} else {
		dirs = append(dirs, ".")
	}

	if !*force {
		fmt.Println("No changes will be applied unless -f is given.")
	}

	replaced := 0

	for _, dir := range dirs {
		if *verbose {
			fmt.Println("Scanning", dir)
		}

		err := filepath.Walk(dir, func (path string, info os.FileInfo, err error) error {
			if err != nil {
				fmt.Errorf("E Error scanning dir %s. Skipping.\n", dir)
				return nil
			}

			if !info.IsDir() {
				contents, err := ioutil.ReadFile(path)

				if err != nil {
					fmt.Errorf("E Error reading %s. Skipping.\n", path)
				} else {
					oldContents := string(contents)
					newContents := strings.Replace(oldContents, needle, replacement, -1)

					if oldContents != newContents {
						replaced++
						fmt.Println("+", path)

						if *force {
							ioutil.WriteFile(path, []byte(newContents), info.Mode())
						}
					} else {
						if *verbose {
							fmt.Println("-", path)
						}
					}
				}
			}

			return nil
		})

		if err != nil {
			panic(err)
		}
	}

	if !*force {
		fmt.Printf("%d file(s) WOULD have been updated.\n", replaced)
	} else {
		fmt.Printf("%d file(s) were updated.\n", replaced)
	}
}
