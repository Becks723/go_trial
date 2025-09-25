package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	for page := 250; page <= 420; page++ {
		fmt.Printf("正在爬取第 %d 页的数据……\n", page)
		crawl(page)
	}
}

/* 文章数据 */
type article struct {
	Date   string   // 日期
	Author string   // 作者
	Title  string   // 标题
	Body   []string // 正文
}

func crawl(page int) {
	// 1. 发送请求
	pageStr := strconv.Itoa(page)
	req, err := http.NewRequest("GET", "https://info22.fzu.edu.cn/lm_list.jsp?totalpage=1088&PAGENUM="+pageStr+"&wbtreeid=1460", nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/140.0.0.0 Safari/537.36")

	client := http.Client{}
	response, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	// 解析网页
	docDetail, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	var articles []article

	docDetail.Find("body > div.sy-content > div > div.right.fr > div.list.fl > ul > li").
		Each(func(i int, s *goquery.Selection) {
			author := s.Find("p > a.lm_a").Text()
			a2 := s.Find("p > a:nth-child(2)")
			href, _ := a2.Attr("href") // 正文链接，给的是相对地址，需要进一步处理
			title, _ := a2.Attr("title")
			date := s.Find("p > span").Text()

			ymd := strings.Split(date, "-")
			year, _ := strconv.Atoi(ymd[0])
			month, _ := strconv.Atoi(ymd[1])
			day, _ := strconv.Atoi(ymd[2])
			// 只爬取 2020-01-01 ~ 2021-09-01 的文章
			if year == 2020 || (year == 2021 && month < 9) || (year == 2021 && month == 9 && day == 1) {
				norm := normalizeHref("https://info22.fzu.edu.cn", href)

				// 继续请求正文
				client2 := http.Client{}
				req, err := http.NewRequest("GET", norm, nil)
				if err != nil {
					log.Fatal(err)
				}
				req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/140.0.0.0 Safari/537.36")
				response, err := client2.Do(req)
				if err != nil {
					log.Fatal(err)
				}
				defer response.Body.Close()

				docDetail2, err := goquery.NewDocumentFromReader(response.Body)
				if err != nil {
					log.Fatal(err)
				}

				// #vsb_content > div > p:nth-child(1)
				// #vsb_content > div > p:nth-child(2)
				// #vsb_content > div > p:nth-child(21)
				// 正文body分成了若干段，用切片储存起来
				var body []string
				docDetail2.Find("#vsb_content > div > p").
					Each(func(i int, s *goquery.Selection) {
						body = append(body, s.Find("span").Text())
					})

				data := article{date, author, title, body}
				articles = append(articles, data)
			}
		})

	for _, v := range articles {
		fmt.Println(v)
	}
}

/* 标准化正文地址 */
func normalizeHref(base string, href string) string {
	href = strings.TrimSpace(href)
	b, _ := url.Parse(base)
	u, _ := url.Parse(href)
	abs := b.ResolveReference(u) // 拼接根链接和相对地址
	return abs.String()
}
