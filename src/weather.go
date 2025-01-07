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
	var weather Weather
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

	if hourly {
		getHourlyForecast(c, locationUrl, numOfHours, &weather)
		printHours(weather.Hourly, numOfHours)
	} else {
		if listAll {
			getHourlyForecast(c, locationUrl, numOfHours, &weather)
		}
		getCurrentWeather(c, locationUrl, &weather)
		printCurrent(weather, listAll, numOfHours)
	}
}
func printCurrent(weather Weather, listAll bool, numOfHours int) {
	currentValues := reflect.ValueOf(weather.Current)
	currentTypes := currentValues.Type()
	todayValues := reflect.ValueOf(weather.Today)
	todayTypes := todayValues.Type()

	// hourlyValues := reflect.ValueOf(weather.Hourly[0])
	// hourlyTypes := hourlyValues.Type()
	var hoursToPrint int
	if numOfHours < len(weather.Hourly) {
		hoursToPrint = numOfHours
	} else {
		hoursToPrint = len(weather.Hourly)
	}

	if listAll {
		fmt.Println("Current Weather:")
		for i := 0; i < currentValues.NumField(); i++ {
			fmt.Println(currentTypes.Field(i).Name, " ", string(Blue), currentValues.Field(i), string(Reset))

		}
		fmt.Println("\nToday:")
		for i := 0; i < todayValues.NumField(); i++ {
			fmt.Println(todayTypes.Field(i).Name, " ", string(Blue), todayValues.Field(i), string(Reset))
		}

		fmt.Println("\nHourly:")
		// fmt.Println(colorize(Green, "Time | Temperature | Precipitation Chance | Feels Like | Wind | Humidity | UV Index | Cloud Cover | Rain Amount"))
		title := [10]string{"Time", "Temperature", "Precip %", "Description", "Feels Like", "Wind", "Humidity", "UV Index", "Cloud Cover", "Rain Amount"}

		// fmt.Printf("%-10v\n", title)

		fmt.Printf(string(Blue))
		for i := 0; i < len(title); i++ {
			fmt.Printf("%-15s", title[i])
		}
		fmt.Printf("%s\n", string(Reset))
		for i := 0; i < hoursToPrint; i++ {
			fields := weather.Hourly[i]
			values := reflect.ValueOf(fields)
			// fmt.Println(whitespace(split, len(title), fields, "|"))
			for i := 0; i < 10; i++ {
				fmt.Printf("%-15s", values.Field(i).String())
			}
			fmt.Println()
		}
	} else {
		fmt.Println("Current Temperature:", string(Blue), weather.Current.Temperature, string(Reset))
	}

}

