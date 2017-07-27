package sql

import (
	tb "github.com/Cepave/open-falcon-backend/common/textbuilder"
)

var QQ = map[string]tb.Transformer{
	"x''": tb.BuildSurrounding(t.S("x'"), t.S("'")),
	"``":  tb.BuildSameSurrounding(t.S("`")),
}

var C = map[string]tb.TextGetter{
	"?":     t.S("?"),
	"where": t.S(" WHERE "),
	"and":   t.S(" AND "),
	"or":    t.S(" OR "),
}

var Post = &post{
	And: tb.BuildJoin(C["and"]),
	Or:  tb.BuildJoin(C["or"]),
}

type post struct {
	And tb.Distiller
	Or  tb.Distiller
}
