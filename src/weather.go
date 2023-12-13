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
	"reflect"
	"regexp"

	"github.com/gocolly/colly/v2"
)


func colorize(color Color, message string) string {
	return string(color) + message + string(Reset)
}

func main() {
	c := colly.NewCollector(
		colly.AllowedDomains("weather.com"),
		colly.MaxDepth(2),
		// colly.Async(),
	)
	var hourly bool
	numOfHours := 24
	var location string
    var daily bool
    var listAll bool
    var temperature bool
    // var current bool
	flag.BoolVar(&hourly, "hourly", false, "Show hourly forecast")
	flag.BoolVar(&daily, "d", false, "Show daily forecast")
    // flag.BoolVar(&current, "c", true, "Show current weather")
	flag.IntVar(&numOfHours, "t", 24, "Specify number of hours")
	flag.StringVar(&location, "l", "none", "Specify location <City> <State> (required)")
    flag.BoolVar(&listAll, "a", false, "List all results")
    flag.BoolVar(&temperature, "temperature", false, "Print current temperature")
	flag.Parse()
    if location == "none" {
        fmt.Println("Need to specify a location <City> State> <Country>")
        return
    }
	if len(flag.Arg(0)) > 0 {
		location = location + " " + flag.Arg(0)
	}
	var locationUrl string
    var data Data
    data = getData(location)
    locationUrl = data.Location.PlaceID[0]
    fmt.Printf("%s, %s\n", data.Location.DisplayName[0], data.Location.Country[0])
    c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", "\nError:", err)
	})

	if len(locationUrl) == 0 {
		fmt.Println("No location")
		return
	}
	// fmt.Printf("args %s", hourly)
	if hourly {
		hourlyForecast(c, locationUrl, numOfHours)
	} else {
        weather := getCurrentWeather(c, locationUrl)
        printCurrent(weather, listAll)
	}
}
func printCurrent(weather Weather, listAll bool) {
    values := reflect.ValueOf(weather)
    types := values.Type()
    if listAll {
        for i := 0; i < values.NumField(); i++ {
            fmt.Println(types.Field(i).Name, " ", string(Blue), values.Field(i), string(Reset))
        }
    } else {
        fmt.Println("Current Temperature:", string(Blue), weather.Temperature, string(Reset))
    }


}

func parseMap(data map[string]interface{}, copy map[string]interface{}) map[string]interface{} {
    for k, v := range data {
        if k == "location" {
            copy[k] = v
            return copy
        } else {
            if reflect.TypeOf(v).String() != "map[string]interface {}"  { 
                continue
            }
            return parseMap(v.(map[string]interface{}), copy)
        }
    }
    return copy
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
func getData(query string) Data {
    resMap := postLocation(query)
    data := mapToStruct(resMap)
    return data
}

func mapToStruct(target map[string]interface{} ) Data {
    locationMap := make(map[string]interface{})
    locationMap = parseMap(target, locationMap)
    // for k, v := range locationMap {
    //     fmt.Printf("key: %s value: %s\n", k, v)
    // }
    locationJson, err := json.Marshal(locationMap)
    if err != nil {
        log.Fatal(err)
    }
    data := Data{}
    err = json.Unmarshal(locationJson, &data)
    if err != nil {
        log.Fatal(err)
    }

    return data

}
func postLocation(query string) map[string]interface{} {
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
	err = json.Unmarshal(body, &target)
	if err != nil {
		fmt.Println("Can not unmarshal JSON")
	}
    return target
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
	hourlyUrl := "https://weather.com/weather/hourbyhour/l/" + locationUrl
	c.Visit(hourlyUrl)
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

func getCurrentWeather(c *colly.Collector, locationUrl string) Weather {
    currentWeather := Weather{}
    timestampRegex := regexp.MustCompile("[0-9]{1,2}:[0-9]{2} [APMapm]{2} [A-Z]+")
    numberThenUnitRegex := regexp.MustCompile("[0-9]+.[0-9]+(| )[A-z]+")
    tag :=  "[data-testid=wxData]"
    c.OnHTML("main#MainContent", func(e *colly.HTMLElement) {
        table := e.DOM.Find("div.TodayDetailsCard--detailsContainer--2yLtL")
        // todayTable := e.DOM.Find(".WeatherTable--columns--6JrVO.WeatherTable--wide--KY3eP")
        // currentWeather.Forecast.Overnight.PrecipPropability = todayTable.Find(":nth-child(4) > span.Column--precip--3JCDO").Text()
        fmt.Println(currentWeather.Forecast.Overnight.PrecipPropability)
        
        currentWeather = Weather {
            Time: timestampRegex.FindString(e.ChildText("span.CurrentConditions--timestamp--1ybTk")),
            Conditions: e.ChildText("div [data-testid=wxPhrase]"),
            Temperature: e.ChildText(".CurrentConditions--tempValue--MHmYY"),
            FeelsLike: e.ChildText("span.TodayDetailsCard--feelsLikeTempValue--2icPt"),
            WindSpeed: table.Find(":nth-child(2) > " + tag + " > span > span + span").Text(),
            Humidity: table.Find(":nth-child(3) > " + tag).Text(),
            Pressure: numberThenUnitRegex.FindString(table.Find(":nth-child(5) > " + tag + " > span").Text()),
            Visibility: table.Find(":nth-child(7) > " + tag).Text(),
            DewPoint: table.Find(":nth-child(4) > " + tag).Text(),
            MoonPhase: table.Find(":nth-child(8) > " + tag).Text(), 
            UvIndex: table.Find(":nth-child(6) > " + tag).Text(),
        }
      })
    todayUrl := "https://weather.com/weather/today/l/" + locationUrl
    c.Visit(todayUrl)
    return currentWeather
}
