package main

/*
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
*/

import (
	"bytes"
	"io/ioutil"
	"testing"
)

// запускаем перед основными функциями по разу чтобы файл остался в памяти в файловом кеше
// ioutil.Discard - это ioutil.Writer который никуда не пишет
func init() {
	SlowSearch(ioutil.Discard)
	FastSearch(ioutil.Discard)
}

// -----
// go test -v

func TestSearch(t *testing.T) {
	slowOut := new(bytes.Buffer)
	SlowSearch(slowOut)
	slowResult := slowOut.String()

	fastOut := new(bytes.Buffer)
	FastSearch(fastOut)
	fastResult := fastOut.String()

	if slowResult != fastResult {
		t.Errorf("results not match\nGot:\n%v\nExpected:\n%v", fastResult, slowResult)
	}
}

// -----
// go test -bench . -benchmem

func BenchmarkSlow(b *testing.B) {
	for i := 0; i < b.N; i++ {
		SlowSearch(ioutil.Discard)
	}
}

func BenchmarkFast(b *testing.B) {
	for i := 0; i < b.N; i++ {
		FastSearch(ioutil.Discard)
	}
}
