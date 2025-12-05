package util

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/spf13/viper"
)

type IPResponse struct {
	Status    string `json:"status"`
	Info      string `json:"info"`
	InfoCode  string `json:"infocode"`
	Province  string `json:"province"`
	City      string `json:"city"`
	AdCode    string `json:"adcode"`
	Rectangle string `json:"rectangle"`
}

// GetLocationByIP 根据IP获取地址信息
func GetLocationByIP(ip string) (string, string, error) {
	data := IPResponse{}
	key := viper.GetString("gaoDe.key")
	url := viper.GetString("gaoDe.ip_url")
	method := "GET"
	params := map[string]string{
		"ip":  ip,
		"key": key,
	}
	res, err := HttpRequest(url, method, nil, params, nil)
	if err != nil {
		return "", "", err
	}
	defer res.Body.Close()
	if res.StatusCode == http.StatusOK {
		return "", "", fmt.Errorf("request failed with status code: %d", res.StatusCode)
	}
	byteData, err := io.ReadAll(res.Body)
	if err != nil {
		return "", "", err
	}
	if err := json.Unmarshal(byteData, &data); err != nil {
		return "", "", err
	}
	return data.Province, data.City, nil
}
