package main

import (
    "fmt"
    "github.com/PuerkitoBio/goquery"
)

 func stringInSlice(str string, list []string) bool {
	for _, v := range list {
		if v == str {
			return true
		}
	}
	return false
 }

func main() {
    ignore := []string{
		"http://i.imgur.com/rbJvwK9.png",
		"http://i.imgur.com/4nAAGI1.png",
		"http://2.bp.blogspot.com/-KRWuh3dIi_U/UIbSuEXVGHI/AAAAAAACUdo/Ny8BDqRMJAU/s1600/KingsLogo.png",
		"http://i.imgur.com/Uinmy02.gif",
		"http://i.imgur.com/pg8ZSoA.jpg",
		"http://i.imgur.com/DGd795E.jpg",
		"http://2.bp.blogspot.com/-iY8hZJCMKzI/UEQnVdjdKcI/AAAAAAABeWg/EPDdxH71TU0/s64/1346643763_9.png",
		}
    urls := []string{}
    x, _ := goquery.NewDocument("http://www.kingsmanga.net/toriko-321/")
    x.Find("img").Each(func(idx int, s *goquery.Selection) {
        v, b := s.Attr("src")
        if b == true {
            urls = append(urls, v)
        }
    })

	for index,element := range urls {
		if !stringInSlice(element,ignore) {
			fmt.Println("Index and element",index,element)
		}
	}
    //fmt.Println(urls)

}
