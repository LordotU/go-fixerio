package gofixerio

type symbols map[string]string
type rates map[string]float64
type ratesNested map[string]rates

type conversionQuery struct {
	From   string  `json:"from"`
	To     string  `json:"to"`
	Amount float64 `json:"amount"`
}
type conversionInfo struct {
	Timestamp int     `json:"timestamp"`
	Rate      float64 `json:"rate"`
}

// ResponseSuccess represents the success of the request
type ResponseSuccess struct {
	Success bool `json:"success"`
}

// ResponseSymbols represents GetSymbols response
type ResponseSymbols struct {
	ResponseSuccess
	Symbols symbols `json:"symbols"`
}

// ResponseLatest represents GetLatest response
type ResponseLatest struct {
	ResponseSuccess
	Timestamp int    `json:"timestamp"`
	Base      string `json:"base"`
	Date      string `json:"date"`
	Rates     rates  `json:"rates"`
}

// ResponseHistorical represents GetHistorical response
type ResponseHistorical struct {
	ResponseSuccess
	ResponseLatest
	Historical bool `json:"historical"`
}

// ResponseConversion represents GetConversion response
type ResponseConversion struct {
	ResponseSuccess
	Query      conversionQuery `json:"query"`
	Info       conversionInfo  `json:"info"`
	Historical bool            `json:"historical"`
	Date       string          `json:"date"`
	Result     float64         `json:"result"`
}

// ResponseTimeseries represents GetTimeseries response
type ResponseTimeseries struct {
	ResponseSuccess
	Timeseries bool        `json:"timeseries"`
	StartDate  string      `json:"start_date"`
	EndDate    string      `json:"end_date"`
	Base       string      `json:"base"`
	Rates      ratesNested `json:"rates"`
}

// ResponseFluctuation represents GetFluctuation response
type ResponseFluctuation struct {
	ResponseSuccess
	Fluctuation bool        `json:"fluctuation"`
	StartDate   string      `json:"start_date"`
	EndDate     string      `json:"end_date"`
	Base        string      `json:"base"`
	Rates       ratesNested `json:"rates"`
}

// ResponseErrorDetail represents details of errored response
type ResponseErrorDetail struct {
	Code int    `json:"code"`
	Type string `json:"type"`
	Info string `json:"info"`
}

// ResponseError represents errored response
type ResponseError struct {
	ResponseSuccess
	Error ResponseErrorDetail `json:"error"`
}
