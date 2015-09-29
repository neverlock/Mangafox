package main

import (
    "os"
    "io"
    "fmt"
    "log"
    "os/exec"
    "strings"
    //"strconv"
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

func downloadFromUrl(url string,index string,mangaName string,volume string,chapter string) {
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

	mkdirOpt := fmt.Sprintf("%s/%s/%s",mangaName,volume,chapter)
        cmd := exec.Command("mkdir","-p",mkdirOpt)

        err := cmd.Run()
        if err != nil {
                log.Fatal(err)
        }

	newName := fmt.Sprintf("%s/%s/%s/%s.%s",mangaName,volume,chapter,index,ext)
//	fmt.Println("Downloading", url, "to", newName)
	fmt.Printf("Downloading..[%s] [%s] [%s] [%s.%s] ",mangaName,volume,chapter,index,ext)

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
/*
	ignore := []string{
		"http://i.imgur.com/rbJvwK9.png",
		"http://i.imgur.com/4nAAGI1.png",
		"http://2.bp.blogspot.com/-KRWuh3dIi_U/UIbSuEXVGHI/AAAAAAACUdo/Ny8BDqRMJAU/s1600/KingsLogo.png",
		"http://i.imgur.com/Uinmy02.gif",
		"http://i.imgur.com/pg8ZSoA.jpg",
		"http://i.imgur.com/DGd795E.jpg",
		"http://2.bp.blogspot.com/-iY8hZJCMKzI/UEQnVdjdKcI/AAAAAAABeWg/EPDdxH71TU0/s64/1346643763_9.png",
		"http://upic.me/i/oh/30y18.png",
		}
*/
	if len(os.Args) == 1 {
		fmt.Printf("Please use: %s manga_name maxchapter\n",os.Args[0])
		fmt.Printf("if you want to download all series of http://mangafox.me/manga/sui_youbi/\n")
		fmt.Printf("Example: %s sui_youbi\n",os.Args[0])
		fmt.Printf("if you want to download sui_youbi/v01/*\n")
		fmt.Printf("Example: %s sui_youbi v01 0\n",os.Args[0])
		os.Exit(0)
	}

	if len(os.Args) == 2 {
		page := []string{}
		URL := fmt.Sprintf("http://mangafox.me/manga/%s",os.Args[1])
		x,_ := goquery.NewDocument(URL)
		x.Find("a").Each(func(idx int,s *goquery.Selection) {
		v, b := s.Attr("href")
		if b == true && strings.Contains(v,os.Args[1]) {
			page = append(page,v)
		//	fmt.Println("link =",v)
		}
		})
		/***** Find all link volume , chapter from first page ******/
		MAX := len(page)-1
		//fmt.Println("start link=",page[MAX])
		for i:=MAX; i>=4 ;i-- {
			//urls := []string{}
			URL := page[i]
			tokens := strings.Split(URL, "/")
			chapter := tokens[len(tokens)-2]
			volume := tokens[len(tokens)-3]


			//fmt.Println("Volume =",volume," Chapter =",chapter)
			/********* find max page from maaga/volume/chapter *******/
			fmt.Printf("### Get Page from [%s] volume [%s] chapter [%s] ###\n",os.Args[1],volume,chapter)
			page := []string{}
			y,_ := goquery.NewDocument(URL)
			y.Find("div[id*='top_center_bar'] select option").Each(func(idx int,s *goquery.Selection) {
			v,b := s.Attr("value")
			if b == true {
				page = append(page,v)
				//fmt.Println("page =",v)
			}
			})
			fmt.Printf("### [%s] / [%s] / [%s] Max page = %s\n",os.Args[1],volume,chapter,page[len(page)-2])
			for j:= 0 ; j<= len(page)-2; j++ {
				newURL := fmt.Sprintf("http://mangafox.me/manga/%s/%s/%s/%s.html",os.Args[1],volume,chapter,page[j])
				z,_:= goquery.NewDocument(newURL)
				z.Find("div[id*='viewer'] img").Each(func(idx int,s *goquery.Selection){
				v,b := s.Attr("src")
				if b == true {
					//fmt.Println("Image=",v)
					downloadFromUrl(v,convertTo000(j),os.Args[1],volume,chapter)
				}
				})
			}
		}
		os.Exit(0)
	}

/*
	if len(os.Args) == 4 {
		//MAX,_ := strconv.Atoi(os.Args[2])
		END,_ := strconv.Atoi(os.Args[3])
		urls := []string{}
		if END == 0 {
			URL := fmt.Sprintf("http://www.kingsmanga.net/%s-%s",os.Args[1],os.Args[2])
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
*/
}
