package veem

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

// CustomerController is the interface for interacting with Veem customers.
type CustomerController interface {
	// Search Veem Contacts
	Search(filters ...Filter) (*SearchCustomersResponse, error)
}

type Customer struct {
	ID             int64  `json:"id,omitempty"`
	Name           string `json:"name,omitempty"`
	FirstName      string `json:"firstName"`
	LastName       string `json:"lastName"`
	Email          string `json:"email"`
	ISOCountryCode string `json:"isoCountryCode"`
	IsContact      bool   `json:"isContact,omitempty"`
}

type SearchCustomersResponse struct {
	Customers []*Customer `json:"content"`

	First            bool  `json:"first"`
	Last             bool  `json:"last"`
	NumberOfElements int   `json:"numberOfElements"`
	TotalElements    int   `json:"totalElements"`
	PageNumber       int32 `json:"number"`
	PageSize         int32 `json:"size"`
	TotalPages       int   `json:"totalPages"`

	controller *customerController
	filters    []Filter
}

func (s *SearchCustomersResponse) Next() (*SearchCustomersResponse, error) {
	if s.Last {
		return nil, errors.New("no more customer pages left")
	}
	return s.controller.Search(
		append(s.filters, WithPageNumber(s.PageNumber+1), WithPageSize(s.PageSize))...,
	)
}

type customerController struct{ *client }

func (c *customerController) Search(filters ...Filter) (*SearchCustomersResponse, error) {
	ep := "veem/v1.1/customers"
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
	out := &SearchCustomersResponse{controller: c, filters: filters}
	return out, c.doIntoWithAuth(req, out)
}
