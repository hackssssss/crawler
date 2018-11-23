package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

const StudyGolangurl = `https://www.qiushibaike.com/hot/`
const Qiushibaikeurl = `https://www.qiushibaike.com/hot/`
const xia55Url = `https://www.55xia.com/`

func TrimString(str string) string {
	return strings.TrimSpace(strings.Trim(strings.Trim(str, "\n"), " "))
}

type Crawl55xia struct {
	Title      string `sql:"column:title" json:"title,omitempty"`             //
	ImgHref    string `sql:"column:img_href" json:"img_href,omitempty"`       //
	DetailPage string `sql:"column:detail_page" json:"detail_page,omitempty"` //
}

func (Crawl55xia) TableName() string {
	return "crawl55xia"
}

func write2DB(title, href, filmDetail string) {
	if title == "" || href == "" || filmDetail == "" {
		return
	}
	data := Crawl55xia{
		Title:      title,
		ImgHref:    href,
		DetailPage: filmDetail,
	}
	err := db.Create(&data).Error
	if err != nil {
		fmt.Println(err)
		return
	}
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
	doc.Find(".row").Eq(0).Find(".index-box").Find(".movie-item").Each(func(i int, content *goquery.Selection) {
		title, _ := content.Find(".movie-item-in .lazy").Attr("title")
		picHref, _ := content.Find(".movie-item-in .lazy").Attr("data-src")
		filmDetail, _ := content.Find(".movie-item-in a").Attr("href")
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
		write2DB(title, picHref, filmDetail)
	})
	fmt.Println("Done")
	return nil
}

var db *gorm.DB

func init() {
	var err error
	db, err = gorm.Open("mysql", "root:1234567890@/crawler?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		fmt.Println("err : ", err)
		os.Exit(0)
	} else {
		fmt.Println("Open success,", db)
	}
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
