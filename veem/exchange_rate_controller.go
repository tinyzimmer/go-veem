package veem

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"
)

// ExchangeRateController is the interface for interacting with Veem exchange rates.
type ExchangeRateController interface {
	// Submits a request to generate an exchange rate quote
	CreateQuote(quote *QuoteRequest) (*Quote, error)
	// Submits a request to generate multiple exchange rate quotes
	CreateMultipleQuotes(quotes []*QuoteRequest) (*BatchQuoteResponse, error)
}

type exchangeRateController struct{ *client }

// Quote represents a quote for an exchange between two currencies.
type Quote struct {
	// The ID for the quote.
	ID string `json:"id"`
	// Time the quote expires, only valid on retrieval.
	Expiry time.Time `json:"expiry"`
	// The source amount.
	FromAmount float64 `json:"fromAmount"`
	// The target amount
	ToAmount float64 `json:"toAmount"`
	// The source currency
	FromCurrency string `json:"fromCurrency"`
	// The target currency
	ToCurrency string `json:"toCurrency"`
	// The rate of the exchange
	Rate float64 `json:"rate"`
}

// QuoteRequest represents a request for a quote.
type QuoteRequest struct {
	// The source amount. Either this or ToAmount can be specified.
	// The other is calculated.
	FromAmount float64 `json:"fromAmount,omitempty"`
	// The target amount. Either this or FromAmount can be specified.
	// The other is calculated.
	ToAmount float64 `json:"toAmount,omitempty"`
	// The source currency
	FromCurrency string `json:"fromCurrency"`
	// The target currency
	ToCurrency string `json:"toCurrency"`
	// The destination country to ensure Veem can support the transfer.
	ToCountry string `json:"toCountry"`
	// The email of recipient to get discounted rate
	RecipientAccountEmail string `json:"recipientAccountEmail,omitempty"`
}

// BatchQuoteResponse is a response to a request for a batch of quotes.
type BatchQuoteResponse struct {
	Quotes   []*Quote        `json:"success"`
	Failures []*QuoteFailure `json:"failure"`
}

type QuoteFailure struct {
	BatchItemID string `json:"batchItemId"`
	ErrorCode   string `json:"errorCode"`
}

func (e *exchangeRateController) CreateQuote(quote *QuoteRequest) (*Quote, error) {
	payload, err := json.Marshal(quote)
	if err != nil {
		return nil, err
	}
	req, err := e.newRequest(http.MethodPost, "veem/v1.1/exchangerates/quotes", bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}
	out := &Quote{}
	return out, e.doIntoWithAuth(req, out)
}

func (e *exchangeRateController) CreateMultipleQuotes(quotes []*QuoteRequest) (*BatchQuoteResponse, error) {
	payload, err := json.Marshal(quotes)
	if err != nil {
		return nil, err
	}
	req, err := e.newRequest(http.MethodPost, "veem/v1.1/exchangerates/quotes/batch", bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}
	out := &BatchQuoteResponse{}
	return out, e.doIntoWithAuth(req, out)
}
