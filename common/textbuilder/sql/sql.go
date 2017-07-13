package sql

import (
	tb "github.com/Cepave/open-falcon-backend/common/textbuilder"
)

// Builds prefixing text getter if the where condition is viable
func Where(condition tb.TextGetter) tb.TextGetter {
	return tb.Prefix(C["where"], condition)
}

func And(getters ...tb.TextGetter) tb.TextGetter {
	return tb.Join(C["and"], getters...)
}

func Or(getters ...tb.TextGetter) tb.TextGetter {
	return tb.Join(C["or"], getters...)
}

func In(column tb.TextGetter, number int) tb.TextGetter {
	return sqlInImpl(column, tb.Repeat(C["?"], number))
}

func InByLen(column tb.TextGetter, arrayObject interface{}) tb.TextGetter {
	return sqlInImpl(column, tb.RepeatByLen(C["?"], arrayObject))
}

func sqlInImpl(column tb.TextGetter, list tb.TextList) tb.TextGetter {
	return tb.Surrounding(
		t.S(column.String()+" IN ( "),
		list.Post().Distill(tb.J[", "]),
		t.S(")"),
	)
}
