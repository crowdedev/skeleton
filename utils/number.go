package utils

import humanize "github.com/dustin/go-humanize"

const NUMBER_FORMAT_IDR = "#.###,##"

const NUMBER_FORMAT_THOUSAND_DOT = "#.###,##"

const NUMBER_FORMAT_THOUSAND_COMMA = "#.###,##"

type (
	Number struct {
	}
)

func (n *Number) FormatInteger(format string, number int) string {
	return humanize.FormatInteger(format, number)
}

func (n *Number) FormatInteger16(format string, number int16) string {
	return humanize.FormatInteger(format, int(number))
}

func (n *Number) FormatInteger32(format string, number int32) string {
	return humanize.FormatInteger(format, int(number))
}

func (n *Number) FormatInteger64(format string, number int64) string {
	return humanize.FormatInteger(format, int(number))
}

func (n *Number) FormatFloat64(format string, number float64) string {
	return humanize.FormatFloat(format, number)
}

func (n *Number) FormatFloat32(format string, number float32) string {
	return humanize.FormatFloat(format, float64(number))
}
