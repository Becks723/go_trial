package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type replyData struct {
	Data struct {
		Cursor struct {
			Pagination_reply struct {
				Next_offset string
			}
		}
		Replies []reply
	}
}

type reply struct {
	Content struct {
		Message string
	}
	Member struct {
		Uname string
	}
}

func main() {
	offset := ""
	for i := 1; i <= 5; i++ {
		fmt.Printf("正在爬取第 %d 页的评论\n", i)

		offset = crawl(offset)
		time.Sleep(2 * time.Second)
	}
}

func crawl(offset string) (nextOffset string) {
	unix := time.Now().Unix()
	nowtime := strconv.FormatInt(unix, 10)
	w_rid, pagination_str := getSign(offset, nowtime)

	baseUrl := "https://api.bilibili.com/x/v2/reply/wbi/main"

	params := url.Values{}
	params.Set("oid", "420981979")
	params.Set("type", "1")
	params.Set("mode", "3")
	params.Set("pagination_str", pagination_str)
	params.Set("plat", "1")
	if offset == "" {
		params.Set("seek_rpid", "")
	}
	params.Set("web_location", "1315875")
	params.Set("w_rid", w_rid)
	params.Set("wts", nowtime)

	fullUrl := fmt.Sprintf("%s?%s", baseUrl, params.Encode())

	req, err := http.NewRequest("GET", fullUrl, nil)
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

	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s", body)

	var rd replyData
	if err := json.Unmarshal(body, &rd); err != nil {
		log.Fatal(err)
	}

	for i, r := range rd.Data.Replies {
		fmt.Printf("%d: user: %s, comment: %s\n", i, r.Member.Uname, r.Content.Message)
	}

	nextOffset = rd.Data.Cursor.Pagination_reply.Next_offset
	return
}

func getSign(offset string, nowtime string) (w_rid string, pagination_str string) {
	pagination_str = fmt.Sprintf(`{"offset":"%s"}`, offset)

	params := url.Values{}
	params.Set("oid", "420981979")
	params.Set("type", "1")
	params.Set("mode", "3")
	params.Set("pagination_str", pagination_str)
	params.Set("plat", "1")
	if offset == "" {
		params.Set("seek_rpid", "")
	}
	params.Set("web_location", "1315875")
	params.Set("wts", nowtime)
	v := params.Encode()
	a := "ea1db124af3c7062474693fa704f4ff8"
	w_rid = getMD5Hash(v + a)
	return
}

func getMD5Hash(s string) string {
	hash := md5.Sum([]byte(s))
	return hex.EncodeToString(hash[:])
}
