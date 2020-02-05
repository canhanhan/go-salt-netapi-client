package cherrypy

import (
	"strconv"
	"strings"
	"time"
)

type saltUnixTime struct {
	time.Time
}

func (t *saltUnixTime) UnmarshalJSON(input []byte) error {
	s, err := strconv.ParseFloat(string(input), 64)
	if err != nil {
		return err
	}

	m := int64(s)
	n := int64((s - float64(m)) * 1000000000)
	t.Time = time.Unix(m, n)
	return nil
}

type saltTime struct {
	time.Time
}

func (t *saltTime) UnmarshalJSON(input []byte) error {
	s := string(input)
	s = strings.Trim(s, "\"")
	v, err := time.Parse("2006, Jan 02 15:04:05.000000", s)
	if err != nil {
		return err
	}

	t.Time = v
	return nil
}
