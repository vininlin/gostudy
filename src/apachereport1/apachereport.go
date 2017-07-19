package main

import (
	"runtime"
	"os"
	"log"
	"bufio"
	"io"
	"safemap"
	"regexp"
	"fmt"
	"path/filepath"
)


var workers = runtime.NumCPU()

func main(){
	runtime.GOMAXPROCS(runtime.NumCPU())
	if len(os.Args) != 2 || os.Args[1] == "-h" || os.Args[1] == "--help" {
		fmt.Printf("usage: %s <file.log>\n", filepath.Base(os.Args[0]))
		os.Exit(1)
	}
	lines := make(chan string,workers*4)
	done := make(chan struct{},workers)
	pageMap := safemap.New()
	go readLines(os.Args[1],lines)
	processLines(done,pageMap,lines)
	waitUntil(done)
	showResults(pageMap)

}


func readLines(filename string,lines chan<- string)  {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal("failed to open the file:", err)
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	for {
		line,err := reader.ReadString('\n')
		if line != "" {
			lines <- line
		}
		if err != nil {
			if err != io.EOF {
				log.Println("failed to finish reading the file:", err)
			}
			break

		}
	}
	close(lines)
}

func processLines(done chan<- struct{},pageMap safemap.SafeMap,lines <-chan string)  {
	getRx := regexp.MustCompile(`GET[ \t]+([^ \t\n]+[.](html|js|css)?)`)
	incrementer := func(value interface{},found bool) interface{} {
			if found {
				return value.(int) + 1
			}
			return 1
	}
	for i := 0 ; i < workers; i++ {
		go func() {
			for line := range lines {
				if matches := getRx.FindStringSubmatch(line); matches != nil {
					fmt.Println("m0=",matches[0])
					fmt.Println("m1=",matches[1])
					pageMap.Update(matches[1],incrementer)
				}
			}
			done <- struct{}{}
		}()
	}
}

func waitUntil(done <-chan struct{})  {
	for i := 0 ; i < workers; i++ {
		<-done
	}
}

func showResults(pageMap safemap.SafeMap)  {
	pages := pageMap.Close()
	for page, count := range pages {
		fmt.Printf("%8d %s\n", count, page)
	}
}