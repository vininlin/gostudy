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
	"time"
	"flag"
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
	log.SetFlags(0)
	var timeoutOpt *int64 = flag.Int64("timeout",0,"seconds (0 means no timeout)")
	flag.Parse()
	if *timeoutOpt < 0 || *timeoutOpt > 240 {
		log.Fatalln("timeout must be in the range [0,240] seconds")
	}
	args := flag.Args()
	if len(args) < 1 {
		log.Fatalln("a regexp to match must be specified")
	}
	pattern := args[0]
	files := args[1:]
	if len(files) < 1 {
		log.Fatalln("must provide at least one filename")
	}
	if lineRx, err := regexp.Compile(pattern); err != nil {
		log.Fatalf("invalid regexp: %s\n", err)
	} else {
		var timeout int64 = 1e9 * 60 * 10
		if *timeoutOpt != 0 {
			timeout = *timeoutOpt * 1e9
		}
		grep(timeout,lineRx,commandLineFiles(files))
	}
}

func grep(timeout int64,lineRx *regexp.Regexp,filenames []string)  {
	jobs := make(chan Job,workers)
	results := make(chan Result,mininum(1000,len(filenames)))
	done := make(chan struct{},workers)

	go addJobs(jobs,filenames,results)
	for i := 0; i < workers; i++ {
		go doJobs(done,lineRx,jobs)
	}
	waitAndProcessResults(timeout,done,results)
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

func waitAndProcessResults(timedout int64,done <-chan struct{},results <-chan Result){
	//返回一个超时channel
	finish := time.After(time.Duration(timedout))
	for working := workers; working > 0; {
		select {
		case result := <-results :
			fmt.Printf("##%s:%d:%s\n", result.filename, result.lino, result.line)
		case <-finish :
			fmt.Println("timed out")
			return
		case <-done :
			fmt.Println("done")
			working--
		}
	}
	for {
		select {
		case result := <-results :
			fmt.Printf("**%s:%d:%s\n", result.filename, result.lino, result.line)
		case <-finish :
			fmt.Println("*timed out")
			return
		default :
			fmt.Println("default")
			return
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
