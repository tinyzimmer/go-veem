package veem

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

// ContactController is the interface for interacting with Veem contacts.
type ContactController interface {
	// Get an account contact by ID
	Get(id int64) (*Contact, error)
	// Get a page of account contacts by email address, first,last name, batchId, and business name
	List(filters ...Filter) (*ListContactsResponse, error)
	// Create a contact
	Create(contact *ContactFull) (*Contact, error)
	// Create a batch of contacts
	CreateBatch(contacts []*ContactFull, includeItems bool) (*BatchOperation, error)
	// Get the status of a batch operation
	GetBatch(batchID int64, includeItems bool) (*BatchOperation, error)
}

type ContactType string

const (
	ContactIncomplete ContactType = "Incomplete"
	ContactBusiness   ContactType = "Business"
	ContactPersonal   ContactType = "Personal"
)

type Contact struct {
	ID               int64  `json:"id,omitempty"`
	BusinessName     string `json:"businessName,omitempty"`
	FirstName        string `json:"firstName"`
	LastName         string `json:"lastName"`
	Email            string `json:"email"`
	ISOCountryCode   string `json:"isoCountryCode"`
	PhoneDialCode    string `json:"dialCode"`
	PhoneNumber      string `json:"phoneNumber"`
	BatchItemID      int64  `json:"batchItemId"`
	ContactAccountID int64  `json:"contactAccountId"`
}

func (c *Contact) ToEntity(t ContactType) *Entity {
	return &Entity{
		BusinessName: c.BusinessName,
		CountryCode:  c.ISOCountryCode,
		Email:        c.Email,
		FirstName:    c.FirstName,
		LastName:     c.LastName,
		Type:         t,
		Phone:        c.PhoneDialCode + c.PhoneNumber,
	}
}

func (c *Contact) toMap() map[string]interface{} {
	out := map[string]interface{}{
		"firstName":      c.FirstName,
		"lastName":       c.LastName,
		"email":          c.Email,
		"isoCountryCode": c.ISOCountryCode,
		"phoneNumber":    c.PhoneNumber,
		"phoneDialCode":  c.PhoneDialCode, // only param that is different on POST/GET
	}
	if c.BusinessName != "" {
		out["businessName"] = c.BusinessName
	}
	if c.ID != 0 {
		out["id"] = c.ID
	}
	if c.BatchItemID != 0 {
		out["batchItemId"] = c.BatchItemID
	}
	if c.ContactAccountID != 0 {
		out["contactAccountId"] = c.ContactAccountID
	}
	return out
}

type ContactFull struct {
	*Contact `json:",inline"`

	Type               ContactType  `json:"type,omitempty"`
	ExternalBusinessID int64        `json:"externalBusinessID,omitempty"`
	BusinessAddress    *Address     `json:"businessAddress,omitempty"`
	BankAccount        *BankAccount `json:"bankAccount,omitempty"`
}

func (c *ContactFull) MarshalJSON() ([]byte, error) {
	out := c.toMap()
	if c.Type != "" {
		out["type"] = c.Type
	}
	if c.ExternalBusinessID != 0 {
		out["externalBusinessID"] = c.ExternalBusinessID
	}
	if c.BusinessAddress != nil {
		out["businessAddress"] = c.BusinessAddress
	}
	if c.BankAccount != nil {
		out["bankAccount"] = c.BankAccount
	}
	return json.Marshal(out)
}

type Address struct {
	Line1         string `json:"line1"`
	Line2         string `json:"line2"`
	City          string `json:"city"`
	StateProvince string `json:"stateProvince"`
	PostalCode    string `json:"postalCode"`
}

