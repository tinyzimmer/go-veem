package veem

type BatchOperation struct {
	BatchID        int64        `json:"batchId"`
	BatchItems     []*BatchItem `json:"batchItems,omitempty"`
	HasErrors      bool         `json:"hasErrors"`
	ProcessedItems int64        `json:"processedItems"`
	Status         string       `json:"status"`
	TotalItems     int64        `json:"totalItems"`
}

type BatchItem struct {
	BatchItemID int64     `json:"batchItemId"`
	ErrorInfo   *APIError `json:"errorInfo,omitempty"`
	Status      string    `json:"status"`
}

type Entity struct {
	BusinessName string      `json:"businessName,omitempty"`
	CountryCode  string      `json:"countryCode"`
	Email        string      `json:"email"`
	FirstName    string      `json:"firstName"`
	LastName     string      `json:"lastName"`
	Type         ContactType `json:"type"`
	Phone        string      `json:"phone"`
}
