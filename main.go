package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

const StudyGolangurl = `https://www.qiushibaike.com/hot/`
const Qiushibaikeurl = `https://www.qiushibaike.com/hot/`
const xia55Url = `https://www.55xia.com/`

func crawlStudyGolang() error {
	resp, err := http.Get(StudyGolangurl)
	if err != nil {
		log.Println("http get fail. err ", err)
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		log.Println("resp.StatusCode err, ", resp.StatusCode)
		if resp.StatusCode == 403 {
			log.Println("403 forbidden")
			time.Sleep(time.Millisecond * 10)
			os.Exit(1)
		}
		return err
	}
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal("NewDocument fail, ", err)
		return err
	}
	doc.Find("div.topic dd.right-info").Each(func(i int, content *goquery.Selection) {
		title := content.Find("div.title").Find("a").Text()
		strings.Trim(title, "\n")
		strings.Trim(title, " ")
		num := content.Find("div.meta").Find("div.num a").Text()
		strings.Trim(num, "\n")
		strings.Trim(num, " ")
		if len(num) == 0 {
			num = "0"
		}
		fmt.Println(title)
		fmt.Println(num)
	})
	return nil
}

//func test() {
//	doc, err := goquery.NewDocument(url)
//}

func TrimString(str string) string {
	return strings.TrimSpace(strings.Trim(strings.Trim(str, "\n"), " "))
}

func crawlQiushibaike() error {
	resp, err := http.Get(Qiushibaikeurl)
	if err != nil {
		log.Println("http get fail. err ", err)
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		log.Println("resp.StatusCode err, ", resp.StatusCode)
		if resp.StatusCode == 403 {
			log.Println("403 forbidden")
			time.Sleep(time.Millisecond * 10)
			os.Exit(1)
		}
		return err
	}
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal("NewDocument fail, ", err)
		return err
	}
	//doc.Find("div.content").Each(func (i int, content *goquery.Selection){
	//	interest := content.Find("")
	//	strings.Trim(title, "\n")
	//	strings.Trim(title, " ")
	//	num := content.Find("div.meta").Find("div.num a").Text()
	//	strings.Trim(num, "\n")
	//	strings.Trim(num, " ")
	//	if len(num) == 0 {
	//		num = "0"
	//	}
	//	fmt.Println(title)
	//	fmt.Println(num)
	//})
	content := doc.Find("div.content").Text()
	TrimString(content)
	fmt.Println(content)

	return nil
}

func crawl55xia() error {
	resp, err := http.Get(xia55Url)
	if err != nil {
		log.Println("http get fail. err ", err)
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		log.Println("resp.StatusCode err, ", resp.StatusCode)
		if resp.StatusCode == 403 {
			log.Println("403 forbidden")
			time.Sleep(time.Millisecond * 10)
			os.Exit(1)
		}
		return err
	}
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal("NewDocument fail, ", err)
		return err
	}
	doc.Find(".row").Eq(0).Find(".index-box").Find(".movie-item").Each(func (i int, content *goquery.Selection){
		title, _ := content.Find(".movie-item-in .lazy").Attr("title")
		picHref,_ := content.Find(".movie-item-in .lazy").Attr("data-src")
		filmDetail,_ := content.Find(".movie-item-in a").Attr("href")
		title = TrimString(title)
		picHref = TrimString(picHref)
		filmDetail = TrimString(filmDetail)
		if len(filmDetail) > 2 {
			filmDetail = filmDetail[2:]
			filmDetail = "https://" + filmDetail
		}
		fmt.Println(title)
		fmt.Println(picHref)
		fmt.Println(filmDetail)
		fmt.Println("==================")
	})
	fmt.Println("Done")
	return nil
}

func main() {
	count := 0
	for {
		if err := crawl55xia(); err != nil {
			time.Sleep(time.Second * 2)
			continue
		}
		count++
		if count >= 1 {
			break
		}
		time.Sleep(time.Second)
	}
}
