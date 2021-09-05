package veem

import (
	"net/url"
	"strconv"
)

type Filter func(*url.Values)

func WithEmail(email string) Filter {
	return func(vals *url.Values) {
		vals.Add("email", email)
	}
}

func WithFirstName(name string) Filter {
	return func(vals *url.Values) {
		vals.Add("firstName", name)
	}
}

func WithLastName(name string) Filter {
	return func(vals *url.Values) {
		vals.Add("lastName", name)
	}
}

func WithBusinessName(name string) Filter {
	return func(vals *url.Values) {
		vals.Add("businessName", name)
	}
}

func WithBatchID(id int64) Filter {
	return func(vals *url.Values) {
		vals.Add("batchId", strconv.Itoa(int(id)))
	}
}

func WithBatchItemIDs(ids ...int64) Filter {
	return func(vals *url.Values) {
		for _, id := range ids {
			vals.Add("batchItemIds", strconv.Itoa(int(id)))
		}
	}
}

func WithPageNumber(num int32) Filter {
	return func(vals *url.Values) {
		vals.Add("pageNumber", strconv.Itoa(int(num)))
	}
}

func WithPageSize(size int32) Filter {
	return func(vals *url.Values) {
		vals.Add("pageSize", strconv.Itoa(int(size)))
	}
}

func WithPaymentIDs(ids ...int64) Filter {
	return func(vals *url.Values) {
		for _, id := range ids {
			vals.Add("paymentIds", strconv.Itoa(int(id)))
		}
	}
}

func WithStatuses(statuses ...string) Filter {
	return func(vals *url.Values) {
		for _, status := range statuses {
			vals.Add("status", status)
		}
	}
}

func WithSortTimeUpdatedAscending() Filter {
	return func(vals *url.Values) {
		vals.Add("sort", "timeUpdated:asc")
	}
}

func WithSortTimeUpdatedDescending() Filter {
	return func(vals *url.Values) {
		vals.Add("sort", "timeUpdated:desc")
	}
}
