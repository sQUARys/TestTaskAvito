package providers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type currencyProvider struct{}

func New() *currencyProvider {
	return &currencyProvider{}
}

func (prov *currencyProvider) Convert(from string, to string, amount float64) (float64, error) {
	format := "https://api.apilayer.com/exchangerates_data/convert?to=%s&from=%s&amount=%f"

	url := fmt.Sprintf(format, to, from, amount)

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, err
	}

	req.Header.Set("apikey", "OW7XhQPjG87aX7vEg5bz0glUbL2BgMeX")

	res, err := client.Do(req)
	if res.Body != nil {
		defer func() {
			err := res.Body.Close()
			if err != nil {
				log.Println(err)
			}
		}()
	}

	jsonMap := make(map[string]interface{})
	body, err := ioutil.ReadAll(res.Body)

	err = json.Unmarshal(body, &jsonMap)
	if err != nil {
		return 0, err
	}

	return jsonMap["result"].(float64), nil
}
