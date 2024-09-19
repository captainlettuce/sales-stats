package types

type AggregationRequest struct {
	AggregationMethods []AggregationMethodType `json:"aggregationMethods"`

	Filter  FilterCriteria   `json:"filter"`
	GroupBy GroupingCriteria `json:"groupBy"`
}

type FilterCriteria struct {
	OrderDate *DateRange `json:"orderDate,omitempty"`

	OrderStatus []OrderStatusType `json:"orderStatus,omitempty"`
	Make        []string          `json:"make,omitempty"`
	Model       []string          `json:"model,omitempty"`
	Version     []string          `json:"version,omitempty"`
	Color       []string          `json:"color,omitempty"`

	// ToDo: Does this really make sense?
	Price *IntegerRange `json:"price,omitempty"`

	ModelYear         *IntegerRange `json:"modelYear,omitempty"`
	MileageKilometers *IntegerRange `json:"mileageKilometers,omitempty"`
}

type DateRange struct {
	Before *Date `json:"before,omitempty"`
	After  *Date `json:"after,omitempty"`
}

type IntegerRange struct {
	Over  *int64 `json:"over,omitempty"`
	Under *int64 `json:"under,omitempty"`
}

type GroupingCriteria struct {
	Make      bool `json:"make,omitempty" db:"make"`
	Model     bool `json:"model,omitempty" db:"model"`
	Version   bool `json:"version,omitempty" db:"version"`
	Color     bool `json:"color,omitempty" db:"color"`
	ModelYear bool `json:"modelYear,omitempty" db:"model_year"`
}
