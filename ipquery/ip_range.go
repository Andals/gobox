package ipquery

import (
	"errors"
)

const (
	IP_RANGE_FIELD_COUNT = 3
)

var ErrorIpRangeNotFound = errors.New("ip range not found")

type IpRange struct {
	Begin uint32
	End   uint32
	Data  []byte
}
