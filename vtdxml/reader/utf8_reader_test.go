package reader

import (
	"github.com/stretchr/testify/assert"
	"io"
	"testing"
	"unicode/utf8"
)

const (
	utf8XmlDoc = `
<?xml version="1.0" encoding="ASCII" standalone="no" ?>
<Clocks timezone="GMT">
	<timehour>11</timehour>
	<timeminute>50</timeminute>
	<timesecond>40</timesecond>
	<timemeridian>p.m.</timemeridian>
	<name>µ$&£</name>
</Clocks>
`
	utf8InvalidXmlDoc = `
<?xml version="1.0" encoding="UTF-8" standalone="no" ?>
<Clocks timezone="GMT">
	<timehour>\xfc\xa1\xa1\xa1\xa1\xa1</timehour>
	<timeminute>50</timeminute>
	<timesecond>40</timesecond>
	<timemeridian>p.m.</timemeridian>
	<name>\xe2\x28\xa1</name>
</Clocks>
`
)

func Test_Utf8Reader_GetChar_Success(t *testing.T) {
	reader, err := NewUtf8Reader([]byte(utf8XmlDoc), 0, len(utf8XmlDoc)-1)
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

func Test_Utf8Reader_GetChar_ParseException(t *testing.T) {
	assert.False(t, utf8.Valid([]byte(utf8InvalidXmlDoc)))

	//reader, err := NewUtf8Reader([]byte(v), 0, len(v)-1)
	//assert.Nil(t, err)
	//assert.NotNil(t, reader)
	//
	//for {
	//	ch, err := reader.GetChar()
	//	if err == io.EOF {
	//		break
	//	} else if err != nil {
	//		assert.EqualError(t, err, "a parse error occurred: invalid ASCII character")
	//	}
	//	fmt.Printf("%c", ch)
	//}
}
