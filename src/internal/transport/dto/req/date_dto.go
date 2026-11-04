package req

import (
	"encoding/json"
	"time"
)

type DateDto time.Time

func (d *DateDto) UnmarshalJSON(data []byte) error {
	var value string
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	date, err := time.Parse("2006-01-02", value)
	if err != nil {
		return err
	}
	*d = DateDto(date)
	return nil
}

func (d DateDto) Time() time.Time {
	return time.Time(d)
}
