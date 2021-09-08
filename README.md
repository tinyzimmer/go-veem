# go-veem

A Go client for the [Veem](https://www.veem.com/) REST API

[![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat-rounded)](https://pkg.go.dev/github.com/tinyzimmer/go-veem)

## Quickstart

### Working with contacts

```go
package main

import (
    "fmt"

    "github.com/tinyzimmer/go-veem/veem"
)

func main() {
    // It is recommended to create and use a sandbox account
    // when first trying out the API.
    client, err := veem.New(&veem.ClientOptions{
        UseSandbox:   true,
        ClientID:     "TINYZIMMER-abcdefgh",
        ClientSecret: "superdupersecret",
    })
    if err != nil {
        panic(err)
    }

    // See the docs for available filters that can be passed to a List
    res, err := client.Contacts().List()
    if err != nil {
        panic(err)
    }
    for _, contact := range res.Contacts {
        fmt.Printf("%+v\n", contact)
    }

    // Create a new contact
    contact, err := client.Contacts().Create(&veem.ContactFull{
        Contact: &veem.Contact{
            ID:             0,
            BusinessName:   "Test Client",
            FirstName:      "Tiny",
            LastName:       "Zimmer",
            Email:          "tiny@zimmer.co",
            ISOCountryCode: "US",
            PhoneDialCode:  "+1",
            PhoneNumber:    "6785555555",
        },
    })

    if err != nil {
        panic(err)
    }

    fmt.Printf("%+v\n", contact)
}
```

### Working with invoices

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
        ClientSecret: "superdupersecret",
    })
    if err != nil {
        panic(err)
    }

    // Ensure a contact for the invoice - you can also get or list
    // contacts to retrieve their details.
    contact, err := client.Contacts().Create(&veem.ContactFull{
        Contact: &veem.Contact{
            ID:             0,
            BusinessName:   "Test Client",
            FirstName:      "Tiny",
            LastName:       "Zimmer",
            Email:          "tiny@zimmer.co",
            ISOCountryCode: "US",
            PhoneDialCode:  "+1",
            PhoneNumber:    "6785555555",
        },
    })

    // Create an attachment to an invoice
    attachment, err := client.Attachments().Upload("./INV-0001.pdf")
    if err != nil {
        panic(err)
    }

    // Download the attachment (obviously not required)
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
    attachment.Type = veem.ExternalInvoiceAttachment
    _, err = client.Invoices().Create(&veem.Invoice{
        Payer: contact.ToEntity(veem.ContactBusiness),
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