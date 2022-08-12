package scrapper

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type extractedJob struct {
	id string
	title string
	review string
	salary string
	grade string
}


// term1 = industry_id >> service == 100
// term2 = city_id >> busan == 4
func Scrape(term1 , term2 string) {
	var baseURL string = "https://www.jobplanet.co.kr/companies?industry_id=" +term1 + "&city_id=" + term2
	var jobs []extractedJob
	c := make(chan []extractedJob)
	totalPages := getPages(baseURL)

	for i := 0; i< totalPages; i++ {
		go getPage(i,baseURL, c)
	}

	for i:=0; i<totalPages; i++ {
		extractedJob:= <-c
		jobs = append(jobs, extractedJob...)
	}

	writeJobs(jobs)
	fmt.Println("Done, extracted", len(jobs))
}

func writeJobs(jobs []extractedJob) {
	file, err := os.Create("jobs.csv")
	checkErr(err)
	w := csv.NewWriter(file)
	defer w.Flush()

	headers := []string{"ID", "Title", "Review", "Salary", "Grade"}

	wErr := w.Write(headers)
	checkErr(wErr)

	for _, job := range jobs {
		jobSlice := []string{"https://www.jobplanet.co.kr/companies/"+ job.id, job.title, job.review, job.salary, job.grade}
		jwErr := w.Write(jobSlice)
		checkErr(jwErr)
	}
}


func CleanString(str string) string{
	return strings.Join(strings.Fields(strings.TrimSpace(str)), "") 
}

func getPage(page int, url string, mainC chan<- []extractedJob){
	var jobs []extractedJob
	c := make(chan extractedJob)
	pageURL := url + "&page=" + strconv.Itoa(page)
	fmt.Println("Requesting", pageURL)
	res, err := http.Get(pageURL)
	checkErr(err)
	checkCode(res)

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkErr(err)

	searchCards := doc.Find(".content_wrap")
	searchCards.Each(func(i int, s *goquery.Selection) {
		go extractJob(s, c)
	})
	for i:=0; i<searchCards.Length(); i++ {
		job := <-c
		jobs = append(jobs, job)
	}
	mainC<- jobs
}

func extractJob(s *goquery.Selection, c chan<- extractedJob){
	btnId := s.Find(".btn_heart1")
	id, _ := btnId.Attr("data-company_id")
	title := CleanString(s.Find(".us_titb_l3>a").Text())
	review :=CleanString(s.Find(".content_col2_4>dt").Text())
	salary := CleanString(s.Find(".content_col2_4>dd>.us_stxt_1").Text())
	grade := s.Find(".content_col2_4>dd>.gfvalue").Text()
	c<- extractedJob{
		id: id, 
		title:title,
		review:review, 
		salary:salary, 
		grade:grade,
	}
}

func getPages(url string) int {
	pages := 0
	res, err := http.Get(url)
	checkErr(err)
	checkCode(res)

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkErr(err)
	doc.Find(".paginnation_new").Each(func(i int,s *goquery.Selection){
		pages = s.Find("a").Length()
	})
	return pages
}

func checkErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func checkCode(res *http.Response) {
	if res.StatusCode != 200 {
		log.Fatalln("Request failed with Status:", res.StatusCode)
	}
}
