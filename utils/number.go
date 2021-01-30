package utils

import humanize "github.com/dustin/go-humanize"

const NUMBER_FORMAT_IDR = "#.###,##"

const NUMBER_FORMAT_THOUSAND_DOT = "#.###,##"

const NUMBER_FORMAT_THOUSAND_COMMA = "#.###,##"

type (
	number struct {
	}
)

func NewNumber() *number {
	return &number{}
}

func (n *number) FormatInteger(format string, number int) string {
	return humanize.FormatInteger(format, number)
}

func (n *number) FormatInteger16(format string, number int16) string {
	return humanize.FormatInteger(format, int(number))
}

func (n *number) FormatInteger32(format string, number int32) string {
	return humanize.FormatInteger(format, int(number))
}

func (n *number) FormatInteger64(format string, number int64) string {
	return humanize.FormatInteger(format, int(number))
}

func (n *number) FormatFloat64(format string, number float64) string {
	return humanize.FormatFloat(format, number)
}

func (n *number) FormatFloat32(format string, number float32) string {
	return humanize.FormatFloat(format, float64(number))
}
