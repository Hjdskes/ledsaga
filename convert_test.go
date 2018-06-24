package sladdfri

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKelvinToMired(t *testing.T) {
	assert := assert.New(t)
	assert.Equal(454, KelvinToMired(2100))
	assert.Equal(454, KelvinToMired(2200))
	assert.Equal(345, KelvinToMired(2900))
	assert.Equal(250, KelvinToMired(4000))
	assert.Equal(250, KelvinToMired(5000))
}

func TestMiredToKelvin(t *testing.T) {
	assert := assert.New(t)
	assert.Equal(2203, MiredToKelvin(500))
	assert.Equal(2203, MiredToKelvin(454))
	assert.Equal(2899, MiredToKelvin(345))
	assert.Equal(4000, MiredToKelvin(250))
	assert.Equal(4000, MiredToKelvin(220))
}

func TestMsToDuration(t *testing.T) {
	assert.Equal(t, 100000, MsToDuration(1000))
}

var hexRGBToColorXYDimTable = []struct {
	rgb  string
	x, y int
	dim  uint8
}{
	{"ff0000", 44506, 21022, 31},
	{"00ff00", 11299, 48941, 66},
	{"0000ff", 8880, 2613, 4},
	{"ffffff", 20943, 21992, 100},
}

func TestHexRGBToColorXYDim(t *testing.T) {
	assert := assert.New(t)
	var x, y int
	var dim uint8
	var err error

	for _, row := range hexRGBToColorXYDimTable {
		x, y, dim, err = HexRGBToColorXYDim(row.rgb)
		assert.Equal(row.x, x)
		assert.Equal(row.y, y)
		assert.Equal(row.dim, dim)
		assert.NoError(err)
	}
}

func TestKelvinToRGB(t *testing.T) {
	assert := assert.New(t)
	// pure white
	r, g, b := KelvinToRGB(6600)
	assert.Equal(1., r)
	assert.Equal(1., g)
	assert.Equal(1., b)
	// blueish
	r, g, b = KelvinToRGB(7000)
	assert.InDelta(.95, r, .01)
	assert.InDelta(.95, g, .01)
	assert.InDelta(1., b, .01)
	// redish
	r, g, b = KelvinToRGB(2000)
	assert.InDelta(1., r, .01)
	assert.InDelta(.53, g, .01)
	assert.InDelta(.05, b, .01)
}
