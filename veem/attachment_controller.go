package veem

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

// AttachmentController is the interface for managing attachments
// to Invoices and Payments.
type AttachmentController interface {
	// Uploads an external attachment for a Payment or Invoice
	Upload(filename string) (*Attachment, error)
	// Downloads the referenced file
	Download(name, referenceID string) (io.ReadCloser, error)
}

type attachmentController struct{ *client }

type Attachment struct {
	Name        string `json:"name"`
	ReferenceID string `json:"referenceId"`
	Type        string `json:"type"`
}

func (a *attachmentController) Upload(filename string) (*Attachment, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	part, err := writer.CreateFormFile("file", filepath.Base(f.Name()))
	if err != nil {
		return nil, err
	}
	if _, err := io.Copy(part, f); err != nil {
		return nil, err
	}
	if err := writer.Close(); err != nil {
		return nil, err
	}
	req, err := a.newRequest(http.MethodPost, "veem/v1.1/attachments", &body)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", writer.FormDataContentType())
	req.Header.Add("X-REQUEST-ID", uuid.New().String())
	res := &Attachment{}
	return res, a.doIntoWithAuth(req, res)
}

func (a *attachmentController) Download(name, referenceID string) (io.ReadCloser, error) {
	ep := fmt.Sprintf("veem/v1.1/attachments?name=%s&referenceId=%s", name, referenceID)
	req, err := a.newRequest(http.MethodGet, ep, nil)
	if err != nil {
		return nil, err
	}
	return a.doWithAuth(req, "application/octet-stream")
}
