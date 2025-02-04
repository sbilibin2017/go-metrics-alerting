package key

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncode(t *testing.T) {
	ke := NewKeyEngine()

	tests := []struct {
		metricType string
		metricName string
		expected   string
	}{
		{"counter", "requests", "counter:requests"},
		{"gauge", "temperature", "gauge:temperature"},
		{"", "metric", ""},
		{"metric", "", ""},
	}

	for _, test := range tests {
		result := ke.Encode(test.metricType, test.metricName)
		assert.Equal(t, test.expected, result)
	}
}

func TestDecode(t *testing.T) {
	ke := NewKeyEngine()

	tests := []struct {
		key      string
		expType  string
		expName  string
		expError error
	}{
		{"counter:requests", "counter", "requests", nil},
		{"gauge:temperature", "gauge", "temperature", nil},
		{"invalidkey", "", "", ErrInvalidKeyFormat},
		{":missingtype", "", "", ErrInvalidKeyFormat},
		{"missingname:", "", "", ErrInvalidKeyFormat},
	}

	for _, test := range tests {
		mType, mName, err := ke.Decode(test.key)
		assert.Equal(t, test.expType, mType)
		assert.Equal(t, test.expName, mName)
		if test.expError != nil {
			assert.Error(t, err)
			assert.EqualError(t, err, test.expError.Error())
		} else {
			assert.NoError(t, err)
		}
	}
}
