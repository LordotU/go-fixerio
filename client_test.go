package gofixerio_test

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	gofixerio "github.com/LordotU/go-fixerio"
)

var referenceSymbols = []string{
	"USD", "EUR", "JPY", "GBP", "AUD", "CAD", "CHF", "CNY", "NZD", "RUB",
}

func getClient() (*gofixerio.FixerIO, error) {
	APIKey := os.Getenv("FIXERIO_API_KEY")
	if APIKey == "" {
		fmt.Println("APIKey shoud be specified for tests running")
		os.Exit(1)
	}

	return gofixerio.New(
		APIKey,
		"eur",
		false,
	)
}

func TestNew(t *testing.T) {
	_, err := gofixerio.New("", "", false)
	assert.EqualError(t, err, "gofixerio error: APIKey is required")

	_, err = gofixerio.New("186f2ecd6621d35c3319bb53e1787d90", "", false)
	assert.Nil(t, err)
}

func TestSetBase(t *testing.T) {
	client, _ := getClient()
	assert.Equal(t, "EUR", client.Base)

	client.SetBase("usd")
	assert.Equal(t, "USD", client.Base)
}

func TestGetSymbols(t *testing.T) {
	client, _ := getClient()

	response, err := client.GetSymbols()
	if err != nil {
		t.Error(err)
	}

	for _, symbol := range referenceSymbols {
		if _, ok := response.Symbols[symbol]; !ok {
			t.Errorf("Reference symbol %s is absent in response %+v", symbol, response.Symbols)
		}
	}
}

func TestGetLatest(t *testing.T) {
	client, _ := getClient()

	response, err := client.GetLatest([]string{})
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, client.Base, response.Base)
	assert.True(t, regexp.MustCompile(`^\d+$`).MatchString(strconv.Itoa(response.Timestamp)))
	_, err = time.Parse("2006-01-02", response.Date)
	assert.Nil(t, err)
	assert.True(t, len(response.Rates) > 0)
}

func TestGetLatest_WithSymbols(t *testing.T) {
	client, _ := getClient()

	response, err := client.GetLatest(referenceSymbols)
	if err != nil {
		t.Error(err)
	}

	assert.True(t, len(response.Rates) == len(referenceSymbols))
}

func TestGetHistorical(t *testing.T) {
	client, _ := getClient()

	yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	response, err := client.GetHistorical(yesterday, []string{})
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, yesterday, response.Date)
}

func TestGetHistorical_WithSymbols(t *testing.T) {
	client, _ := getClient()

	yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	response, err := client.GetHistorical(yesterday, referenceSymbols)
	if err != nil {
		t.Error(err)
	}

	assert.True(t, len(response.Rates) == len(referenceSymbols))
}

func TestGetConversion(t *testing.T) {
	client, _ := getClient()

	amount := float64(100)
	date := time.Now().AddDate(0, 0, -1).Format("2006-01-02")

	_, err := client.GetConversion(referenceSymbols[0], referenceSymbols[1], amount, date)
	assert.EqualError(t, err, "gofixerio error: fixer.io request error with info: Access Restricted - Your current Subscription Plan does not support this API Function.")
}

func TestGetTimeseries(t *testing.T) {
	client, _ := getClient()

	beforeYesterday := time.Now().AddDate(0, 0, -2).Format("2006-01-02")
	yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")

	_, err := client.GetTimeseries(beforeYesterday, yesterday)
	assert.EqualError(t, err, "gofixerio error: fixer.io request error with info: Access Restricted - Your current Subscription Plan does not support this API Function.")
}

func TestGetFluctuation(t *testing.T) {
	client, _ := getClient()

	beforeYesterday := time.Now().AddDate(0, 0, -2).Format("2006-01-02")
	yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")

	_, err := client.GetFluctuation(beforeYesterday, yesterday)
	assert.EqualError(t, err, "gofixerio error: fixer.io request error with info: Access Restricted - Your current Subscription Plan does not support this API Function.")
}

func ExampleNew() {
	fixerio, _ := gofixerio.New("your API key here", "EUR", false)

	latestRates, err := fixerio.GetLatest([]string{})
	if err != nil {
		log.Panic(err)
	}

	log.Printf("%+v", latestRates)
}
