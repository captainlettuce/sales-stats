package types

type AggregationResult struct {
	Avg int64 `json:"avg,omitempty" db:"avg"`
	Max int64 `json:"max,omitempty" db:"max"`
	Min int64 `json:"min,omitempty" db:"min"`
	Sum int64 `json:"sum,omitempty" db:"sum"`

	OrderDate Date `json:"orderDate,omitempty" db:"od"`

	Make      *string `json:"make,omitempty" db:"make"`
	Model     *string `json:"model,omitempty" db:"model"`
	Version   *string `json:"version,omitempty" db:"version"`
	Color     *string `json:"color,omitempty" db:"color"`
	ModelYear *int64  `json:"modelYear,omitempty" db:"model_year"`
}
