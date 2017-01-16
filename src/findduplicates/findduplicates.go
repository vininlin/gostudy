package main

import (
	"sync"
	"os"
	"runtime"
	"path/filepath"
	"log"
	"crypto/sha1"
	"io"
	"fmt"
	"sort"
)

const maxSizeOfSmallFile  = 1024 * 30
const maxGoroutines  = 100

type pathsInfo struct {
	size int64
	paths []string
}

type fileInfo struct {
	sha1 []byte
	size int64
	path string
}

func main()  {
	runtime.GOMAXPROCS(runtime.NumCPU())
	if len(os.Args) == 1 || os.Args[1] == "-h" || os.Args[1] == "--help" {
		fmt.Printf("usage: %s <path>\n", filepath.Base(os.Args[0]))
		os.Exit(1)
	}
	infoChan := make(chan fileInfo,maxGoroutines * 2)
	go findDuplications(infoChan,os.Args[1])
	pathData := mergeResults(infoChan)
	outputResults(pathData)


}

func findDuplications(infoChan chan fileInfo,dirname string){
	waiter := &sync.WaitGroup{}
	filepath.Walk(dirname,makeWalkFunc(infoChan,waiter))
	waiter.Wait()
	close(infoChan)
}

func makeWalkFunc(infoChan chan fileInfo,waiter *sync.WaitGroup) func(string,os.FileInfo,error)error{
	return func(path string,info os.FileInfo,err error) error {
		if err == nil && info.Size() > 0 && (info.Mode() & os.ModeType == 0) {
			if info.Size() < maxSizeOfSmallFile || runtime.NumGoroutine() > maxGoroutines {
				processFile(path,info,infoChan,nil)
			}else{
				waiter.Add(1)
				go processFile(path,info,infoChan, func() {waiter.Done()})
			}
		}
		return nil
	}
}

func processFile(filename string,info os.FileInfo,infoChan chan fileInfo,done func())  {
	if done != nil {
		defer done()
	}
	file,err := os.Open(filename)
	if err != nil {
		log.Println("error:",err)
	}
	defer file.Close()
	hash := sha1.New()
	if size,err := io.Copy(hash,file); size != info.Size() || err != nil {
		if err != nil {
			log.Println("error:",err)
		} else {
			log.Println("error: failed to read the whole file:", filename)
		}
		return
	}
	infoChan <- fileInfo{hash.Sum(nil),info.Size(),filename}
}

func mergeResults(infoChan <-chan fileInfo)map[string]*pathsInfo  {
	pathData := make(map[string]*pathsInfo)
	format := fmt.Sprintf("%%016X:%%%dX", sha1.Size*2)
	fmt.Println(format)
	for info := range infoChan {
		key := fmt.Sprint(format,info.size,info.sha1)
		value,found := pathData[key]
		if !found {
			value = &pathsInfo{size: info.size}
			pathData[key] = value
		}
		value.paths = append(value.paths,info.path)
	}
	return pathData
}

func outputResults(pathData map[string]*pathsInfo)  {
	keys := make([]string,0,len(pathData))
	for key := range pathData {
		keys = append(keys,key)
	}
	sort.Strings(keys)
	for _,key := range keys {
		value := pathData[key]
		if len(value.paths) > 1 {
			fmt.Printf("%d duplicate files (%s bytes):\n",
				len(value.paths), commas(value.size))
			sort.Strings(value.paths)
			for _,name := range value.paths {
				fmt.Printf("\t%s\n", name)
			}
		}
	}
}

func commas(x int64) string {
	value := fmt.Sprint(x)
	for i := len(value) - 3; i > 0 ; i -= 3 {
		value = value[:i] + "," + value[i:]
	}
	return value
}


