package chglog

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDotGet(t *testing.T) {
	assert := assert.New(t)
	now := time.Now()

	type Nest struct {
		Str  string
		Int  int
		Time time.Time
	}

	type Sample struct {
		Str  string
		Int  int
		Date time.Time
		Nest Nest
	}

	sample := Sample{
		Str:  "sample_string",
		Int:  12,
		Date: now,
		Nest: Nest{
			Str:  "nest_string",
			Int:  34,
			Time: now,
		},
	}

	var val interface{}
	var ok bool

	// .Str
	val, ok = dotGet(&sample, "Str")
	assert.True(ok)
	assert.Equal(val, "sample_string")

	// Lowercase
	val, ok = dotGet(&sample, "str")
	assert.True(ok)
	assert.Equal(val, "sample_string")

	// Int
	val, ok = dotGet(&sample, "Int")
	assert.True(ok)
	assert.Equal(val, 12)

	// Time
	val, ok = dotGet(&sample, "Date")
	assert.True(ok)
	assert.Equal(val, now)

	// Nest
	val, ok = dotGet(&sample, "Nest.Str")
	assert.True(ok)
	assert.Equal(val, "nest_string")

	val, ok = dotGet(&sample, "Nest.Int")
	assert.True(ok)
	assert.Equal(val, 34)

	val, ok = dotGet(&sample, "Nest.Time")
	assert.True(ok)
	assert.Equal(val, now)

	val, ok = dotGet(&sample, "nest.int")
	assert.True(ok)
	assert.Equal(val, 34)

	// Notfound
	val, ok = dotGet(&sample, "not.found")
	assert.False(ok)
	assert.Nil(val)
}

func TestCompare(t *testing.T) {
	assert := assert.New(t)

	type sample struct {
		a        interface{}
		op       string
		b        interface{}
		expected bool
	}

	table := []sample{
		{0, "<", 1, true},
		{0, ">", 1, false},
		{1, ">", 0, true},
		{1, "<", 0, false},
		{"a", "<", "b", true},
		{"a", ">", "b", false},
		{time.Unix(1518018017, 0), "<", time.Unix(1518018043, 0), true},
		{time.Unix(1518018017, 0), ">", time.Unix(1518018043, 0), false},
	}

	for _, sa := range table {
		actual, err := compare(sa.a, sa.op, sa.b)
		assert.Nil(err)
		assert.Equal(sa.expected, actual)
	}
}
