package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/tidwall/gjson"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
)

var (
	GET_WEATHER_FOR_AMAP_URL       = "https://restapi.amap.com/v3/weather/weatherInfo"
	GET_WEATHER_FOR_AMAP_TOKEN     = "ef1fd7c39e929320a38ff7185ae6fcff"
	SEND_MSG_TO_DINGTALK_ROBOT_URL = "https://oapi.dingtalk.com/robot/send?access_token=ffaabe93a835ff732b8053c0cd54c1e8315a8f906ddc0cc722dad5e833ff281c"
	//var SEND_MSG_TO_DINGTALK_ROBOT_TOKEN = "ffaabe93a835ff732b8053c0cd54c1e8315a8f906ddc0cc722dad5e833ff281c"
	GET_WEATHER_FOR_AMAP_CITYID = "110100" //北京代码
)

func main() {

	//var c = cron.New()
	////spec := "0 40 8 1/1 * ? "
	//spec := "0/3 * * * * * "
	//_, _ = c.AddFunc(spec, sendTodayMsg)
	//c.Start()
	//select {}
	sendTodayMsg()

}

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
			DayWeather   string `json:"dayweather"`
			NightWeather string `json:"nightweather"`
			DayTemp      string `json:"daytemp"`
			NightTemp    string `json:"nighttemp"`
			DayWind      string `json:"daywind"`
			NightWind    string `json:"nightwind"`
			DayPower     string `json:"daypower"`
			NightPower   string `json:"nightpower"`
		} `json:"casts"`
	} `json:"forecasts"`
}

type CityCode string
type DayOffset int8

const (
	Today    DayOffset = 0
	Tomorrow DayOffset = 1
)

//模块化包装
type AmapWeather struct {
	url   *url.URL
	token string
}

//工厂方法，组装url，减少传递过程中的参数
func NewAmapWeather(weatherUrl string, token string) (*AmapWeather, error) {
	parse, err := url.Parse(weatherUrl)
	if err != nil {
		log.Println("sender init failed:" + err.Error())
		return nil, err
	}
	return &AmapWeather{
		url:   parse,
		token: token,
	}, nil
}

//具体实现类
func (a *AmapWeather) GetWeather(city CityCode, dayOffset DayOffset) (*DayWeather, error) {
	params := url.Values{}
	params.Set("key", a.token)
	params.Set("city", string(city))
	params.Set("extensions", "all")
	a.url.RawQuery = params.Encode()
	urlPath := a.url.String()
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
	fmt.Println(string(body))
	//将获得的数据解组为结构
	var amapResp amapApiResponse
	err = json.Unmarshal(body, &amapResp)
	if err != nil {
		log.Println("amapResp Unmarshal fail:" + err.Error())
		return nil, err
	}
	s := amapResp.Forecasts[0]
	s2 := s.Casts[dayOffset]
	return &DayWeather{
		Date:         s2.Date,
		Week:         s2.Week,
		DayWeather:   s2.DayWeather,
		NightWeather: s2.NightWeather,
		DayTemp:      s2.DayTemp,
		NightTemp:    s2.NightTemp,
		DayWind:      s2.DayWind,
		NightWind:    s2.NightWind,
		DayPower:     s2.DayPower,
		NightPower:   s2.NightPower,
	}, nil
}

func sendTodayMsg() {
	//获取天气数据
	//msg, _ := getWeatherForAmap(GET_WEATHER_FOR_AMAP_URL, GET_WEATHER_FOR_AMAP_TOKEN, GET_WEATHER_FOR_AMAP_CITYID)
	//fmt.Println(msg)
	//组装URL+token
	amapWeather, err := NewAmapWeather(GET_WEATHER_FOR_AMAP_URL, GET_WEATHER_FOR_AMAP_TOKEN)
	if err != nil {
		log.Println("get weather error:" + err.Error())
	}
	//获取今天数据
	amapWeatherinfo, err2 := amapWeather.GetWeather(CityCode(GET_WEATHER_FOR_AMAP_CITYID), Today)
	if err2 != nil {

	}
	fmt.Println(amapWeatherinfo.Date + amapWeatherinfo.DayPower + amapWeatherinfo.NightTemp)
	//发送今天天气消息
	//robot, _ := sendWeatherMsgToDingTalkRobot(SEND_MSG_TO_DINGTALK_ROBOT_URL, msg)
	//fmt.Println(robot)
}

