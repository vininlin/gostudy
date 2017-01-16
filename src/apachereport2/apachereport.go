package main

import (
	"runtime"
	"sync"
	"os"
	"log"
	"path/filepath"
	"bufio"
	"io"
	"regexp"
	"fmt"
)

var workers = runtime.NumCPU()

type pageMap struct {
	countForPage map[string]int
	mutex *sync.RWMutex
}

func main(){
	runtime.GOMAXPROCS(runtime.NumCPU())
	if len(os.Args) != 2 || os.Args[1] == "-h" || os.Args[1] == "--help" {
		fmt.Printf("usage: %s <file.log>\n", filepath.Base(os.Args[0]))
		os.Exit(1)
	}
	lines := make(chan string,workers * 4)
	done := make(chan struct{},workers)
	pageMap := NewPageMap()
	go readLines(os.Args[1],lines)
	getRx := regexp.MustCompile(`GET[ \t]+([^ \t\n]+[.]html?)`)
	for i := 0; i < workers; i++ {
		go processLines(done,getRx,pageMap,lines)
	}
	waitUtil(done)
	showResult(pageMap)

}

func NewPageMap() *pageMap {
	return &pageMap{make(map[string]int),new (sync.RWMutex)}
}

func (pm *pageMap) Increatement(page string){
	pm.mutex.Lock()
	defer pm.mutex.Unlock()
	pm.countForPage[page]++
}

func (pm *pageMap) Len() int{
	pm.mutex.RLock()
	defer pm.mutex.RUnlock()
	return len(pm.countForPage)
}

func readLines(filename string,lines chan <- string) {
	file,err := os.Open(filename)
	if err != nil {
		log.Fatal("failed to open the file:", err)
	}
	defer  file.Close()
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

func processLines(done chan<- struct{},getRx *regexp.Regexp,pageMap *pageMap,lines <-chan string){
	for line := range lines {
		if matches := getRx.FindStringSubmatch(line); matches != nil {
			pageMap.Increatement(matches[1])
		}
	}
	done <- struct {}{}
}

func waitUtil(done <-chan struct{})  {
	for i := 0; i < workers; i++ {
		<-done
	}
}

func showResult(pageMap *pageMap){
	for page,count := range pageMap.countForPage {
		fmt.Printf("%8d %s\n", count, page)
	}
}

