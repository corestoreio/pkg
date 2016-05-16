package geoip

import (
	"encoding/json"
	"net"
	"net/http"
	"time"

	"github.com/corestoreio/csfw/net/geoip/internal"
)

type mmws struct {
	userID     string
	licenseKey string
	client     *http.Client
}

func newMMWS(userID, licenseKey string, timeout time.Duration) *mmws {
	if timeout == 0 {
		timeout = time.Second * 20
	}
	ws := &mmws{
		userID:     userID,
		licenseKey: licenseKey,
		client:     &http.Client{Timeout: timeout},
	}
	return ws
}

func (mm *mmws) Country(ipAddress net.IP) (*Country, error) {
	return mm.fetch("https://geoip.maxmind.com/geoip/v2.1/country/", ipAddress)
}

func (mm *mmws) Close() error {
	return nil
}

func (a *mmws) City(ipAddress net.IP) (internal.Response, error) {
	return a.fetch("https://geoip.maxmind.com/geoip/v2.1/city/", ipAddress)
}

func (a *mmws) Insights(ipAddress net.IP) (internal.Response, error) {
	return a.fetch("https://geoip.maxmind.com/geoip/v2.1/insights/", ipAddress)
}

func (a *mmws) fetch(prefix string, ipAddress net.IP) (internal.Response, error) {
	var response internal.Response
	req, err := http.NewRequest("GET", prefix+ipAddress.String(), nil)
	if err != nil {
		return response, err
	}

	// authorize the request
	// http://dev.maxmind.com/geoip/geoip2/web-services/#Authorization
	req.SetBasicAuth(a.userID, a.licenseKey)

	// execute the request

	resp, err := a.client.Do(req)
	if err != nil {
		return response, err
	}
	defer resp.Body.Close()

	// handle errors that may occur
	// http://dev.maxmind.com/geoip/geoip2/web-services/#Response_Headers
	if resp.StatusCode >= 400 && resp.StatusCode < 600 {
		v := internal.Error{}
		err := json.NewDecoder(resp.Body).Decode(&v)
		if err != nil {
			return response, err
		}

		return response, v
	}

	// parse the response body
	// http://dev.maxmind.com/geoip/geoip2/web-services/#Response_Body

	err = json.NewDecoder(resp.Body).Decode(&response)
	return response, err
}