////从高德获取所有天气数据
//func getWeatherForAmap(weatherUrl string, weatherToken string, cityID string) (string, error) {
//	params := url.Values{}
//	Url, err := url.Parse(weatherUrl)
//	if err != nil {
//		log.Println("sender init failed: " + err.Error())
//		return "", err
//	}
//	params.Set("key", weatherToken)
//	params.Set("city", cityID)
//	params.Set("extensions", "all")
//	Url.RawQuery = params.Encode()
//	urlPath := Url.String()
//	fmt.Println(urlPath)
//	resp, err := http.Get(urlPath)
//	defer func(Body io.ReadCloser) {
//		err := Body.Close()
//		if err != nil {
//		}
//	}(resp.Body)
//	body, err := ioutil.ReadAll(resp.Body)
//	var weatherMsg = string(body)
//	fmt.Println(weatherMsg)
//	return weatherMsg, nil
//}

//向钉钉机器人发送当天天气消息
func sendWeatherMsgToDingTalkRobot(msgUrl string, msg string) (string, error) {
	client := &http.Client{}
	//拼接json参数
	data1 := make(map[string]interface{})
	data1["msgtype"] = "text"
	data2 := make(map[string]interface{})
	data2["isAtAll"] = false
	data1["at"] = data2
	data3 := make(map[string]interface{})
	//处理信息，从所有数据中拿去当天天气消息
	fmt.Println("===========================start process msg again after sendMsg==================================")
	afterMsg := dingTalkAfterSendMsgProcess(msg)
	fmt.Println(afterMsg)
	fmt.Println("===========================process msg again after sendMsg is already==================================")
	data3["content"] = "提醒：" + afterMsg
	data1["text"] = data3

	bytesData, err := json.Marshal(data1)
	if err != nil {
		return msg, err
	}

	req, _ := http.NewRequest(http.MethodPost, msgUrl, bytes.NewReader(bytesData))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := client.Do(req)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	fmt.Println(string(body))
	return msg, nil
}

//把天气数据处理为当天数据
func dingTalkAfterSendMsgProcess(msg string) string {

	week := gjson.Get(msg, `forecasts.#.casts.0.week`)
	dayWeather := gjson.Get(msg, `forecasts.#.casts.0.dayweather`)
	nightWeather := gjson.Get(msg, `forecasts.#.casts.0.nightweather`)
	highTemp := gjson.Get(msg, `forecasts.#.casts.0.daytemp`)
	lowTemp := gjson.Get(msg, `forecasts.#.casts.0.nighttemp`)
	dayWind := gjson.Get(msg, `forecasts.#.casts.0.daywind`)
	dayPower := gjson.Get(msg, `forecasts.#.casts.0.daypower`)
	nightWind := gjson.Get(msg, `forecasts.#.casts.0.nightwind`)
	nightPower := gjson.Get(msg, `forecasts.#.casts.0.nightpower`)
	today := time.Now().Format("2006-01-02")
	finallyMsg := "今天是" + today + "星期" + week.String()[2:len(week.String())-2] + "，今天白天" + dayWeather.String()[2:len(dayWeather.String())-2] + "，" +
		dayWind.String()[2:len(dayWind.String())-2] + "风" + dayPower.String()[2:len(dayPower.String())-2] + "级，" +
		"今天夜间" + nightWeather.String()[2:len(nightWeather.String())-2] + "，" + nightWind.String()[2:len(nightWind.String())-2] + "风" +
		nightPower.String()[2:len(nightPower.String())-2] + "级，" +
		"最高气温：" + highTemp.String()[2:len(highTemp.String())-2] + "度，最低气温：" + lowTemp.String()[2:len(lowTemp.String())-2] + "度"
	return finallyMsg
}
