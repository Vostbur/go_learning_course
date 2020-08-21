package main

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

func ExecutePipeline(jobs ...job) {
	Channels := make([]chan interface{}, len(jobs)+1)
	for i := range Channels {
		Channels[i] = make(chan interface{}, MaxInputDataLen)
	}

	wg := &sync.WaitGroup{}
	for i, j := range jobs {
		wg.Add(1)
		go func(i int, j job, in chan interface{}, out chan interface{}, wg *sync.WaitGroup) {
			defer wg.Done()
			j(in, out)
			if i == 0 {
				close(out)
			}
		}(i, j, Channels[i], Channels[i+1], wg)
		time.Sleep(time.Millisecond)
	}
	wg.Wait()
}

func SingleHash(in, out chan interface{}) {
	wg := &sync.WaitGroup{}
	mu := &sync.Mutex{}

	for i := range in {
		wg.Add(1)
		go func(d interface{}, wg *sync.WaitGroup) {
			defer wg.Done()
			data := fmt.Sprintf("%v", d)
			fmt.Println(d, "SingleHash data", data)

			mu.Lock()
			md5 := DataSignerMd5(data)
			mu.Unlock()
			fmt.Println(d, "SingleHash md5(data)", md5)

			wgIns := &sync.WaitGroup{}
			var crc32Md5, crc32 string
			wgIns.Add(2)
			go func(crc32Md5 *string, wgIns *sync.WaitGroup) {
				defer wgIns.Done()
				*crc32Md5 = DataSignerCrc32(md5)
			}(&crc32Md5, wgIns)
			go func(crc32 *string, wgIns *sync.WaitGroup) {
				defer wgIns.Done()
				*crc32 = DataSignerCrc32(data)
			}(&crc32, wgIns)
			wgIns.Wait()

			fmt.Println(d, "SingleHash crc32(md5(data))", crc32Md5)
			fmt.Println(d, "SingleHash crc32(data)", crc32)

			result := crc32 + "~" + crc32Md5
			out <- result
			fmt.Println(d, "SingleHash result", result)
		}(i, wg)
	}
	wg.Wait()
	close(out)
}

func MultiHash(in, out chan interface{}) {
	wg := &sync.WaitGroup{}
	mu := &sync.Mutex{}

	for d := range in {
		wg.Add(1)
		go func(d interface{}, wg *sync.WaitGroup) {
			defer wg.Done()
			data := fmt.Sprintf("%v", d)
			result := ""
			wgIns := &sync.WaitGroup{}
			for i := 0; i < 6; i++ {
				wgIns.Add(1)
				go func(r *string, i int, data string, wgIns *sync.WaitGroup) {
					defer wgIns.Done()
					th := strconv.Itoa(i)
					crc32 := DataSignerCrc32(th + data)
					mu.Lock()
					*r += crc32
					mu.Unlock()
					fmt.Println(d, "MultiHash: crc32(th+step1))", th, crc32)
				}(&result, i, data, wgIns)
				time.Sleep(100 * time.Millisecond)
			}
			wgIns.Wait()
			out <- result
			fmt.Println(d, "MultiHash: result", result)
			fmt.Println()
		}(d, wg)
	}
	wg.Wait()
	close(out)
}

func CombineResults(in, out chan interface{}) {
	var s []string
	for i := range in {
		data := fmt.Sprintf("%v", i)
		s = append(s, data)
	}
	sort.Strings(s)
	result := strings.Join(s, "_")
	fmt.Println(result)
	out <- result
	close(out)
}
