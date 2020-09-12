package harvest

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

func HttpGet(url string, account_id string, token string) ([]byte, error) {
	log.WithFields(log.Fields{
		"url": url,
	}).Debug("HttpGet")
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("Harvest-Account-ID", account_id)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("User-Agent", "a9sharvest")

	client := &http.Client{
		Timeout: time.Second * 10,
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GET %s returned status %d body %v", url, resp.StatusCode, string(body))
	}

	return body, nil
}
