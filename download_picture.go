package main

import (
	"fmt"
	"github.com/gocolly/colly"
	// "github.com/afly3/colly"
	// "io/ioutil"
	// "net/http"
	// "github.com/PuerkitoBio/goquery"
	// "github.com/levigross/grequests"
	// "os"
	// "sync"
	// "strings"
	// "time"
	// "strconv"
)


func main() {
	
	// main_url := "http://ps.youxiake.com/gallery/"
	main_url := "http://ps.youxiake.com/album/556144/"
	c := colly.NewCollector(
		// Visit only root url and urls which start with "e" or "h" on httpbin.org
		// colly.URLFilters(
		// 	regexp.MustCompile("http://ps.youxiake\\.com/.*"),
		// ),
		colly.AllowedDomains("ps.youxiake.com"),
	)
	filter_c := colly.NewCollector(
		// Visit only root url and urls which start with "e" or "h" on httpbin.org
		// colly.URLFilters(
		// 	regexp.MustCompile("http://ps.youxiake\\.com/.*"),
		// ),
		colly.AllowedDomains("ps.youxiake.com"),
	)
	// On every a element which has href attribute call callback
	c.OnHTML("a[href^=\"/album\"]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		// Print link
		fmt.Printf("Link found: %s\n", e.Request.AbsoluteURL(link))
		// Only those links are visited which are matched by  any of the URLFilter regexps
		c.Visit(e.Request.AbsoluteURL(link+"/"))
	})

	c.OnHTML(".icon-bg.like", func(e *colly.HTMLElement){
		text := e.Text
			fmt.Printf("album text: %s %#v\n", text, e.Request.URL.RequestURI())

		// filter_c.Visit(e.Request.URL.RequestURI())
	})

	filter_c.OnHTML(".img-con", func(e *colly.HTMLElement){
		src := e.ChildAttr("img", "src")
		fmt.Printf("image url: %s\n", src)

	})

	// Before making a request print "Visiting ..."
	// c.OnRequest(func(r *colly.Request) {
		// fmt.Println("Visiting", r.URL.String())
	// })


	// Start scraping on http://httpbin.org
	// c.Async = true
	c.Visit(main_url)
	// c.Wait()
}

// func main2(){
// 		doc, err := goquery.NewDocument("http://ps.youxiake.com/gallery/")
// 		if err != nil {
// 			fmt.Errorf("download error:%#v", err)
// 			os.Exit(-1)
// 		}
// 		doc.Find("div a").Each(func(i int, s *goquery.Selection){
// 			// src, exists := s.Find("a").Attr("href")
// 			src, exists := s.Attr("target")
// 			if exists && src == "_blank" {
// 				href,_ := s.Attr("href")
// 			   fmt.Println("===", href)
// 			}
// 		})
//     return
// }

// var wg sync.WaitGroup
// func main3(){
// 	now := time.Now()
// 	initalUrls := []string{"http://www.zngirls.com/girl/18071/album/",}
// 	for _, url := range initalUrls {
// 		doc, err := goquery.NewDocument(url)
// 		if err != nil {
// 			fmt.Errorf("download error:%#v", err)
// 			os.Exit(-1)
// 		}
// 		doc.Find(".gli_link").Each(func(i int, s *goquery.Selection){
// 			src, exists := s.Find("img").Attr("src")
// 			atl, _ := s.Find("img").Attr("alt")
// 			fmt.Printf("begin download image album %v src %v\n", atl, src)

// 			if(exists){
// 				wg.Add(1)
// 				go func ( src string) {
// 					defer wg.Done()

// 					n :=0 
// 					s := strings.Replace(src, "cover/", "", 1)
// 					ss := strings.Split(s, "/")
// 					fm := strings.Join(ss[:len(ss)-1], "/")
// 					sf0 := fm + "/%d.jpg"
// 					sfn := fm + "/%03d.jpg"

// 					for{
// 						s := ""
// 						if n==0 {
// 							s = fmt.Sprintf(sf0, n)
// 						}else{
// 							s = fmt.Sprintf(sfn, n)
// 						}

// 						fmt.Printf("begin download: %v\n", s)
// 						res, _ := grequests.Get(s, &grequests.RequestOptions{
// 							Headers:map[string]string {
// 								"Referer":"http://www.zngirls.com",
// 								"User-Agent":"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36"}})
// 						if res.StatusCode != 200 {
// 							fmt.Printf("download is error, exit the download for ablum:%s\n", src)	
// 						}
// 						length := res.Header.Get("Content-Length")
// 						slen,_ := strconv.Atoi(length)
// 						if slen < 4100 {
// 							fmt.Printf("download is error, exit the download for ablum:%s\n", src)	
// 							break
// 						}

// 						index := strings.Index(s, "gallery")
// 						if index == -1 {
// 							fmt.Errorf("invalid http addr, can't find the key word of gallery, parse error:%s\n", src)
// 						}
// 						ss2 := strings.Split(string(s[index:]), "/")
// 						dirname := strings.Join(ss2[:len(ss2)-1], "/")
// 						if _, err := os.Stat(dirname); err != nil {
// 							fmt.Printf("create the download dir:%s\n", dirname)
// 							os.MkdirAll(dirname, 0744)
// 						}
// 						filename := strings.Join(ss2, "/")
// 						res.DownloadToFile(filename)
// 						fmt.Printf("successfully download the image to dir: %s\n", filename)
// 						n++
// 					}
// 				}(src)
// 			}
// 		})
// 	}

// 	wg.Wait()
// 	fmt.Printf("download task complete, time consuming: %#v\n", time.Now().Sub(now))

// }