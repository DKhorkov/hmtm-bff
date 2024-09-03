package summer_test_test

import (
	"hmtm_bff/internal/summer"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSummerDouble(test *testing.T) {
	assert.Equal(test, 4, summer.Double(2), "Error with doubling")
}
