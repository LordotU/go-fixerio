// Package gofixerio providers simple wrapper for "Foreign exchange rates and currency conversion JSON API" (https://fixer.io).
//
// Using this package is easy as 1-2-3:
//
//     fixerio := gofixerio.New("your API key here", "EUR", "false")
//     latestRates, err := fixerio.GetLatest()
//     if err != nil {
//             log.Panic(err)
//     }
//     log.Printf("%+v", latestRates)
//
package gofixerio

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

// FixerIO represents a client instance
type FixerIO struct {
	APIKey string
	Base   string // optional
	Secure bool   // optional
}

// BaseURL is the main URL for accessing https://fixer.io API
const BaseURL = "data.fixer.io"

// New eturns pointer to a new instance of FixerIO
func New(APIKey string, Base string, Secure bool) (*FixerIO, error) {
	if APIKey == "" {
		return nil, wrapError(errors.New("APIKey is required"))
	}

	return &FixerIO{
		APIKey: APIKey,
		Base:   strings.ToUpper(Base),
		Secure: Secure,
	}, nil
}

// SetBase sets new base currency for all subsequent requests
func (f *FixerIO) SetBase(base string) *FixerIO {
	f.Base = strings.ToUpper(base)
	return f
}

// GetSymbols returns pointer to ResponseSymbols struct which contains symbols map (map[string]string)
func (f *FixerIO) GetSymbols() (*ResponseSymbols, error) {
	url := f.getURL("symbols", map[string]string{})
	var response ResponseSymbols
	err := f.makeRequest(url, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// GetLatest returns pointer to ResponseLatest struct which contains latest rates map (map[string]float64) for given symbols arg
func (f *FixerIO) GetLatest(symbols []string) (*ResponseLatest, error) {
	url := f.getURL("latest", map[string]string{
		"symbols": strings.Join(symbols[:], ","),
	})
	var response ResponseLatest
	err := f.makeRequest(url, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// GetHistorical returns pointer to ResponseHistorical struct which contains historical rates map (map[string]float64) for given date and symbols args
func (f *FixerIO) GetHistorical(date string, symbols []string) (*ResponseHistorical, error) {
	url := f.getURL(date, map[string]string{
		"symbols": strings.Join(symbols[:], ","),
	})
	var response ResponseHistorical
	err := f.makeRequest(url, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// GetConversion returns pointer to ResponseConversion struct which contains conversion data for given from, to, amount and date args
func (f *FixerIO) GetConversion(from string, to string, amount float64, date string) (*ResponseConversion, error) {
	url := f.getURL("convert", map[string]string{
		"from":   from,
		"to":     to,
		"amount": strconv.FormatFloat(float64(amount), 'f', -1, 64),
		"date":   date,
	})
	var response ResponseConversion
	err := f.makeRequest(url, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// GetTimeseries returns pointer to ResponseTimeseries struct which contains timeseries rates map (map[string][string]float64) for given startDate and endDate arg
func (f *FixerIO) GetTimeseries(startDate string, endDate string) (*ResponseTimeseries, error) {
	url := f.getURL("timeseries", map[string]string{
		"start_date": startDate,
		"end_data":   endDate,
	})
	var response ResponseTimeseries
	err := f.makeRequest(url, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// GetFluctuation returns pointer to ResponseTimeseries struct which contains fluctuation rates map (map[string][string]float64) for given startDate and endDate arg
func (f *FixerIO) GetFluctuation(startDate string, endDate string) (*ResponseFluctuation, error) {
	url := f.getURL("fluctuation", map[string]string{
		"start_date": startDate,
		"end_data":   endDate,
	})
	var response ResponseFluctuation
	err := f.makeRequest(url, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

func (f *FixerIO) getURL(method string, params map[string]string) string {
	var url bytes.Buffer

	if f.Secure {
		url.WriteString("https://")
	} else {
		url.WriteString("http://")
	}

	url.WriteString(BaseURL)
	url.WriteString("/api")
	url.WriteString("/" + method)
	url.WriteString("?access_key=")
	url.WriteString(string(f.APIKey))
	url.WriteString("&base=")
	url.WriteString(f.Base)

	for name, value := range params {
		if value != "" {
			url.WriteString("&" + name + "=" + value)
		}
	}

	return url.String()
}

func (f *FixerIO) makeRequest(url string, result interface{}) error {
	response, err := http.Get(url)
	if err != nil {
		return wrapError(err)
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return wrapError(err)
	}

	var responseError ResponseError
	err = json.Unmarshal(body, &responseError)
	if err != nil {
		return wrapError(err)
	}

	if !responseError.Success && responseError.Error.Code != 0 {
		if responseError.Error.Info != "" {
			return wrapError(errors.New("fixer.io request error with info: " + responseError.Error.Info))
		}
		if responseError.Error.Type != "" {
			return wrapError(errors.New("fixer.io request error with type: " + responseError.Error.Type))
		}
		return wrapError(errors.New("fixer.io request error with code: " + strconv.Itoa(responseError.Error.Code)))
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return wrapError(err)
	}

	return nil
}

func wrapError(err error) error {
	return errors.Wrap(err, "gofixerio error")
}
