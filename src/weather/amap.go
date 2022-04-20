package weather

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

type DayWeather struct {
	Date         string `json:"date"`
	Week         string `json:"week"`
	DayWeather   string `json:"dayweather"`
	NightWeather string `json:"nightweather"`
	DayTemp      string `json:"daytemp"`
	NightTemp    string `json:"nighttemp"`
	DayWind      string `json:"daywind"`
	NightWind    string `json:"nightwind"`
	DayPower     string `json:"daypower"`
	NightPower   string `json:"nightpower"`
}

type amapApiResponse struct {
	Status    string `json:"status"`
	Count     string `json:"count"`
	Info      string `json:"info"`
	Infocode  string `json:"infocode"`
	Forecasts []struct {
		City       string `json:"city"`
		Adcode     string `json:"adcode"`
		Province   string `json:"province"`
		Reporttime string `json:"reporttime"`
		Casts      []struct {
			Date         string `json:"date"`
			Week         string `json:"week"`
			Dayweather   string `json:"dayweather"`
			Nightweather string `json:"nightweather"`
			Daytemp      string `json:"daytemp"`
			Nighttemp    string `json:"nighttemp"`
			Daywind      string `json:"daywind"`
			Nightwind    string `json:"nightwind"`
			Daypower     string `json:"daypower"`
			Nightpower   string `json:"nightpower"`
		} `json:"casts"`
	} `json:"forecasts"`
}

type CityCode string
type DayOffset int8

const (
	Today    DayOffset = 0
	Tomorrow DayOffset = 1
)

// 模块化 整个包
type AmapWeather struct {
	//  模块依赖的数据、或者其他模块 这些东西从工厂方法里传进来
	url   *url.URL
	token string
}

// 工厂方法 工厂方法在main里组装
func NewAmapWeather(weatherUrl string, token string) (*AmapWeather, error) {
	_url, err := url.Parse(weatherUrl)
	if err != nil {
		log.Println("sender init failed: " + err.Error())
		return nil, err
	}
	return &AmapWeather{
		url:   _url,
		token: token,
	}, nil
}

// 模块具体功能 形参、返回值 一定要确切
func (a *AmapWeather) Get(city CityCode, dayOffset DayOffset) (*DayWeather, error) {
	params := url.Values{}
	params.Set("key", a.token)
	params.Set("city", string(city))
	params.Set("extensions", "all")
	thisUrl, _ := url.Parse(a.url.String())
	thisUrl.RawQuery = params.Encode()
	urlPath := thisUrl.String()
	fmt.Println(urlPath)
	resp, err := http.Get(urlPath)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
		}
	}(resp.Body)
	body, err := ioutil.ReadAll(resp.Body)
	//将获得的数据解组为结构
	var amapResp amapApiResponse
	err = json.Unmarshal(body, &amapResp)
	if err != nil {
		return nil, err
	}

	f := amapResp.Forecasts[0]
	//f.City
	c := f.Casts[0]

	return &DayWeather{
		Date:         c.Date,
		Week:         c.Week,
		DayWeather:   c.Dayweather,
		NightWeather: c.Nightweather,
		DayTemp:      c.Daytemp,
		NightTemp:    c.Nighttemp,
		DayWind:      "",
		NightWind:    "",
		DayPower:     "",
		NightPower:   "",
	}, nil
}
