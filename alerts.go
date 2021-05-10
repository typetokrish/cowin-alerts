package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type district struct {
	DistrictId   int    `json:"district_id"`
	DistrictName string `json:"district_name"`
}

type districtResponse struct {
	Districts []district `json:"districts"`
}

type fee struct {
	Vaccine string `json:"vaccine"`
	Fee     string `json:"fee"`
}
type session struct {
	SessionId         string   `json:"session_id"`
	Date              string   `json:"date"`
	AvailableCapacity int      `json:"available_capacity"`
	MinAgeLimit       int      `json:"min_age_limit"`
	Vaccine           string   `json:"vaccine"`
	Slots             []string `json:"slots"`
	From              string   `json:"from"`
	To                string   `json:"to"`
}
type center struct {
	CenterId      int    `json:"center_id"`
	CenterName    string `json:"name"`
	CenterAddress string `json:"address"`
	BlockName     string `json:"block_name"`
	Pincode       int    `json:"pincode"`
	FeeType       string `json:"fee_type"`
	VaccineFee    []fee
	Sessions      []session `json:"sessions"`
}

type CenterResponse struct {
	Centers []center `json:"centers"`
}

/**
 * Fetch the Districts from the CoWin Portal for Kerala
 */
func getDistricts() []district {
	//Api URL for the Districts
	url := "https://cdn-api.co-vin.in/api/v2/admin/location/districts/17"

	httpClient := &http.Client{Timeout: time.Second * 15}
	req, errR := http.NewRequest(http.MethodGet, url, nil)

	if errR != nil {
		fmt.Println("Request Error", errR)

	}
	//Define http headers
	req.Header.Add("content-type", "application/json")
	req.Header.Add("user-agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.72 Safari/537.36")

	resp, errDo := httpClient.Do(req)
	if errDo != nil {
		fmt.Println("Cannot Fetch the Districts Now", errDo)
	}
	defer resp.Body.Close()

	responseBytes, errResp := ioutil.ReadAll(resp.Body)
	fmt.Println(string(responseBytes))

	if errResp != nil {
		fmt.Println("Response Error", errResp)
	}
	//fmt.Println(responseBytes)

	//parse the api respons to a fixed struct / json
	var apiResponse districtResponse
	errParse := json.Unmarshal(responseBytes, &apiResponse)
	if errParse != nil {
		fmt.Println("Json Parse Error", errParse)
	}
	return apiResponse.Districts
}

func getAvailablCentersByDistrict(district district) []center {
	date := time.Now().Format("02-01-2006")
	url := "https://cdn-api.co-vin.in/api/v2/appointment/sessions/public/calendarByDistrict"
	url = url + "?district_id=" + strconv.Itoa(district.DistrictId) + "&date=" + date

	httpClient := &http.Client{Timeout: time.Second * 30}
	req, errR := http.NewRequest(http.MethodGet, url, nil)

	if errR != nil {
		fmt.Println("Request Error", errR)

	}
	//Define http headers
	req.Header.Add("content-type", "application/json")
	req.Header.Add("user-agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.72 Safari/537.36")

	resp, errDo := httpClient.Do(req)
	if errDo != nil {
		fmt.Println("Cannot Fetch the Centers Now", errDo)
	}
	defer resp.Body.Close()

	responseBytes, errResp := ioutil.ReadAll(resp.Body)
	//fmt.Println(string(responseBytes))

	if errResp != nil {
		fmt.Println("Response Error", errResp)
	}

	//parse the api respons to a fixed struct / json
	var apiResponse CenterResponse
	errParse := json.Unmarshal(responseBytes, &apiResponse)
	if errParse != nil {
		fmt.Println("Json Parse Error", errParse)
	}
	return apiResponse.Centers

}

func notifyTelgramChannel(district district, centers []center) {

	botKey := ""
	chatId := "@cowinalertskerala2021"
	var message string
	message = "Slots opened for district " + district.DistrictName
	message = message + "["
	for _, center := range centers {
		message = message + "  " + center.CenterName + ", "
	}
	message = url.QueryEscape(message + "]")

	url := "https://api.telegram.org/bot" + botKey + "/sendMessage?chat_id=" + chatId + "&text=" + message
	fmt.Println(url)

	httpClient := &http.Client{Timeout: time.Second * 30}
	req, errR := http.NewRequest(http.MethodGet, url, nil)

	if errR != nil {
		fmt.Println("Request Error", errR)

	}
	//Define http headers
	req.Header.Add("content-type", "application/json")
	req.Header.Add("user-agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.72 Safari/537.36")

	resp, errDo := httpClient.Do(req)
	if errDo != nil {
		fmt.Println("Cannot Fetch the Centers Now", errDo)
	}
	defer resp.Body.Close()

	responseBytes, errResp := ioutil.ReadAll(resp.Body)
	//fmt.Println(string(responseBytes))

	if errResp != nil {
		fmt.Println("Telegram Response Error", errResp)
	}

	fmt.Println("Telegram API", string(responseBytes))

}

func main() {
	fmt.Println("CoWin Alerts Main Started")
	districts := getDistricts()

	for _, district := range districts {
		centers := getAvailablCentersByDistrict(district)
		fmt.Println("Cnters", centers)
		if len(centers) == 0 {
			fmt.Println("There is no centers available for District", district)
		} else {
			fmt.Println("Notifying Channel", district)
			notifyTelgramChannel(district, centers)
		}
	}

}
