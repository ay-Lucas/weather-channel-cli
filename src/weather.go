package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gocolly/colly/v2"
)

type Hour struct {
	Time         string
	Temperature  string
	PrecipChance string
	Condition    string
	FeelsLike    string
	Wind         string
	Humidity     string
	UVIndex      string
	CloudCover   string
	RainAmount   string
}

type Params struct {
	Query        string `json:"query"`
	Language     string `json:"language"`
	LocationType string `json:"locationType"`
}
type Body struct {
	Name   string `json:"name"`
	Params Params `json:"params"`
}

type Color string

const (
	Black  Color = "\u001b[30m"
	Red    Color = "\u001b[31m"
	Green  Color = "\u001b[32m"
	Yellow Color = "\u001b[33m"
	Blue   Color = "\u001b[34m"
	Reset  Color = "\u001b[0m"
)

func colorize(color Color, message string) string {
	return string(color) + message + string(Reset)
}

func main() {
	c := colly.NewCollector(
		colly.AllowedDomains("weather.com"),
		colly.MaxDepth(2),
		colly.Async(),
	)
	var hourly bool
	var numOfHours int
	var location string
	var locationUrl string
	flag.BoolVar(&hourly, "d", false, "Show hourly forecast")
	flag.IntVar(&numOfHours, "t", 24, "Specify number of hours")
	flag.StringVar(&location, "l", "none", "Location <City> <State>")
	flag.Parse()
	// readSettings()
	if len(flag.Arg(0)) > 0 {
		location = location + " " + flag.Arg(0)
	}
	fmt.Println(location)
	// os.Exit(1)
	//
	if location != "none" {
		locationUrl = postLocation(location)
		// locationUrl = handleRes(res)
	}
	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", "\nError:", err)
	})

	c.Wait()
	if len(locationUrl) == 0 {
		fmt.Println("No location")
		return
	}

	if hourly {
		hourlyForecast(c, locationUrl, numOfHours)
	} else {
		currentWeather(c, locationUrl)
	}
}

func writeSettings() {
	settingsDir := "/home/lucas/.local/share/weather-thing/"
	settingsFile := "data.json"
	// cwd, err := os.Getwd()
	fmt.Println(os.UserHomeDir())
	// if err != nil {
	// 	log.Fatal(err)
	// }
	if !doesPathExist(settingsDir) {
		err := os.Mkdir(settingsDir, 0755)
		if err != nil {
			log.Println(err)
		}
	}
	if !doesPathExist(settingsDir + settingsFile) {
		file, err := os.Create(settingsDir + settingsFile)
		if err != nil {
			log.Println(err)
		}
		defer file.Close()
	}
}

func doesPathExist(path string) bool {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		log.Println(err)
		return false
	}
	return true
}

func readSettings(path string) {
}

func isFileExists() {
}

func postLocation(query string) string {
	posturl := "https://weather.com/api/v1/p/redux-dal"
	bodyJson := []Body{
		{
			Name: "getSunV3LocationSearchUrlConfig",
			Params: Params{
				Query:        query,
				Language:     "en-US",
				LocationType: "locale",
			},
		},
	}
	jsonData, err := json.Marshal(bodyJson)
	// fmt.Printf("Json: %s", jsonData)
	if err != nil {
		log.Fatal(err)
	}

	response, err := http.Post(posturl, "application/json",
		bytes.NewBuffer(jsonData))
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()

	var target map[string]interface{}
	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("JSON reading error")
	}
	err1 := json.Unmarshal(body, &target)
	if err1 != nil {
		fmt.Println("Can not unmarshal JSON")
	}
	data := target["dal"].(map[string]interface{})["getSunV3LocationSearchUrlConfig"].(map[string]interface{})["language:en-US;locationType:locale;query:"+query].(map[string]interface{})["data"].(map[string]interface{})["location"]
	placeId := data.(map[string]interface{})["placeId"].([]interface{})[0]
	location := data.(map[string]interface{})["displayName"].([]interface{})[0]

	fmt.Println(placeId)
	// fmt.Println(data)
	fmt.Println(location)
	id := fmt.Sprintf("%v", placeId)
	return id
}

