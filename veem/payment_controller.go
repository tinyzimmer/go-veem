package veem

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

// PaymentController is the interface for interacting with Veem payments.
type PaymentController interface {
	// Get a payment by ID
	Get(id int64) (*Payment, error)
	// Get payments for this account with filters
	List(filters ...Filter) (*ListPaymentsResponse, error)
	// Create a new payment
	Create(payment *DraftPayment) (*Payment, error)
	// Create a batch of payments
	CreateBatch(payments []*DraftPayment, includeItems bool) (*BatchOperation, error)
	// Get the status of a batch operation
	GetBatch(batchID int64, includeItems bool) (*BatchOperation, error)
	// Approve a payment
	Approve(id int64) (*Payment, error)
	// Cancel a payment
	Cancel(id int64) (*Payment, error)
}

type Payment struct {
	Attachments            []*Attachment           `json:"attachments,omitempty"`
	BatchItemID            int64                   `json:"batchItemId,omitempty"`
	ClaimLink              string                  `json:"claimLink,omitempty"`
	CCEmails               []string                `json:"ccEmails,omitempty"`
	DueDate                time.Time               `json:"dueDate,omitempty"`
	ExchangeRateQuoteId    string                  `json:"exchangeRateQuoteId,omitempty"`
	ExternalInvoiceRefId   string                  `json:"externalInvoiceRefId,omitempty"`
	ID                     int64                   `json:"id,omitempty"`
	Notes                  string                  `json:"notes,omitempty"`
	Payee                  *Entity                 `json:"payee"`
	PayeeAmount            *Amount                 `json:"payeeAmount"`
	PaymentAction          string                  `json:"paymentAction,omitempty"`
	PaymentApproval        *PaymentApproval        `json:"paymentApproval,omitempty"`
	PaymentApprovalRequest *PaymentApprovalRequest `json:"paymentApprovalRequest,omitempty"`
	PurposeOfPayment       string                  `json:"purposeOfPayment,omitempty"`
	PushPaymentInfo        *PushPaymentInfo        `json:"pushPaymentInfo,omitempty"`
	Status                 string                  `json:"status"`
	TimeCreated            time.Time               `json:"timeCreated"`
	TimeUpdated            time.Time               `json:"timeUpdated"`
}

type PaymentApproval struct {
	ApprovalStatus         string          `json:"approvalStatus"`
	ApproverNumber         int64           `json:"approverNumber"`
	ApproverNumberRequired int64           `json:"approverNumberRequired"`
	UserApprovals          []*UserApproval `json:"userApprovalList"`
}

type UserApproval struct {
	ApprovalStatus string `json:"approvalStatus"`
	Email          string `json:"email"`
	FirstName      string `json:"firstName"`
	MiddleName     string `json:"middleName"`
	LastName       string `json:"lastName"`
}

type PaymentApprovalRequest struct {
	AccountID int64 `json:"accountId"`
	UserID    int64 `json:"userId"`
}

type PushPaymentInfo struct {
	Amount          *Amount `json:"amount"`
	PushPaymentInfo string  `json:"pushPaymentInfo"`
	Reference       string  `json:"reference"`
}

type DraftPayment struct {
	Amount               *Amount       `json:"amount"`
	ApproveAutomatically bool          `json:"approveAutomatically,omitempty"`
	Attachments          []*Attachment `json:"attachments,omitempty"`
	CCEmails             []string      `json:"ccEmails,omitempty"`
	DueDate              *time.Time    `json:"dueDate,omitempty"`
	ExchangeRateQuoteId  string        `json:"exchangeRateQuoteId,omitempty"`
	ExternalInvoiceRefId string        `json:"externalInvoiceRefId,omitempty"`
	Notes                string        `json:"notes,omitempty"`
	Payee                *Entity       `json:"payee"`
	PurposeOfPayment     string        `json:"purposeOfPayment,omitempty"`
}

type paymentControler struct{ *client }

type ListPaymentsResponse struct {
	Payments []*Payment `json:"content"`

	First            bool  `json:"first"`
	Last             bool  `json:"last"`
	NumberOfElements int   `json:"numberOfElements"`
	TotalElements    int   `json:"totalElements"`
	PageNumber       int32 `json:"number"`
	PageSize         int32 `json:"size"`
	TotalPages       int   `json:"totalPages"`

	controller *paymentControler
	filters    []Filter
}

func (p *paymentControler) Get(id int64) (*Payment, error) {
	ep := fmt.Sprintf("veem/v1.1/payments/%d", id)
	req, err := p.newRequest(http.MethodGet, ep, nil)
	if err != nil {
		return nil, err
	}
	payment := &Payment{}
	return payment, p.doIntoWithAuth(req, payment)
}

func (p *paymentControler) List(filters ...Filter) (*ListPaymentsResponse, error) {
	ep := "veem/v1.1/payments"
	if len(filters) > 0 {
		vals := &url.Values{}
		for _, f := range filters {
			f(vals)
		}
		ep = fmt.Sprintf("%s?%s", ep, vals.Encode())
	} else {
		filters = make([]Filter, 0)
	}
	req, err := p.newRequest(http.MethodGet, ep, nil)
	if err != nil {
		return nil, err
	}
	out := &ListPaymentsResponse{controller: p, filters: filters}
	return out, p.doIntoWithAuth(req, out)
}

func (g *ListPaymentsResponse) Next() (*ListPaymentsResponse, error) {
	if g.Last {
		return nil, errors.New("no more payment pages left")
	}
	return g.controller.List(
		append(g.filters, WithPageNumber(g.PageNumber+1), WithPageSize(g.PageSize))...,
	)
}

func (p *paymentControler) Create(payment *DraftPayment) (*Payment, error) {
	payload, err := json.Marshal(payment)
	if err != nil {
		return nil, err
	}
	req, err := p.newRequest(http.MethodPost, "veem/v1.1/payments", bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}
	out := &Payment{}
	return out, p.doIntoWithAuth(req, out)
}

func (p *paymentControler) CreateBatch(payments []*DraftPayment, includeItems bool) (*BatchOperation, error) {
	payload, err := json.Marshal(payments)
	if err != nil {
		return nil, err
	}
	ep := fmt.Sprintf("veem/v1.1/payments/batch?includeItems=%t", includeItems)
	req, err := p.newRequest(http.MethodPost, ep, bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}
	out := &BatchOperation{}
	return out, p.doIntoWithAuth(req, out)
}

func (p *paymentControler) GetBatch(batchID int64, includeItems bool) (*BatchOperation, error) {
	ep := fmt.Sprintf("veem/v1.1/payments/batch/%d?includeItems=%t", batchID, includeItems)
	req, err := p.newRequest(http.MethodGet, ep, nil)
	if err != nil {
		return nil, err
	}
	out := &BatchOperation{}
	return out, p.doIntoWithAuth(req, out)
}

func (p *paymentControler) Approve(id int64) (*Payment, error) {
	ep := fmt.Sprintf("veem/v1.1/payments/%d/approve", id)
	req, err := p.newRequest(http.MethodPost, ep, nil)
	if err != nil {
		return nil, err
	}
	payment := &Payment{}
	return payment, p.doIntoWithAuth(req, payment)
}

func (p *paymentControler) Cancel(id int64) (*Payment, error) {
	ep := fmt.Sprintf("veem/v1.1/payments/%d/cancel", id)
	req, err := p.newRequest(http.MethodPost, ep, nil)
	if err != nil {
		return nil, err
	}
	payment := &Payment{}
	return payment, p.doIntoWithAuth(req, payment)
}
