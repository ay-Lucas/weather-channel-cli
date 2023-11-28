package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gocolly/colly/v2"
)

type Response struct {
	Dal struct {
		GetSunV3LocationSearchURLConfig struct {
			LanguageEnUSLocationTypeLocaleQueryMiami struct {
				StatusText string `json:"statusText"`
				Data       struct {
					Location struct {
						Address           []string  `json:"address"`
						AdminDistrict     []string  `json:"adminDistrict"`
						AdminDistrictCode []any     `json:"adminDistrictCode"`
						City              []string  `json:"city"`
						Country           []string  `json:"country"`
						CountryCode       []string  `json:"countryCode"`
						DisplayName       []string  `json:"displayName"`
						IanaTimeZone      []string  `json:"ianaTimeZone"`
						Latitude          []float64 `json:"latitude"`
						Locale            []struct {
							Locale3 any    `json:"locale3"`
							Locale4 any    `json:"locale4"`
							Locale1 string `json:"locale1"`
							Locale2 string `json:"locale2"`
						} `json:"locale"`
						Longitude            []float64 `json:"longitude"`
						Neighborhood         []any     `json:"neighborhood"`
						PlaceID              []string  `json:"placeId"`
						PostalCode           []string  `json:"postalCode"`
						PostalKey            []string  `json:"postalKey"`
						DisputedArea         []bool    `json:"disputedArea"`
						DisputedCountries    []any     `json:"disputedCountries"`
						DisputedCountryCodes []any     `json:"disputedCountryCodes"`
						DisputedCustomers    []any     `json:"disputedCustomers"`
						DisputedShowCountry  [][]bool  `json:"disputedShowCountry"`
						IataCode             []string  `json:"iataCode"`
						IcaoCode             []string  `json:"icaoCode"`
						LocID                []string  `json:"locId"`
						LocationCategory     []any     `json:"locationCategory"`
						PwsID                []string  `json:"pwsId"`
						Type                 []string  `json:"type"`
					} `json:"location"`
				} `json:"data"`
				Status  int  `json:"status"`
				Loading bool `json:"loading"`
				Loaded  bool `json:"loaded"`
			} `json:"language:en-US;locationType:locale;query:miami"`
		} `json:"getSunV3LocationSearchUrlConfig"`
	} `json:"dal"`
}
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
	flag.BoolVar(&hourly, "d", false, "Show hourly forecast")
	flag.IntVar(&numOfHours, "t", 24, "Specify number of hours")
	flag.StringVar(&location, "l", "none", "Location <City> <State>")
	flag.Parse()
	readSettings()
	res := postLocation()
	locationUrl := handleRes(res)
	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", "\nError:", err)
	})

	c.Wait()

	if hourly {
		hourlyForecast(c, numOfHours)
	} else {
		currentWeather(c, locationUrl)
	}
}

func readSettings() {
	settingsDir := "~/.local/share/weather-thing"
	// cwd, err := os.Getwd()
	// if err != nil {
	// log.Fatal(err)
	// }
	// fmt.Println(dirs)
	os.ReadDir(settingsDir)
	stat, err := os.Stat(settingsDir)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(stat)
}

func writeSettings() {
	// settingsDir := "~/.local/share/weather-thing"
}

func isFileExists() {
}

func postLocation() *http.Response {
	posturl := "https://weather.com/api/v1/p/redux-dal"
	bodyJson := []Body{
		{
			Name: "getSunV3LocationSearchUrlConfig",
			Params: Params{
				Query:        "miami",
				Language:     "en-US",
				LocationType: "locale",
			},
		},
	}
	json_data, err := json.Marshal(bodyJson)
	// fmt.Printf("Json: %s", json_data)
	if err != nil {
		log.Fatal(err)
	}

	resp, err := http.Post(posturl, "application/json",
		bytes.NewBuffer(json_data))
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	return resp
}

func handleRes(resp *http.Response) string { // returns URL for querried locatoin
	var response Response

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("JSON reading error")
	}
	err1 := json.Unmarshal(body, &response)
	if err1 != nil {
		fmt.Println("Can not unmarshal JSON")
	}

	placeId := response.Dal.GetSunV3LocationSearchURLConfig.LanguageEnUSLocationTypeLocaleQueryMiami.Data.Location.PlaceID[0]
	fmt.Println(placeId)
	return placeId
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

func hourlyForecast(c *colly.Collector, numOfHours int) {
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

	c.Visit("https://weather.com/weather/hourbyhour/l/73a00f6a54fd626905dc1ae45c472c895a92da1e34591b680d45ec516f941349")
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
