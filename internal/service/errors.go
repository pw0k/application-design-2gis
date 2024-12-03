package service

import (
	"application-design-test/internal/model"
	"errors"
	"fmt"
)

var ErrQuotaUnavailable = errors.New("quota unavailable")

type QuotaUnavailableError struct {
	Order *model.Order
}

func newQuotaUnavailableError(order *model.Order) *QuotaUnavailableError {
	return &QuotaUnavailableError{
		Order: order,
	}
}

func (e *QuotaUnavailableError) Error() string {
	return fmt.Sprintf("qota is unavailable for specific date, order %v", e.Order)
}

func (e *QuotaUnavailableError) Is(target error) bool {
	return target == ErrQuotaUnavailable
}
