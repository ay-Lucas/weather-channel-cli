# Weather CLI Tool

A command-line application built in Go to fetch and display real-time weather forecasts using a combination of API integration and web scraping.

## Disclaimer: Potential Issues with Future Functionality

## This tool relies on the structure of Weather.com’s website and API responses. If Weather.com updates its website, modifies its API, or introduces new restrictions, the tool may no longer work as intended

## Features

- Retrieve current, hourly, and daily weather forecasts for any location
- Customizable forecast options
- Clean, colorized output
- Planned support for saving user preferences in a configuration file

---

## Prerequisites

- Go installed
- Internet connection

---

## Installation

1. Clone the repository:

   ```bash
   git clone https://github.com/yourusername/weather-cli-tool.git
   cd weather-cli-tool
   ```

2. Build the application:

   ```bash
   go build -o weather-cli .
   ```

---

## Usage

Run the CLI tool with the following options:

```bash
./weather-cli [flags]
```

### Flags

| Flag                          | Description                                                               | Default Value |
| ----------------------------- | ------------------------------------------------------------------------- | ------------- |
| `-l <city> <state> <country>` | Specify the location (e.g., "New York NY US"). **Required**.              | nil nil US    |
| `-hourly`                     | Display the hourly weather forecast.                                      | `false`       |
| `-d`                          | Display the daily weather forecast.                                       | `false`       |
| `-t <hours>`                  | Specify the number of hours for the hourly forecast.                      | `24`          |
| `-a`                          | List all available weather details (current, today, and hourly combined). | `false`       |
| `-temperature`                | Print only the current temperature.                                       | `false`       |

### Examples

#### Current Weather

Retrieve the current weather for a specific location:

```bash
./weather-cli -l Los Angeles, CA
```

#### Hourly Weather

Display the hourly weather forecast for the next 12 hours:

```bash
./weather-cli -l Los Angeles, CA -hourly -t 12
```

#### Daily Forecast

Show the daily weather forecast:

```bash
./weather-cli -l Los Angeles, CA -d
```

#### Full Weather Report

Get all details (current, today’s breakdown, and hourly forecast):

```bash
./weather-cli -l Los Angeles, CA -a
```

---

## Technical Details

- **API Integration**: Reverse-engineered Weather.com’s API to fetch location and weather data.
- **Web Scraping**: Utilized the `Colly` library to scrape additional weather details.
- **Customization**: Designed CLI flags for a personalized user experience.

---

## Contribution

Feel free to fork the repository and create a pull request for new features or bug fixes.

---

## License

This project is licensed under the [MIT License](LICENSE).
