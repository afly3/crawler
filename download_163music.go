package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/gocolly/colly"
	"github.com/tidwall/gjson"
	// "net/http"
	// "github.com/PuerkitoBio/goquery"
	// "github.com/levigross/grequests"
	// "sync"
	// "strings"
	// "time"
	// "strconv"
)

// 实际中应该用更好的变量名
var (
	dir string
	url string
)

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func main() {
	downloadURL := "http://music.163.com/song/media/outer/url?id="

	flag.StringVar(&dir, "dir", ".", "the dir download")
	flag.StringVar(&url, "url", "", "the url of 163 music")
	flag.Parse()
	if 7 >= len(url) {
		flag.Usage()
		return
	}

	// 另一种绑定方式
	// q = flag.Bool("q", false, "suppress non-error messages during configuration testing")

	mainURL := url
	// mainURL := "https://music.163.com/playlist?id=1988668862"
	if exsit, err := PathExists(dir); !exsit {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			fmt.Println("dir can't create, %+v", err)
			return
		}
	}
	c := colly.NewCollector()
	// #song-list-pre-cache
	c.OnHTML("#song-list-pre-cache .f-hide", func(e *colly.HTMLElement) {
		e.ForEach("a[href]", func(index int, item *colly.HTMLElement) {
			name := item.Text
			href := item.Attr("href")
			id := href[9:]
			fmt.Println("name:", name, "    id:", id)
			durl := downloadURL + id
			filename := dir + "/" + name + ".mp3"
			fmt.Println("begin download mp3:", name)
			file, err := os.Create(filename)
			defer file.Close()
			client := http.DefaultClient
			client.Timeout = time.Second * 60
			resp, err := client.Get(durl)
			if resp.ContentLength <= 0 {
				fmt.Println("length is error!")
			}
			raw := resp.Body
			defer raw.Close()
			_, err = io.Copy(file, raw)
			if err != nil {
				fmt.Println(err)
			}
		})
		// array := gjson.Parse(songJson).Array()
		// for _, item := range array {
		// 	name := item.Get("name").String()
		// 	id := item.Get("id").String()

		// }
		// Only those links are visited which are matched by  any of the URLFilter regexps
		// c.Visit(e.Request.AbsoluteURL(link+"/"))
	})

	// On every a element which has href attribute call callback
	c.OnHTML("#song-list-pre-data", func(e *colly.HTMLElement) {
		return
		songJson := e.Text
		array := gjson.Parse(songJson).Array()
		for _, item := range array {
			name := item.Get("name").String()
			id := item.Get("id").String()
			durl := downloadURL + id
			filename := name + ".mp3"
			fmt.Println("begin download mp3:", name)
			file, err := os.Create(filename)
			defer file.Close()
			client := http.DefaultClient
			client.Timeout = time.Second * 60
			resp, err := client.Get(durl)
			if resp.ContentLength <= 0 {
				fmt.Println("length is error!")
			}
			raw := resp.Body
			defer raw.Close()
			_, err = io.Copy(file, raw)
			if err != nil {
				fmt.Println(err)
			}

		}
	})

	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	// Start scraping on http://httpbin.org
	// c.Async = true
	c.Visit(mainURL)
	c.Wait()
}
