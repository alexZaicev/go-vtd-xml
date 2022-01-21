package reader

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	xmlDoc = `
<?xml version="1.0" encoding="ASCII" standalone="no" ?>
<Clocks timezone="GMT">
	<timehour>11</timehour>
	<timeminute>50</timeminute>
	<timesecond>40</timesecond>
	<timemeridian>p.m.</timemeridian>
</Clocks>
`
)

func Test_AsciiReader_GetChar_Success(t *testing.T) {
	reader, err := NewAsciiReader([]byte(xmlDoc), 0, len(xmlDoc))
	assert.Nil(t, err)
	assert.NotNil(t, reader)

	for {
		ch, err := reader.GetChar()
		if err != nil {
			break
		}
		fmt.Printf("%c", ch)
	}

}
