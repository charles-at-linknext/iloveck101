package main

import (
	"bufio"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"net/http"
	"os"
	"os/user"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
	"github.com/coreos/go-log/log"
	"github.com/spf13/cobra"
)

var (
	dir      string
	threadId = regexp.MustCompile(`thread-(\d*)-`)
	imageId  = regexp.MustCompile(`([^\/]+)\.(png|jpg)`)
)

func worker(linkChan chan string, wg *sync.WaitGroup) {
	defer wg.Done()

	for target := range linkChan {
		resp, err := http.Get(target)
		if err != nil {
			log.Println(err)
			continue
		}
		defer resp.Body.Close()

		m, _, err := image.Decode(resp.Body)
		if err != nil {
			log.Println(err)
			continue
		}

		// Ignore small images
		bounds := m.Bounds()
		if bounds.Size().X > 300 && bounds.Size().Y > 300 {
			imgInfo := imageId.FindStringSubmatch(target)
			out, _ := os.Create(dir + "/" + imgInfo[1] + "." + imgInfo[2])
			defer out.Close()
			switch imgInfo[2] {
			case "jpg":
				jpeg.Encode(out, m, nil)
			case "png":
				png.Encode(out, m)
			}
		}
	}
}

func crawler(target string, workerNum int) {
	doc, err := goquery.NewDocument(target)
	if err != nil {
		panic(err)
	}

	usr, _ := user.Current()
	title := doc.Find("h1#thread_subject").Text()
	dir = fmt.Sprintf("%v/Pictures/iloveck101/%v - %v", usr.HomeDir, threadId.FindStringSubmatch(target)[1], title)

	os.MkdirAll(dir, 0755)

	linkChan := make(chan string)
	wg := new(sync.WaitGroup)
	for i := 0; i < workerNum; i++ {
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

// [todo] - Holy shit function, should refactor it!
func printGoogleResult(keyword string, page int) (hrefs []string) {
	client := &http.Client{}
	queryUrl := fmt.Sprintf("https://www.google.com.tw/search?espv=210&es_sm=119&q=%v+site:ck101.com&start=%v", keyword, page*10)
	req, err := http.NewRequest("GET", queryUrl, nil)
	if err != nil {
		panic(err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_9_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/31.0.1650.57 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	doc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		panic(err)
	}

	hrefs = make([]string, 0)

	// Print result list
	doc.Find("li.g h3.r a").Each(func(i int, s *goquery.Selection) {
		title := s.Text()
		href, exist := s.Attr("href")
		if exist {
			hrefs = append(hrefs, href)
			fmt.Printf("[%v] %v\n", i, title)
		}
	})

	// Print pages
	for i := page - 3; i <= page+2; i++ {
		if i >= 0 {
			if i == page {
				fmt.Printf("[%v] ", i)
			} else {
				fmt.Printf("%v ", i)
			}
		}
	}
	fmt.Println()

	return hrefs
}

func main() {

	var postUrl string
	var workerNum int

	rootCmd := &cobra.Command{
		Use:   "iloveck101",
		Short: "Download all the images in given post url",
		Run: func(cmd *cobra.Command, args []string) {
			crawler(postUrl, workerNum)
		},
	}
	rootCmd.Flags().StringVarP(&postUrl, "url", "u", "http://ck101.com/thread-2876990-1-1.html", "Url of post")
	rootCmd.Flags().IntVarP(&workerNum, "worker", "w", 10, "Number of workers")

	searchCmd := &cobra.Command{
		Use:   "search",
		Short: "Download all the images in given post url",
		Run: func(cmd *cobra.Command, args []string) {
			page := 0
			keyword := args[0]
			hrefs := printGoogleResult(keyword, page)

			scanner := bufio.NewScanner(os.Stdin)
			quit := false

			for !quit {
				fmt.Print("ck101> ")

				if !scanner.Scan() {
					break
				}

				line := scanner.Text()
				parts := strings.Split(line, " ")
				cmd := parts[0]
				args := parts[1:]

				switch cmd {
				case "quit":
					quit = true
				case "n":
					page = page + 1
					hrefs = printGoogleResult(keyword, page)
				case "p":
					if page > 0 {
						page = page - 1
					}
					hrefs = printGoogleResult(keyword, page)
				case "d":
					index, err := strconv.ParseUint(args[0], 0, 0)
					if err != nil {
						fmt.Println(err)
						continue
					}
					if int(index) >= len(hrefs) {
						fmt.Println("Invalid index")
						continue
					}

					// Only support url with format ck101.com/thread-xxx
					if threadId.Match([]byte(hrefs[index])) {
						crawler(hrefs[index], 10)
						fmt.Println("Done!")
					} else {
						fmt.Println("Unsupport url", hrefs[index])
					}
				default:
					fmt.Println("unrecognized command:", cmd, args)
				}
			}
		},
	}

	rootCmd.AddCommand(searchCmd)
	rootCmd.Execute()
}
