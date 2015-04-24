package main

import (
    "os"
    "io"
    "fmt"
    "strings"
    "strconv"
    "net/http"
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

func downloadFromUrl(url string) {
	tokens := strings.Split(url, "/")
	fileName := tokens[len(tokens)-1]
	fmt.Println("Downloading", url, "to", fileName)

	// TODO: check file existence first with io.IsExist
	output, err := os.Create(fileName)
	if err != nil {
		fmt.Println("Error while creating", fileName, "-", err)
		return
	}
	defer output.Close()

	response, err := http.Get(url)
	if err != nil {
		fmt.Println("Error while downloading", url, "-", err)
		return
	}
	defer response.Body.Close()

	n, err := io.Copy(output, response.Body)
	if err != nil {
		fmt.Println("Error while downloading", url, "-", err)
		return
	}

	fmt.Println(n, "bytes downloaded.")
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
	if len(os.Args) == 1 {
		fmt.Printf("Please use: %s manga_name maxchapter\n",os.Args[0])
		fmt.Printf("if you want to download http://www.kingsmanga.net/toriko-321/ from 1-312\n")
		fmt.Printf("Example: %s toriko 312\n",os.Args[0])
		os.Exit(0)
	}

	if len(os.Args) == 3 {
		fmt.Println(os.Args[2])
		volume,_ := strconv.Atoi(os.Args[2])
		URL := fmt.Sprintf("http://www.kingsmanga.net/%s-%d",os.Args[1],volume)
		fmt.Println(URL)
		x, _ := goquery.NewDocument(URL)
		x.Find("img").Each(func(idx int, s *goquery.Selection) {
		v, b := s.Attr("src")
		if b == true {
			urls = append(urls, v)
		}
		})

		for index,element := range urls {
			if !stringInSlice(element,ignore) {
				fmt.Println("Index and element",index,element)
				downloadFromUrl(element)
			}
		}
	}
}
