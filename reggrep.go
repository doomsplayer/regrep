package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var dir = flag.String("path", os.Getenv(`PWD`), "path/to/grep")
var pattern = flag.String("pattern", ".*", "regexp match pattern")
var verbose = flag.Bool("verbose", false, "")
var hide = flag.Bool("hide", false, "")
var suffix = flag.String("suffix", "", "")
var recursion = flag.Bool("recursion", false, "")
var outputnumber = 1

func main() {
	flag.Parse()
	if *verbose {
		fmt.Println(`pattern is :`, *pattern)
	}
	do(*dir)
}
func e(err error) {
	if err != nil {
		panic(err)
		os.Exit(1)
	}

}

func do(in string) {
	dir, err := ioutil.ReadDir(in)
	e(err)
	reg, err := regexp.Compile(*pattern)
	if err != nil {
		fmt.Println(`regexp illegal`)
	}

	for _, v := range dir {
		if !*hide {
			if v.Name()[0] == '.' {
				continue
			}
		}

		if *verbose {
			fmt.Println(`open: `, in+"/"+v.Name(), `   ---------isDIR: `, v.IsDir())
		}
		if !v.IsDir() {
			if func() bool {
				suf := strings.Split(*suffix, "/")
				retbool := false

				for _, u := range suf {
					if strings.HasSuffix(v.Name(), u) {
						return true
					}
				}
				return retbool
			}() {

				f, err := os.OpenFile(in+"/"+v.Name(), os.O_RDONLY, 0666)
				e(err)
				bufr := bufio.NewReader(f)
				iline := 1
				for line, _, err := bufr.ReadLine(); err == nil; func() { line, _, err = bufr.ReadLine(); iline++ }() {
					for _, u := range reg.FindAllString(string(line), -1) {
						relpath, _ := filepath.Rel(os.Getenv(`PWD`), filepath.Clean(in+"/"+v.Name()))
						fmt.Println("\033[31m", outputnumber, "\033[49;37m", `at`, "\033[44;37m", relpath, "\033[49;37m", "\033[32m", `line`, iline, "\033[37;40m", `:`, "\033[30;47m", u, "\033[37;40m", "\n")
						outputnumber++
					}
				}
				f.Close()
			}
		} else {
			if *recursion {
				do(in + "/" + v.Name())
			}

		}
	}
}
