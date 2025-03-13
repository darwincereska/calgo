package calculator

type Value struct {
	Number float64
	Unit string
	Type ValueType
	IsError bool
	ErrorMsg string
	Label string // For labeled expressions
}

type ValueType int

const (
	NumberType ValueType = iota
	CurrencyType
	UnitType
	TimeType
	PercentageType
	DateType
	TemperatureType
	CommentType
	LabelType
)

type Context struct {
	Variables map[string]Value
	TimeZone string
	PPI float64
	EmSize float64
}

type Expression interface {
	Evaluate(ctx *Context) Value
}