package main

import (
	"runtime"
	"regexp"
	"os"
	"log"
	"bufio"
	"bytes"
	"io"
	"path/filepath"
	"fmt"
)

var workers = runtime.NumCPU()

type Result struct {
	filename string
	lino int
	line string
}

type Job struct {
	filename string
	results chan<- Result
}

func (job Job) Do(lineRx *regexp.Regexp){
	file, err := os.Open(job.filename)
	if err != nil {
		log.Printf("error:%s\n",err)
		return
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	for lino := 1; ; lino++ {
		line, err := reader.ReadBytes('\n')
		line = bytes.TrimRight(line,"\n\r")
		if lineRx.Match(line) {
			job.results <- Result{job.filename,lino,string(line)}
		}
		if err != nil {
			if err != io.EOF {
				log.Print("error: %s",err)
			}
			break
		}
	}
}

func main(){
	runtime.GOMAXPROCS(runtime.NumCPU())
	if len(os.Args) < 3 || os.Args[1] == "-h" || os.Args[1] == "--help" {
		fmt.Printf("usage: %s <regexp> <files>\n", filepath.Base(os.Args[0]))
		os.Exit(1)
	}
	if lineRx, err := regexp.Compile(os.Args[1]); err != nil {
		log.Fatalf("invalid regexp: %s\n", err)
	} else {
		grep(lineRx,commandLineFiles(os.Args[2:]))
	}
}

func grep(lineRx *regexp.Regexp,filenames []string)  {
	jobs := make(chan Job,workers)
	results := make(chan Result,mininum(1000,len(filenames)))
	done := make(chan struct{},workers)

	go addJobs(jobs,filenames,results)
	for i := 0; i < workers; i++ {
		go doJobs(done,lineRx,jobs)
	}
	waitAndProcessResults(done,results)
}

func commandLineFiles(files []string) []string  {
	if runtime.GOOS == "windows" {
		args := make([]string,0,len(files))
		for _, name := range files {
			if matches, err := filepath.Glob(name); err != nil {
				args = append(args,name)
			}else if matches != nil {
				args = append(args,matches...)
			}
		}
		return args
	}
	return files
}

func addJobs(jobs chan<- Job,filenames []string,results chan<- Result){
	for _, filename := range filenames {
		jobs <- Job{filename,results}
	}
	close(jobs)
}

func doJobs(done chan<- struct{}, lineRx *regexp.Regexp,jobs <-chan Job){
	for job := range jobs {
		job.Do(lineRx)
	}
	done <- struct{}{}
}

func waitAndProcessResults(done <-chan struct{},results <-chan Result){
	for working := workers; working > 0; {
		select {
		case result := <-results:
			fmt.Printf("%s:%d:%s\n", result.filename, result.lino, result.line)
		case <-done:
			working--
		}
	}
DONE:
	for {
		select {
		case result := <-results:
			fmt.Printf("%s:%d:%s\n", result.filename, result.lino, result.line)
		default :
			break DONE
		}
	}

}

func mininum(x int, ys ...int) int{
	for _, y := range ys {
		if y < x {
			x = y
		}
	}
	return x
}
