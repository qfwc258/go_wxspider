package common

import (
	"net/http"
	"io/ioutil"
)

var client = &http.Client{}

func GetMethod(url string)string{
	request, _ := http.NewRequest("GET", url, nil)
    response, _ := client.Do(request)
    defer response.Body.Close()
    if response.StatusCode == 200 {
        str, _ := ioutil.ReadAll(response.Body)
        bodyStr := string(str)
        return bodyStr
	}
	return ""
}