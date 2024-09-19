package types

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

type Date struct {
	time.Time
}

func (d *Date) MarshalJSON() ([]byte, error) {
	if d == nil {
		return nil, nil
	}
	return json.Marshal(d.Format(time.DateOnly))
}

func (d *Date) UnmarshalJSON(data []byte) error {
	// Remove leading and trailing `"`
	t, err := time.Parse(time.DateOnly, string(data)[1:len(data)-1])
	if err != nil {
		return err
	}

	*d = Date{t}

	return nil
}

func (d *Date) Value() (driver.Value, error) {
	return driver.Value(d.Format(time.DateOnly)), nil
}

func (d *Date) Scan(src interface{}) error {
	if src == nil {
		return nil
	}
	switch src := src.(type) {
	case time.Time:
		*d = Date{src}
		return nil
	}

	return fmt.Errorf("cannot convert %T to Date", src)
}
