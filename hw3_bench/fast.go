package main

import (
	"bufio"
	"bytes"
	"fmt"
	"go_learning_course/hw3_bench/model"
	"io"
	"os"
	"strconv"
	"strings"
)

// > go test -bench . -benchmem
// goos: windows
// goarch: amd64
// pkg: go_learning_course/hw3_bench
// BenchmarkSlow-8               40          29963090 ns/op        19018695 B/op    195846 allocs/op
// BenchmarkFast-8              550           2148149 ns/op         1385620 B/op      8226 allocs/op

func FastSearch(out io.Writer) {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	seenBrowsers := make(map[string]bool, 1000)
	var foundUsers bytes.Buffer

	sc := bufio.NewScanner(file)
	sc.Split(bufio.ScanLines)

	i := -1
	for sc.Scan() {
		line := sc.Text()
		i++

		if !strings.Contains(line, "Android") && !strings.Contains(line, "MSIE") {
			continue
		}

		user := &model.User{}
		err := user.UnmarshalJSON([]byte(line))
		if err != nil {
			panic(err)
		}

		isAndroid := false
		isMSIE := false

		browsers := user.Browsers

		for _, browser := range browsers {

			if strings.Contains(browser, "Android") {
				isAndroid = true
				if _, exists := seenBrowsers[browser]; !exists {
					seenBrowsers[browser] = true
				}
			}

			if strings.Contains(browser, "MSIE") {
				isMSIE = true
				if _, exists := seenBrowsers[browser]; !exists {
					seenBrowsers[browser] = true
				}
			}
		}

		if !(isAndroid && isMSIE) {
			continue
		}

		email := user.Email
		index := strings.LastIndex(email, "@")
		email = email[:index] + " [at] " + email[index+1:]
		foundUsers.Write([]byte("[" + strconv.Itoa(i) + "] " + user.Name + " <" + email + ">\n"))
	}

	if err := sc.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading input:", err)
	}

	fmt.Fprintln(out, "found users:\n"+foundUsers.String())
	fmt.Fprintln(out, "Total unique browsers", len(seenBrowsers))
}
