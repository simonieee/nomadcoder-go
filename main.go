package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/labstack/echo"
	"github.com/simonieee/learngo/scrapper"
)

const fileName string = "jobs.csv"

func handleHome(c echo.Context) error {
	return c.File("home.html")
}

func handleScrape(c echo.Context) error {
	defer os.Remove(fileName)
	term1 := strings.ToLower(scrapper.CleanString(c.FormValue("service"))) 
	term2 := strings.ToLower(scrapper.CleanString(c.FormValue("city"))) 
	fmt.Println(c.FormValue("service"), c.FormValue("city"))
	scrapper.Scrape(term1, term2)
	return c.Attachment(fileName, fileName)
}

func main() {
	e := echo.New()
	e.GET("/",handleHome)
	e.POST("/scrape",handleScrape)
	e.Logger.Fatal(e.Start(":1323"))
}