func parseMap(data map[string]interface{}, copy map[string]interface{}) map[string]interface{} {
	for k, v := range data {
		if k == "location" {
			copy[k] = v
			return copy
		} else {
			if reflect.TypeOf(v).String() != "map[string]interface {}" {
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

func mapToStruct(target map[string]interface{}) Data {
	locationMap := make(map[string]interface{})
	locationMap = parseMap(target, locationMap)
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

func getHourlyForecast(c *colly.Collector, locationUrl string, numOfHours int, weather *Weather) {
	c.OnHTML("div.HourlyForecast--DisclosureList--MQWP6", func(e *colly.HTMLElement) {
		e.ForEach("details", func(i int, h *colly.HTMLElement) {
			hour := Hour{}
			hour.Temperature = h.ChildText("span.DetailsSummary--tempValue--jEiXE")
			hour.Time = h.ChildText("h3.DetailsSummary--daypartName--kbngc")
			hour.Description = h.ChildText("span.DetailsSummary--extendedData--307Ax")

			hour.PrecipChance = h.ChildText("div.DetailsSummary--precip--1a98O span")

			detailsTable := h.ChildTexts("ul.DetailsTable--DetailsTable--3Bt2T li.DetailsTable--listItem--Z-5Vi div.DetailsTable--field--CPpc_ span.DetailsTable--value--2YD0-")
			hour.FeelsLike = detailsTable[0]
			hour.Wind = detailsTable[1]
			hour.Humidity = detailsTable[2]
			hour.UVIndex = detailsTable[3]
			hour.CloudCover = detailsTable[4]
			hour.RainAmount = detailsTable[5]
			weather.Hourly = append(weather.Hourly, hour)
		})
	})
	hourlyUrl := "https://weather.com/weather/hourbyhour/l/" + locationUrl
	c.Visit(hourlyUrl)
}

func printHours(hour []Hour, num int) {
	fmt.Println(colorize(Green, "Temperature\tFeels Like\tPrecipitation"))
	for i := 0; i < len(hour) || i < num; i++ {
		w := hour[i]
		fmt.Println(colorize(Blue, w.Time) + "\t\t" + w.FeelsLike + "\t\t" + w.PrecipChance)
	}
}
func getCurrentWeather(c *colly.Collector, locationUrl string, weather *Weather) {
	timestampRegex := regexp.MustCompile("[0-9]{1,2}:[0-9]{2} [APMapm]{2} [A-Z]+")
	numberThenUnitRegex := regexp.MustCompile("[0-9]+.[0-9]+(| )[A-z]+")
	percentageRegex := regexp.MustCompile("[0-9]{1,3}%")
	tag := "[data-testid=wxData]"
	c.OnHTML("main#MainContent", func(e *colly.HTMLElement) {
		table := e.DOM.Find("div.TodayDetailsCard--detailsContainer--2yLtL")
		todayTable := e.DOM.Find("ul.WeatherTable--columns--6JrVO.WeatherTable--wide--KY3eP")
		weather.Today.Morning.Temperature = todayTable.Find("li:nth-child(1) > a > div> span[data-testid=TemperatureValue]").Text()
		weather.Today.Afternoon.Temperature = todayTable.Find("li:nth-child(1) > a > div> span[data-testid=TemperatureValue]").Text()
		weather.Today.Evening.Temperature = todayTable.Find("li:nth-child(1) > a > div> span[data-testid=TemperatureValue]").Text()
		weather.Today.Overnight.Temperature = todayTable.Find("li:nth-child(1) > a > div> span[data-testid=TemperatureValue]").Text()

		weather.Today.Morning.PrecipProb = percentageRegex.FindString(todayTable.Find("li:nth-child(1) > a > div.Column--precip--3JCDO").Text())
		weather.Today.Afternoon.PrecipProb = percentageRegex.FindString(todayTable.Find("li:nth-child(2) > a > div.Column--precip--3JCDO").Text())
		weather.Today.Evening.PrecipProb = percentageRegex.FindString(todayTable.Find("li:nth-child(3) > a > div.Column--precip--3JCDO").Text())
		weather.Today.Overnight.PrecipProb = percentageRegex.FindString(todayTable.Find("li:nth-child(4) > a > div.Column--precip--3JCDO").Text())

		weather.Current.Time = timestampRegex.FindString(e.ChildText("span.CurrentConditions--timestamp--1ybTk"))
		weather.Current.Description = e.ChildText("div [data-testid=wxPhrase]")
		weather.Current.Temperature = e.ChildText(".CurrentConditions--tempValue--MHmYY")
		weather.Current.FeelsLike = e.ChildText("span.TodayDetailsCard--feelsLikeTempValue--2icPt")
		weather.Current.WindSpeed = table.Find(":nth-child(2) > " + tag + " > span > span + span").Text()
		weather.Current.Humidity = table.Find(":nth-child(3) > " + tag).Text()
		weather.Current.Pressure = numberThenUnitRegex.FindString(table.Find(":nth-child(5) > " + tag + " > span").Text())
		weather.Current.Visibility = table.Find(":nth-child(7) > " + tag).Text()
		weather.Current.DewPoint = table.Find(":nth-child(4) > " + tag).Text()
		weather.Current.MoonPhase = table.Find(":nth-child(8) > " + tag).Text()
		weather.Current.UvIndex = table.Find(":nth-child(6) > " + tag).Text()
	})
	todayUrl := "https://weather.com/weather/today/l/" + locationUrl
	c.Visit(todayUrl)
}
