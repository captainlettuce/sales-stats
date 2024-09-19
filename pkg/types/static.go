package types

import (
	"encoding/json"
	"errors"
	"fmt"
	"slices"
)

type AggregationMethodType string

const AggregationMethodAvg AggregationMethodType = "avg"
const AggregationMethodMax AggregationMethodType = "max"
const AggregationMethodMin AggregationMethodType = "min"
const AggregationMethodSum AggregationMethodType = "sum"

func AggregationMethod(s string) (AggregationMethodType, error) {
	if slices.Contains([]AggregationMethodType{
		AggregationMethodAvg,
		AggregationMethodMax,
		AggregationMethodMin,
		AggregationMethodSum,
	}, AggregationMethodType(s)) {
		return AggregationMethodType(s), nil
	}
	return "", fmt.Errorf("invalid aggregation method '%s'", s)
}

func (am *AggregationMethodType) UnmarshalJSON(bytes []byte) error {
	if am == nil {
		return errors.New("trying to unmarshal AggregationMethodType into nil pointer")
	}

	var s string
	err := json.Unmarshal(bytes, &s)
	if err != nil {
		return err
	}

	*am, err = AggregationMethod(s)
	return err
}

type OrderStatusType string

const OrderStatusReturned OrderStatusType = "RETURNED"
const OrderStatusCompleted OrderStatusType = "COMPLETED"
const OrderStatusCancelled OrderStatusType = "CANCELLED"

func OrderStatus(s string) (OrderStatusType, error) {
	if slices.Contains([]OrderStatusType{
		OrderStatusCancelled,
		OrderStatusCompleted,
		OrderStatusReturned,
	}, OrderStatusType(s)) {
		return OrderStatusType(s), nil
	}
	return "", fmt.Errorf("invalid order status '%s'", s)
}

func (os *OrderStatusType) UnmarshalJSON(bytes []byte) error {
	if os == nil {
		return errors.New("trying to unmarshal OrderStatusType into nil pointer")
	}
	var err error
	var s string

	err = json.Unmarshal(bytes, &s)
	if err != nil {
		return err
	}

	*os, err = OrderStatus(s)
	return err
}

var (
	ErrInvalidArgument = errors.New("invalid argument")
	ErrUnexpectedError = errors.New("unexpected error")
)
