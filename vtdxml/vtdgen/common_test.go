package vtdgen

import (
	"fmt"
	"testing"
)

var xmlDoc = `
<?xml version="1.0" encoding="ASCII" standalone="no" ?>
<Clocks timezone="GMT">
	<timehour>11</timehour>
	<timeminute>50</timeminute>
	<timesecond>40</timesecond>
	<timemeridian>p.m.</timemeridian>
</Clock
`

func Test(t *testing.T) {
	v := VtdGen{
		encoding: FormatUtf16LE,
		offset:   15,
		xmlDoc:   []byte(xmlDoc),
	}

	fmt.Println(v.getPrevOffset())
}
