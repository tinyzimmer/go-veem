# go-veem

A Go client for the [Veem](https://www.veem.com/) REST API

[![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat-rounded)](https://pkg.go.dev/github.com/tinyzimmer/go-veem)

## Quickstart

```go
package main

import (
	"github.com/tinyzimmer/go-veem/veem"
)

func main() {
    // It is recommended to create and use a sandbox account
    // when first trying out the API.
	client, err := veem.New(&veem.ClientOptions{
		UseSandbox:   true,
		ClientID:     "TINYZIMMER-abcdefgh",
		ClientSecret: "superdupersecrt",
	})
	if err != nil {
		panic(err)
	}

    // Retrieve the full country currency map
    countryMap, err := client.Meta().CountryCurrencyMap(true)
    if err != nil {
        panic(err)
    }
    for _, country := range countryMap {
        fmt.Println(country)
    }

    // Create an attachment to an invoice or payment
    attachment, err := client.Attachments().Upload("./INV-0001.pdf")
    if err != nil {
        panic(err)
    }

    // Download the attachment
    rdr, err := client.Attachments().Download(attachment.Name, attachment.ReferenceID)
    if err != nil {
        panic(err)
    }
    f, err := os.Create("downloaded.pdf")
    if err != nil {
        panic(err)
    }
    defer f.Close()
    io.Copy(f, rdr)

    // Create an invoice
    _, err := client.Invoices().Create(&veem.Invoice{
		Payer: &veem.Entity{
			BusinessName: "TinyZimmerTech LTD",
			CountryCode:  "US",
			Email:        "me@example.com",
			FirstName:    "Tiny",
			LastName:     "Zimmer",
			Type:         "Business",
			Phone:        "6785555555",
		},
		Amount: &veem.Amount{
			Currency: "USD",
			Number:   1000,
		},
		Attachments: []*veem.Attachment{attachment},
	})
	if err != nil {
		panic(err)
	}
}

```

See the Godoc for more information. More examples will come later.