package main

import (
	"fmt"
	"net/http"
	"time"
	"io"
	"os"
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

func main() {
	// https://music.163.com/song?id=1357785909

	// download_url := "http://music.163.com/song/media/outer/url?id=436514312.mp3"
	downloadURL := "http://music.163.com/song/media/outer/url?id="
	//mainURL := "https://music.163.com/album?id=36457807"
	//mainURL := "https://music.163.com/artist?id=5771"
	mainURL := "https://music.163.com/playlist?id=368529707"
	// mainURL := "https://music.163.com/playlist?id=2319166657"
	c := colly.NewCollector(
	// Visit only root url and urls which start with "e" or "h" on httpbin.org
	// colly.URLFilters(
	// 	regexp.MustCompile("http://ps.youxiake\\.com/.*"),
	// ),
	// colly.AllowedDomains("music.163.com"),
	)
// #song-list-pre-cache
	c.OnHTML("#song-list-pre-cache .f-hide", func(e *colly.HTMLElement) {
		e.ForEach("a[href]", func(index int, item *colly.HTMLElement){
			name := item.Text
			href := item.Attr("href")
			id := href[9:]
			fmt.Println("name:", name, "    id:", id)
			durl := downloadURL+id
			filename := name + ".mp3"
			fmt.Println("begin download mp3:", name)
			file, err := os.Create(filename)
			defer file.Close()
			client := http.DefaultClient
			client.Timeout = time.Second*60
			resp, err := client.Get(durl)
			if resp.ContentLength <= 0 {
				fmt.Println("length is error!")
			}
			raw := resp.Body
			defer raw.Close()
			_,err = io.Copy(file, raw)
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

		// Only those links are visited which are matched by  any of the URLFilter regexps
		// c.Visit(e.Request.AbsoluteURL(link+"/"))
	})

	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
		// r.Headers.Add("Referrer Policy", "no-referrer-when-downgrade")
		// r.Headers.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/71.0.3578.98 Safari/537.36")
	})

	c.OnResponse(func(r *colly.Response) {
		// fmt.Println("response body: ", string(r.Body[:]))
	})

	// Start scraping on http://httpbin.org
	// c.Async = true
	c.Visit(mainURL)
	c.Wait()
}
