package main

import (
	"bufio"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
)

var url string
var wordlist string
var output string
var threads int
var verbose bool
var headers string
var urls []string

var wg sync.WaitGroup

func main() {
	flag.StringVar(&url, "u", "", "The target url")
	flag.StringVar(&wordlist, "w", "", "The wordlist path")
	flag.IntVar(&threads, "t", 5, "The threads number")
	flag.BoolVar(&verbose, "v", false, "Enable verbose mode")
	flag.StringVar(&headers, "H", "", "Set a custom header")
	flag.StringVar(&output, "o", "", "The output file")
	flag.Parse()

	wg.Add(threads)
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

	if headers == "" {
		for i := 0; i < threads; i++ {
			go getStatusCode(s)
		}
	} else {
		for i := 0; i < threads; i++ {
			go getStatusCodeHaders(s)
		}
	}
	wg.Wait()

	if output != "" {
		crr, _ := os.Create(output)
		crr.Close()

		out, err := os.OpenFile(output, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			panic(err)
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
			//fmt.Printf("%v/%v - %v\n", url, path, req.Status)
			urls = append(urls, fmt.Sprintf("%v/%v\n", url, path))
		}
		if verbose {
			fmt.Printf("%v/%v - %v\n", url, path, req.Status)
		}
	}
}

func getStatusCodeHaders(file *bufio.Scanner) {
	defer wg.Done()
	for file.Scan() {
		path := fmt.Sprint(strings.ReplaceAll(fmt.Sprint(file.Text()), " ", ""))
		client := &http.Client{}
		req, err := http.NewRequest("GET", fmt.Sprintf("%v/%v", url, path), nil)
		if err != nil {
			panic(err)
		}
		req.Header.Add(strings.Split(headers, ":")[0], strings.Split(headers, ":")[1])
		resp, _ := client.Do(req)
		if resp.StatusCode == 200 {
			fmt.Printf("%v/%v - %v\n", url, path, resp.Status)
			urls = append(urls, fmt.Sprintf("%v/%v\n", url, path))
		}
		if verbose {
			fmt.Printf("%v/%v - %v\n", url, path, resp.Status)
		}
	}
}
