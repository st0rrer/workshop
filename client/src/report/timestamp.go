package report

import (
	"fmt"
	"strconv"
	"time"
)

type timestamp time.Time

func (t *timestamp) MarshalJSON() ([]byte, error) {
	ts := time.Time(*t).UnixNano() / int64(time.Millisecond)

	return []byte(fmt.Sprint(ts)), nil
}

func (t *timestamp) UnmarshalJSON(data []byte) error {

	if string(data) == "null" {
		return nil
	}

	milliseconds, err := strconv.ParseInt(string(data), 0, 64)
	if err != nil {
		return err
	}

	*t = timestamp(time.Unix(0, milliseconds*int64(time.Millisecond)))
	return nil
}
