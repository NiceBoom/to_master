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

func main() {
	var GET_WEATHER_FOR_AMAP_URL = "https://restapi.amap.com/v3/weather/weatherInfo"
	var GET_WEATHER_FOR_AMAP_TOKEN = "ef1fd7c39e929320a38ff7185ae6fcff"
	var SEND_MSG_TO_DINGTALK_ROBOT_URL = "https://oapi.dingtalk.com/robot/send?access_token=ffaabe93a835ff732b8053c0cd54c1e8315a8f906ddc0cc722dad5e833ff281c"
	//var SEND_MSG_TO_DINGTALK_ROBOT_TOKEN = "ffaabe93a835ff732b8053c0cd54c1e8315a8f906ddc0cc722dad5e833ff281c"
	var GET_WEATHER_FOR_AMAP_CITYID = "110100" //北京代码
	fmt.Println("=============start get weather msg=====================")
	msg, _ := getWeatherForAmap(GET_WEATHER_FOR_AMAP_URL, GET_WEATHER_FOR_AMAP_TOKEN, GET_WEATHER_FOR_AMAP_CITYID)
	fmt.Println(msg)
	fmt.Println("============= get weather msg is already =====================")
	fmt.Println("============= start send weather msg to dingTalk robot =====================")
	robot, _ := sendWeatherMsgToDingTalkRobot(SEND_MSG_TO_DINGTALK_ROBOT_URL, msg)
	fmt.Println(robot)
	fmt.Println("============= send weather msg to dingTalk robot of already =====================")
}

//从高德获取所有天气数据
func getWeatherForAmap(weatherUrl string, weatherToken string, cityID string) (string, error) {
	params := url.Values{}
	Url, err := url.Parse(weatherUrl)
	if err != nil {
		log.Println("sender init failed: " + err.Error())
		return "", err
	}
	params.Set("key", weatherToken)
	params.Set("city", cityID)
	params.Set("extensions", "all")
	Url.RawQuery = params.Encode()
	urlPath := Url.String()
	fmt.Println(urlPath)
	resp, err := http.Get(urlPath)
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
		}
	}(resp.Body)
	body, err := ioutil.ReadAll(resp.Body)
	var weatherMsg = string(body)
	fmt.Println(weatherMsg)
	return weatherMsg, nil
}

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
