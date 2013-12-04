package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/user"
	"regexp"
	"sync"

	"github.com/PuerkitoBio/goquery"
)

var (
	dir      string
	threadId = regexp.MustCompile(`thread-(\d*)-`)
	imageId  = regexp.MustCompile(`img/(.*)`)
)

func worker(linkChan chan string, wg *sync.WaitGroup) {
	defer wg.Done()

	for url := range linkChan {
		out, _ := os.Create(dir + "/" + imageId.FindStringSubmatch(url)[1])
		defer out.Close()

		resp, _ := http.Get(url)
		defer resp.Body.Close()

		io.Copy(out, resp.Body)
	}
}

func main() {
	url := "http://ck101.com/thread-2876990-1-1.html"

	doc, err := goquery.NewDocument(url)
	if err != nil {
		panic(err)
	}

	usr, _ := user.Current()
	title := doc.Find("h1#thread_subject").Text()
	dir = fmt.Sprintf("%v/Pictures/iloveck101/%v - %v", usr.HomeDir, threadId.FindStringSubmatch(url)[1], title)

	os.MkdirAll(dir, 0755)

	linkChan := make(chan string)
	wg := new(sync.WaitGroup)
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go worker(linkChan, wg)
	}

	doc.Find("div[itemprop=articleBody] img").Each(func(i int, img *goquery.Selection) {
		imgUrl, _ := img.Attr("file")
		linkChan <- imgUrl
	})

	close(linkChan)
	wg.Wait()
}
