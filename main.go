package main

import (
	"bufio"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"sync"
	"time"
)

var green = "\033[92m"
var red = "\033[91m"
var reset = "\033[0m"

const ASCII_ART = `
_._     _,-'""` + "`" + `-._
(,-.` + "`" + `._,'(       |\` + "-/" + `|
    ` + "`" + `-.-' \ )-` + "`" + `( , o o)
          ` + "`" + `-    \` + "_`" + `"'` + "`" + `-'

Credit : https://github.com/softwaretobi/go-checkerproxie/
`

func CheckProxy(proxy string, timeout time.Duration, outputFile string, wg *sync.WaitGroup, mutex *sync.Mutex, verifiedCount *int, nonVerifiedCount *int) {
	defer wg.Done()

	proxyURL, err := url.Parse("http://" + proxy)
	if err != nil {
		fmt.Println("Invalid proxy URL:", proxy)
		return
	}

	proxyFunc := http.ProxyURL(proxyURL)
	client := &http.Client{
		Transport: &http.Transport{Proxy: proxyFunc},
		Timeout:   timeout,
	}

	_, err = client.Get("http://google.com")
	if err == nil {
		mutex.Lock()
		*verifiedCount++
		f, err := os.OpenFile(outputFile, os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Println("Error writing to file:", err)
			mutex.Unlock()
			return
		}
		defer f.Close()
		_, _ = f.WriteString(proxy + "\n")
		mutex.Unlock()
		fmt.Println(green + "Valid proxy : " + proxy + reset)
	} else {
		mutex.Lock()
		*nonVerifiedCount++
		mutex.Unlock()
		fmt.Println(red + "Invalid proxy : " + proxy + reset)
	}
}

func main() {
	if len(os.Args) != 5 {
		PrintHelp()
		os.Exit(1)
	}

	proxyFile := os.Args[1]

	outputFile := os.Args[3]
	timeoutMs := os.Args[4]

	timeout, err := time.ParseDuration(timeoutMs + "ms")
	if err != nil {
		fmt.Println("Invalid timeout value:", err)
		return
	}

	file, err := os.Open(proxyFile)
	if err != nil {
		fmt.Println("Error reading proxy file:", err)
		return
	}
	defer file.Close()

	_ = os.Remove(outputFile)
	outFile, err := os.Create(outputFile)
	if err != nil {
		fmt.Println("Error creating output file:", err)
		return
	}
	defer outFile.Close()

	var proxies []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		proxies = append(proxies, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading proxies:", err)
		return
	}

	totalProxies := len(proxies)
	var wg sync.WaitGroup
	var mutex sync.Mutex
	verifiedCount := 0
	nonVerifiedCount := 0

	for _, proxy := range proxies {
		wg.Add(1)
		go CheckProxy(proxy, timeout, outputFile, &wg, &mutex, &verifiedCount, &nonVerifiedCount)
	}

	wg.Wait()

	fmt.Printf("\nThe test of the proxies is finished. Valid proxies have been saved in the %s\n", outputFile)
	fmt.Printf("Final results:\n")
	fmt.Printf("Verified proxies   : %d / %d\n", verifiedCount, totalProxies)
	fmt.Printf("Unverified proxies : %d / %d\n", nonVerifiedCount, totalProxies)
	fmt.Println(ASCII_ART)
}

func PrintHelp() {
	fmt.Println("Usage: go run proxy_checker.go <proxy_file> <proxy_type> <output_file> <timeout>")
	fmt.Println("Arguments:")
	fmt.Println("  <proxy_file> : Path of the file containing the proxies to be tested")
	fmt.Println("  <proxy_type> : Type of proxy to search for.")
	fmt.Println("  <output_file> : Path of the output file to save valid proxies")
	fmt.Println("  <timeout>    : Timeout in milliseconds for each test request")
	fmt.Println()
	fmt.Println("Example:")
	fmt.Println("  go run proxy_checker.go proxies.txt http valids.txt 5000")
}
