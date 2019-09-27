package main

import (
	"testing"

	"github.com/blend/go-sdk/assert"
)

func Test_createWeights(t *testing.T) {
	assert := assert.New(t)

	weights := createWeights(katakana)
	assert.Len(weights, len(katakana))
	for key, value := range weights {
		assert.NotEmpty(katakana[key])
		assert.Equal(1.0, value)
	}
}

func Test_increaseWeight(t *testing.T) {
	assert := assert.New(t)

	weights := createWeights(katakana)

	increaseWeight(weights, "ヂ")
	assert.Equal(2.0, weights["ヂ"])
}

func Test_decreaseWeight(t *testing.T) {
	assert := assert.New(t)

	weights := createWeights(katakana)
	decreaseWeight(weights, "ヂ")
	assert.Equal(0.5, weights["ヂ"])
	decreaseWeight(weights, "ヂ")
	assert.Equal(0.25, weights["ヂ"])
	decreaseWeight(weights, "ヂ")
	assert.Equal(0.125, weights["ヂ"])
	decreaseWeight(weights, "ヂ")
	assert.Equal(0.125, weights["ヂ"])
}

func Test_random(t *testing.T) {
	assert := assert.New(t)

	weights := map[string]float64{
		"foo": 128.0,
		"bar": 16.0,
		"baz": 2.0,
		"boo": 0.25,
	}
	values := map[string]string{
		"foo": "alpha",
		"bar": "beta",
		"baz": "charlie",
		"boo": "delta",
	}

	counts := make(map[string]int)

	for x := 0; x < 1024; x++ {
		key, value := random(weights, values)
		assert.Equal(values[key], value)
		counts[key] = counts[key] + 1
	}

	assert.True(counts["foo"] > counts["bar"])
	assert.True(counts["bar"] > counts["baz"])
	assert.True(counts["baz"] > counts["boo"])
}

func Test_applyLimit(t *testing.T) {
	assert := assert.New(t)

	values := map[string]string{
		"foo": "f",
		"bar": "b",
		"moo": "m",
		"war": "w",
		"car": "c",
	}

	values = applyLimit(values, 3)
	assert.Len(values, 3)
}