type BankAccount struct {
	AccountNumber         string   `json:"bankAccountNumber,omitempty"`
	RoutingNumber         string   `json:"routingNumber,omitempty"`
	BankName              string   `json:"bankName,omitempty"`
	BankAddress           *Address `json:"bankAddress,omitempty"`
	BankCNaps             string   `json:"bankCnaps,omitempty"`
	BankCode              string   `json:"bankCode,omitempty"`
	BankIFSCBranchCode    string   `json:"bankIfscBranchCode,omitempty"`
	BankInstitutionNumber string   `json:"bankInstitutionNumber,omitempty"`
	BeneficiaryName       string   `json:"beneficiaryName,omitempty"`
	BranchCode            string   `json:"branchCode,omitempty"`
	BSBBankCode           string   `json:"bsbBankCode,omitempty"`
	CLABE                 string   `json:"clabe,omitempty"`
	CurrencyCode          string   `json:"currencyCode,omitempty"`
	IBAN                  string   `json:"iban,omitempty"`
	ISOCountryCode        string   `json:"isoCountryCode,omitempty"`
	SortCode              string   `json:"sortCode,omitempty"`
	SwiftBIC              string   `json:"swiftBic,omitempty"`
	TransitCode           string   `json:"transitCode,omitempty"`
}

type contactController struct{ *client }

type ListContactsResponse struct {
	Contacts []*Contact `json:"content"`

	First            bool  `json:"first"`
	Last             bool  `json:"last"`
	NumberOfElements int   `json:"numberOfElements"`
	TotalElements    int   `json:"totalElements"`
	PageNumber       int32 `json:"number"`
	PageSize         int32 `json:"size"`
	TotalPages       int   `json:"totalPages"`

	controller *contactController
	filters    []Filter
}

func (c *contactController) Get(id int64) (*Contact, error) {
	ep := fmt.Sprintf("veem/v1.1/contacts/%d", id)
	req, err := c.newRequest(http.MethodGet, ep, nil)
	if err != nil {
		return nil, err
	}
	contact := &Contact{}
	return contact, c.doIntoWithAuth(req, contact)
}

func (c *contactController) List(filters ...Filter) (*ListContactsResponse, error) {
	ep := "veem/v1.1/contacts"
	if len(filters) > 0 {
		vals := &url.Values{}
		for _, f := range filters {
			f(vals)
		}
		ep = fmt.Sprintf("%s?%s", ep, vals.Encode())
	} else {
		filters = make([]Filter, 0)
	}
	req, err := c.newRequest(http.MethodGet, ep, nil)
	if err != nil {
		return nil, err
	}
	out := &ListContactsResponse{controller: c, filters: filters}
	return out, c.doIntoWithAuth(req, out)
}

func (g *ListContactsResponse) Next() (*ListContactsResponse, error) {
	if g.Last {
		return nil, errors.New("no more contact pages left")
	}
	return g.controller.List(
		append(g.filters, WithPageNumber(g.PageNumber+1), WithPageSize(g.PageSize))...,
	)
}

func (c *contactController) Create(contact *ContactFull) (*Contact, error) {
	payload, err := json.Marshal(contact)
	if err != nil {
		return nil, err
	}
	req, err := c.newRequest(http.MethodPost, "veem/v1.1/contacts", bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}
	out := &Contact{}
	return out, c.doIntoWithAuth(req, out)
}

func (c *contactController) CreateBatch(contacts []*ContactFull, includeItems bool) (*BatchOperation, error) {
	payload, err := json.Marshal(contacts)
	if err != nil {
		return nil, err
	}
	ep := fmt.Sprintf("veem/v1.1/contacts/batch?includeItems=%t", includeItems)
	req, err := c.newRequest(http.MethodPost, ep, bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}
	out := &BatchOperation{}
	return out, c.doIntoWithAuth(req, out)
}

func (c *contactController) GetBatch(batchID int64, includeItems bool) (*BatchOperation, error) {
	ep := fmt.Sprintf("veem/v1.1/contacts/batch/%d?includeItems=%t", batchID, includeItems)
	req, err := c.newRequest(http.MethodGet, ep, nil)
	if err != nil {
		return nil, err
	}
	out := &BatchOperation{}
	return out, c.doIntoWithAuth(req, out)
}
