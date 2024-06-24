package stringutil

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsValidAddressErr(t *testing.T) {
	type testCase struct {
		value    string
		positive bool
	}

	for _, tc := range []testCase{
		{
			value:    "0.0.0.0:80",
			positive: true,
		},
		{
			value:    "localhost:80",
			positive: true,
		},
		{
			value:    "service-1.dc1.ru:6060",
			positive: true,
		},
		{
			// cannot use underscores in hostnames (see https://stackoverflow.com/a/2183140/2361497)
			value:    "service_1.dc1.ru:6060",
			positive: false,
		},
		{
			value:    "123service:6060",
			positive: true,
		},
		{
			value:    "123service.dc1.ru:6060",
			positive: true,
		},
		{
			value:    "1:6060",
			positive: true,
		},
		{
			value:    "1.dc1.ru:6060",
			positive: true,
		},
		{
			value:    "192.168.0.255:6060",
			positive: true,
		},
		{
			value:    "192.168.0.300:6060",
			positive: false,
		},
		{
			value:    "1.dc1.499:6060",
			positive: false,
		},
		{
			value:    "1.dc1.255:6060",
			positive: false,
		},
	} {
		tc := tc
		t.Run(tc.value, func(t *testing.T) {
			err := IsValidAddressErr(tc.value)
			if tc.positive {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}
}
