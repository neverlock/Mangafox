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
		"http://i.imgur.com/rbJvwK9.png",
		"http://i.imgur.com/4nAAGI1.png",
		"http://2.bp.blogspot.com/-KRWuh3dIi_U/UIbSuEXVGHI/AAAAAAACUdo/Ny8BDqRMJAU/s1600/KingsLogo.png",
		"http://i.imgur.com/Uinmy02.gif",
		"http://i.imgur.com/pg8ZSoA.jpg",
		"http://i.imgur.com/DGd795E.jpg",
		"http://2.bp.blogspot.com/-iY8hZJCMKzI/UEQnVdjdKcI/AAAAAAABeWg/EPDdxH71TU0/s64/1346643763_9.png",
		"http://upic.me/i/oh/30y18.png",
		}
	if len(os.Args) == 1 {
		fmt.Printf("Please use: %s manga_name maxchapter\n",os.Args[0])
		fmt.Printf("if you want to download http://www.kingsmanga.net/toriko-321/ from 1-312\n")
		fmt.Printf("Example: %s toriko 312\n",os.Args[0])
		fmt.Printf("if you want to download http://www.kingsmanga.net/toriko-321/ only 312\n")
		fmt.Printf("Example: %s toriko 312 0\n",os.Args[0])
		os.Exit(0)
	}

	if len(os.Args) == 3 {
		MAX,_ := strconv.Atoi(os.Args[2])

		for i:=1 ; i <= MAX ; i++ {
			urls := []string{}
			URL := fmt.Sprintf("http://www.kingsmanga.net/%s-%d",os.Args[1],i)
			x, _ := goquery.NewDocument(URL)
			x.Find("img").Each(func(idx int, s *goquery.Selection) {
			v, b := s.Attr("src")
			if b == true {
				urls = append(urls, v)
			}
			})

			for index,element := range urls {
				if !stringInSlice(element,ignore) {
					//fmt.Println("Index and element",index,element)
					downloadFromUrl(element,convertTo000(index),os.Args[1],convertTo000(i))
				}
			}
		}
		os.Exit(0)
	}

	if len(os.Args) == 4 {
		MAX,_ := strconv.Atoi(os.Args[2])
		END,_ := strconv.Atoi(os.Args[3])
		urls := []string{}
		if END == 0 {
			URL := fmt.Sprintf("http://www.kingsmanga.net/%s-%d",os.Args[1],MAX)
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
					downloadFromUrl(element,convertTo000(index),os.Args[1],convertTo000(MAX))
				}
			}
		} else {
			fmt.Println("Error")
		}
	}
}
