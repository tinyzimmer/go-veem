package veem

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// InvoiceController is the interface for interacting with Veem invoices.
type InvoiceController interface {
	// Post an invoice and send to a receiver
	Create(inv *Invoice) (*Invoice, error)
	// Retrieve an invoice
	Get(id int64) (*Invoice, error)
	// Cancel an invoice
	Cancel(id int64) (*Invoice, error)
}

type Invoice struct {
	Payer                *Entity       `json:"payer"`
	Amount               *Amount       `json:"amount"`
	Attachments          []*Attachment `json:"attachments,omitempty"`
	CCEmails             []string      `json:"ccEmails,omitempty"`
	DueDate              *time.Time    `json:"dueDate,omitempty"`
	ExchangeRateQuoteId  string        `json:"exchangeRateQuoteId,omitempty"`
	ExternalInvoiceRefId string        `json:"externalInvoiceRefId,omitempty"`
	Notes                string        `json:"notes,omitempty"`
	PurposeOfPayment     string        `json:"purposeOfPayment,omitempty"`

	// Populated on retrieval
	ID          int64      `json:"id,omitempty"`
	Status      string     `json:"status,omitempty"`
	TimeCreated *time.Time `json:"timeCreated,omitempty"`
	ClaimLink   string     `json:"claimLink,omitempty"`
}

type invoiceController struct{ *client }

func (i *invoiceController) Create(inv *Invoice) (*Invoice, error) {
	payload, err := json.Marshal(inv)
	if err != nil {
		return nil, err
	}
	req, err := i.newRequest(http.MethodPost, "veem/v1.1/invoices", bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}
	out := &Invoice{}
	return out, i.doIntoWithAuth(req, out)
}

func (i *invoiceController) Get(id int64) (*Invoice, error) {
	req, err := i.newRequest(http.MethodGet, fmt.Sprintf("veem/v1.1/invoices/%d", id), nil)
	if err != nil {
		return nil, err
	}
	out := &Invoice{}
	return out, i.doIntoWithAuth(req, out)
}

func (i *invoiceController) Cancel(id int64) (*Invoice, error) {
	req, err := i.newRequest(http.MethodPost, fmt.Sprintf("veem/v1.1/invoices/%d/cancel", id), nil)
	if err != nil {
		return nil, err
	}
	out := &Invoice{}
	return out, i.doIntoWithAuth(req, out)
}