func getLocation(c *colly.Collector) string {
	var locationUrl string

	c.OnHTML("div.styles--OverflowNav--AWKwe.styles--overflowNav--cdZvZ", func(h *colly.HTMLElement) {
		url := h.ChildAttrs("a.ListItem--listItem--25ojW.styles--listItem--2CkF3.Button--default--2gfm1", "href")
		// url := h.ChildAttr("a.ListItem--listItem--25ojW.styles--listItem--2CkF3.Button--default--2gfm1", "href")
		fmt.Println(url)
		// split := strings.Split(url, "/")
		// locationUrl = split[len(split)-1]
		fmt.Println(locationUrl)
	})

	c.Visit("https://weather.com")

	return locationUrl
}

func hourlyForecast(c *colly.Collector, locationUrl string, numOfHours int) {
	hour := make([]Hour, 0)
	c.OnHTML("div.HourlyForecast--DisclosureList--MQWP6", func(e *colly.HTMLElement) {
		e.ForEach("details", func(i int, h *colly.HTMLElement) {
			item := Hour{}
			item.Temperature = h.ChildText("span.DetailsSummary--tempValue--jEiXE")
			item.Time = h.ChildText("h3.DetailsSummary--daypartName--kbngc")
			item.Condition = h.ChildText("div.DetailsSummary--condition--2JmHb")

			item.PrecipChance = h.ChildText("div.DetailsSummary--precip--1a98O span")

			detailsTable := h.ChildTexts("ul.DetailsTable--DetailsTable--3Bt2T li.DetailsTable--listItem--Z-5Vi div.DetailsTable--field--CPpc_ span.DetailsTable--value--2YD0-")
			item.FeelsLike = detailsTable[0]
			item.Wind = detailsTable[1]
			item.Humidity = detailsTable[2]
			item.UVIndex = detailsTable[3]
			item.CloudCover = detailsTable[4]
			item.RainAmount = detailsTable[5]
			hour = append(hour, item)
		})

		printHours(hour, numOfHours)
	})

	c.Visit("https://weather.com/weather/hourbyhour/l/" + locationUrl)
}

func printHours(hour []Hour, num int) {
	fmt.Println(colorize(Green, "Temperature\tFeels Like\tPrecipitation"))
	for i := 0; i < len(hour) || i < num; i++ {
		w := hour[i]
		// fmt.Printf("\n%s - Temperature: %s\nFeels Like: %s\nCondition: %s\nWind: %s\nHumidity: %s\nUVIndex: %s\nCloud Cover: %s\nPrecipitation: %s\n\n", w.Time, w.Temperature, w.Condition, w.FeelsLike, w.Wind, w.Humidity, w.UVIndex, w.CloudCover, w.RainAmount)
		// fmt.Println(colorize(Blue, w.Time) + " | Feels Like: " + w.FeelsLike + " Precipitation %: " + w.PrecipChance)
		fmt.Println(colorize(Blue, w.Time) + "\t\t" + w.FeelsLike + "\t\t" + w.PrecipChance)
	}
}

func currentWeather(c *colly.Collector, locationUrl string) {
	// c.OnHTML("div.styles--OverflowNav--AWKwe.styles--overflowNav--cdZvZ", func(h *colly.HTMLElement) {
	// 	url = h.ChildAttr("a.ListItem--listItem--25ojW.styles--listItem--2CkF3.styles--active--1ihhY.Button--default--2gfm1", "href")
	// 	fmt.Println(url)
	// })
	c.OnHTML("span.CurrentConditions--tempValue--MHmYY", func(e *colly.HTMLElement) {
		fmt.Println("Current Temperature:", string(Blue), e.Text, string(Reset))
	})
	todayUrl := "https://weather.com/weather/today/l/" + locationUrl

	c.Visit(todayUrl)

	c.Wait()
}
