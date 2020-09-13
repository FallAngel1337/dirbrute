package main

import (
	"bufio"
	"flag"
	"fmt"
	"net/http"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"
)

var url string
var wordlist string
var verbose bool
var output string
var urls []string

var wg sync.WaitGroup

func main() {
	flag.StringVar(&url, "u", "", "The target url")
	flag.StringVar(&wordlist, "w", "", "The wordlist path")
	flag.BoolVar(&verbose, "v", false, "Enable verbose mode")
	flag.StringVar(&output, "o", "", "The output file")
	flag.Parse()

	cpus := runtime.NumCPU()
	wg.Add(cpus)

	f, err := os.Open(wordlist)
	defer f.Close()

	if err != nil {
		panic("Could not read file:" + wordlist)
	}

	s := bufio.NewScanner(f)

	testRequest, _ := http.Get(url + "/8u3xrj981u2s4r98u2s31j9du498rdh8sj8j9k2su9kia39uej1")
	if testRequest.StatusCode == 200 {
		fmt.Println("It returs a 200, maybe its not good to perform a fuzzing! :(")
		os.Exit(1)
	}

	if verbose {
		start := time.Now()
		fmt.Print("[+] Starting...\n\n")
		for i := 0; i < cpus; i++ {
			go getStatusCode(s)
		}
		stop := time.Since(start)
		time.Sleep(time.Second)
		fmt.Print("\n[!] Time spent: ", stop)
	} else {
		for i := 0; i < cpus; i++ {
			go getStatusCode(s)
		}
	}
	wg.Wait()

	if output != "" {
		crr, _ := os.Create(output)
		crr.Close()

		out, err := os.OpenFile(output, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			panic("WOW!")
		}
		defer out.Close()

		for _, u := range urls {
			out.WriteString(u)
		}
	}
}

func getStatusCode(file *bufio.Scanner) {
	defer wg.Done()
	for file.Scan() {
		path := fmt.Sprint(strings.ReplaceAll(fmt.Sprint(file.Text()), " ", ""))
		req, err := http.Get(fmt.Sprintf("%v/%v", url, path))

		if err != nil {
			panic("Some error occurred!")
			os.Exit(1)
		}

		if req.StatusCode == 200 {
			fmt.Printf("%v/%v - %v\n", url, path, req.Status)
			urls = append(urls, fmt.Sprintf("%v/%v\n", url, path))
		}
	}
}
