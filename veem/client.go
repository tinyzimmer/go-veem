package veem

import (
	"net/http"
	"net/url"
)

type Client interface {
	Meta() MetaController
	Attachments() AttachmentController
	Contacts() ContactController
	Customers() CustomerController
	ExchangeRates() ExchangeRateController
	Invoices() InvoiceController
	Payments() PaymentController
}

var sandboxURL = mustParseURL("https://sandbox-api.veem.com")
var liveURL = mustParseURL("https://api.veem.com")

type ClientOptions struct {
	// Use the sandbox API.
	UseSandbox bool
	// ClientID and ClientSecret for authenticating with Veem
	ClientID, ClientSecret string
}

func New(opts *ClientOptions) (Client, error) {
	url := liveURL
	if opts.UseSandbox {
		url = sandboxURL
	}
	c := &client{opts: opts, apiURL: url, client: &http.Client{}}
	var err error
	c.token, err = c.getAccessToken()
	if err != nil {
		return nil, err
	}
	return c, nil
}

type client struct {
	opts   *ClientOptions
	apiURL *url.URL
	client *http.Client
	token  *AccessTokenResponse
}

func (c *client) Meta() MetaController                  { return &metaController{c} }
func (c *client) Attachments() AttachmentController     { return &attachmentController{c} }
func (c *client) Contacts() ContactController           { return &contactController{c} }
func (c *client) Customers() CustomerController         { return &customerController{c} }
func (c *client) ExchangeRates() ExchangeRateController { return &exchangeRateController{c} }
func (c *client) Invoices() InvoiceController           { return &invoiceController{c} }
func (c *client) Payments() PaymentController           { return &paymentControler{c} }
