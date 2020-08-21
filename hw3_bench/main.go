package main

import "fmt"

func main() {
	fmt.Print(`
	go test -v - чтобы проверить что ничего не сломалось
	go test -bench . -benchmem - для просмотра производительности

	go test -bench .
	go test -bench . -benchmem

	go test -bench . -benchmem -cpuprofile="cpu.out" -memprofile="mem.out" -memprofilerate=1

	Linux:
	go tool pprof hw3_bench.test cpu.out
	go tool pprof hw3_bench.test mem.out

	(pprof) top
	(pprof) list FastSearch
	(pprof) web
	(pprof) quit
`)
}
