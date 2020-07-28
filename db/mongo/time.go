package mongo

import (
	"encoding/json"
	"time"
)

type DateTime string

func (d DateTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.String())
}

func (d DateTime) String() string {
	return string(d)
}

func (d DateTime) Time() time.Time {
	t, _ := time.Parse("2006-01-02 15:04:05", d.String())
	return t
}

func NewDateTimeFromTime(t time.Time) DateTime {
	return DateTime(t.Format("2006-01-02 15:04:05"))
}
