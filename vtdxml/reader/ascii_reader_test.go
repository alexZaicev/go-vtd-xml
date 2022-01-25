package reader

import (
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	asciiXmlDoc = `
<?xml version="1.0" encoding="ASCII" standalone="no" ?>
<Clocks timezone="GMT">
	<timehour>11</timehour>
	<timeminute>50</timeminute>
	<timesecond>40</timesecond>
	<timemeridian>p.m.</timemeridian>
</Clocks>
`
	asciiInvalidXmlDoc = `
<?xml version="1.0" encoding="ASCII" standalone="no" ?>
<Clocks timezone="GMT">
	<timehour>11</timehour>
	<timeminute>50</timeminute>
	<timesecond>40</timesecond>
	<timemeridian>p.m.</timemeridian>
	<name>µ$&£</name>
</Clocks>
`
)

func Test_AsciiReader_GetChar_Success(t *testing.T) {
	reader, err := NewAsciiReader([]byte(asciiXmlDoc), 0, len(asciiXmlDoc)-1)
	assert.Nil(t, err)
	assert.NotNil(t, reader)

	for {
		_, err := reader.GetChar()
		if err == io.EOF {
			break
		}
		assert.Nil(t, err)
	}
}

func Test_AsciiReader_GetChar_ParseException(t *testing.T) {
	reader, err := NewAsciiReader([]byte(asciiInvalidXmlDoc), 0, len(asciiInvalidXmlDoc)-1)
	assert.Nil(t, err)
	assert.NotNil(t, reader)

	for {
		_, err := reader.GetChar()
		if err == io.EOF {
			break
		} else if err != nil {
			assert.EqualError(t, err, "a parse error occurred: invalid ASCII character")
		}
	}
}
