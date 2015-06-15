package main

import (
    "os"
    "io"
    "fmt"
    "log"
    "os/exec"
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

func convertTo000(index int) string{
	if index < 10 {
		return fmt.Sprintf("00%d",index)
	}
	if index < 99 {
		return fmt.Sprintf("0%d",index)
	}
	if index < 999 {
		return fmt.Sprintf("%d",index)
	}
	return "000"
}

func downloadFromUrl(url string,index string,mangaName string,path string) {
/*
	tokens := strings.Split(url, "/")
	newName := tokens[len(tokens)-1]
	fmt.Println("Downloading", url, "to", newName)
*/
	tokens := strings.Split(url, "/")
	fileNameExt := tokens[len(tokens)-1]
	fileName := strings.Split(fileNameExt,".")
	ext := fileName[len(fileName)-1]
//	fmt.Println("Downloading", url, "to", fileNameExt)

	mkdirOpt := fmt.Sprintf("%s/%s",mangaName,path)
        cmd := exec.Command("mkdir","-p",mkdirOpt)

        err := cmd.Run()
        if err != nil {
                log.Fatal(err)
        }

	newName := fmt.Sprintf("%s/%s/%s.%s",mangaName,path,index,ext)
//	fmt.Println("Downloading", url, "to", newName)
	fmt.Printf("Downloading..[%s] [%s] [%s.%s] ",mangaName,path,index,ext)

	// TODO: check file existence first with io.IsExist
	output, err := os.Create(newName)
	if err != nil {
		fmt.Println("Error while creating", newName, "-", err)
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
		}

	if len(os.Args) == 1 {
		fmt.Printf("Please use: %s manga_name maxchapter\n",os.Args[0])
		fmt.Printf("if you want to download all series of http://www.niceoppai.net/Saikin-Kono-Sekai-wa-Watashi-dake-no-Mono-ni-Narimashita\n")
		fmt.Printf("Example: %s Saikin-Kono-Sekai-wa-Watashi-dake-no-Mono-ni-Narimashita\n",os.Args[0])
		fmt.Printf("if you want to download http://www.niceoppai.net/Saikin-Kono-Sekai-wa-Watashi-dake-no-Mono-ni-Narimashita/22/ only 22\n")
		fmt.Printf("Example: %s Saikin-Kono-Sekai-wa-Watashi-dake-no-Mono-ni-Narimashita 22 0\n",os.Args[0])
		os.Exit(0)
	}

	if len(os.Args) == 2 {
		page := []string{}
		URL := fmt.Sprintf("http://www.niceoppai.net/%s",os.Args[1])
		x,_ := goquery.NewDocument(URL)
		x.Find("a").Each(func(idx int,s *goquery.Selection) {
		v, b := s.Attr("href")
		if b == true && strings.Contains(v,os.Args[1]) {
			page = append(page,v)
//			fmt.Println("link =",v)
		}
		})
		MAX := len(page)-1
		//fmt.Println("start link=",page[MAX])
		//have 2 junk link 
		for i:=MAX; i>=2 ;i-- {
			fmt.Printf("Link = %s\n",page[i])
			urls := []string{}
			//URL := page[i]
			URL := fmt.Sprintf("%s?all",page[i])
			tokens := strings.Split(URL, "/")
			chapter := tokens[len(tokens)-2]

			x, _ := goquery.NewDocument(URL)
			x.Find("img").Each(func(idx int, s *goquery.Selection) {
			v, b := s.Attr("src")
			//fmt.Println("link =",v)
			if b == true {
				urls = append(urls, v)
			}
			})
			fmt.Printf("### Get data from [%s] volume [%s] ###\n",os.Args[1],chapter)
			for index,element := range urls {
				if !stringInSlice(element,ignore) {
					//fmt.Println("Index and element",index,element)
					downloadFromUrl(element,convertTo000(index),os.Args[1],chapter)
				}
			}
		}
		os.Exit(0)
	}

	if len(os.Args) == 4 {
		//MAX,_ := strconv.Atoi(os.Args[2])
		END,_ := strconv.Atoi(os.Args[3])
		urls := []string{}
		if END == 0 {
			fmt.Printf("[Get data]http://www.niceoppai.net/%s/%s/?all\n",os.Args[1],os.Args[2])
			URL := fmt.Sprintf("http://www.niceoppai.net/%s/%s/?all",os.Args[1],os.Args[2])
			x, _ := goquery.NewDocument(URL)
			x.Find("img").Each(func(idx int, s *goquery.Selection) {
			v, b := s.Attr("src")
			if b == true {
				urls = append(urls, v)
			}
			})
			for index,element := range urls {
				if !stringInSlice(element,ignore) {
					//fmt.Printf("Index[%d]:=%s\n",index,element)
					downloadFromUrl(element,convertTo000(index),os.Args[1],os.Args[2])
				}
			}
		} else {
			fmt.Println("Error")
		}
	}
}
