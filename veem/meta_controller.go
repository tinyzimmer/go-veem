package veem

import (
	"fmt"
	"net/http"
)

// MetaController is the interface for accessing veem metadata.
type MetaController interface {
	// Returns a list of countries supported and currencies for each
	CountryCurrencyMap(bankFields bool) ([]*CountryCurrentMap, error)
}

type metaController struct{ *client }

type CountryCurrentMap struct {
	BankFields                []string          `json:"bankFields"`
	Country                   string            `json:"country"`
	CountryName               string            `json:"countryName"`
	InvoiceAttachmentRequired bool              `json:"invoiceAttachmentRequired"`
	PurposeOfPaymentRequired  bool              `json:"purposeOfPaymentRequired"`
	PurposeOfPaymentInfo      []*PaymentPurpose `json:"purposeOfPaymentInfo"`
	ReceivingCurrencies       []string          `json:"receivingCurrencies"`
	SendingCurrencies         []string          `json:"sendingCurrencies"`
}

type PaymentPurpose struct {
	CountryCode string `json:"countryCode"`
	Description string `json:"description"`
	Industry    string `json:"industry"`
	SubIndustry string `json:"subindustry"`
	PurposeCode string `json:"purposeCode"`
}

func (m *metaController) CountryCurrencyMap(bankFields bool) ([]*CountryCurrentMap, error) {
	ep := fmt.Sprintf("veem/public/v1.1/country-currency-map?bankFields=%t", bankFields)
	req, err := m.newRequest(http.MethodGet, ep, nil)
	if err != nil {
		return nil, err
	}
	out := make([]*CountryCurrentMap, 0)
	return out, m.doIntoWithAuth(req, &out)
}
