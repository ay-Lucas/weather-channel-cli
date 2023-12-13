package main
type Hour struct {
    Time string
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
type Data struct {
    Location struct {
        Address           []string  `json:"address"`
        AdminDistrict     []any     `json:"adminDistrict"`
        AdminDistrictCode []any     `json:"adminDistrictCode"`
        City              []string  `json:"city"`
        Country           []string  `json:"country"`
        CountryCode       []string  `json:"countryCode"`
        DisplayName       []string  `json:"displayName"`
        IanaTimeZone      []string  `json:"ianaTimeZone"`
        Latitude          []float64 `json:"latitude"`
        Locale            []struct {
            Locale1 any    `json:"locale1"`
            Locale2 string `json:"locale2"`
            Locale3 any    `json:"locale3"`
            Locale4 any    `json:"locale4"`
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
}
type Location struct {
	Location struct {
		City    []string `json:"city"`
		Country []string `json:"country"`
		PlaceId []string `json:"placeId"`
	} `json:"location"`
}
type Weather struct {
    Time string
    Conditions string
    FeelsLike string
    Temperature string
    Humidity string
    Pressure string
    Visibility string
    WindSpeed string
    DewPoint string
    UvIndex string
    MoonPhase string
    Forecast Forecast
}
type Today struct {
    PrecipPropability string
    Temperature string
}
type Forecast struct {
    Morning Today
    Afternoon Today
    Evening Today
    Overnight Today
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

