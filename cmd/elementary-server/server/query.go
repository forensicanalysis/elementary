package server

import (
	"fmt"

	"github.com/xwb1989/sqlparser"
)

func expandQuery(sql string) (string, error) {
	stmt, err := sqlparser.Parse(sql)
	if err != nil {
		return "", fmt.Errorf("%w ('%s')", err, sql)
	}

	tb := sqlparser.NewTrackedBuffer(func(buf *sqlparser.TrackedBuffer, node sqlparser.SQLNode) {
		if colIdent, ok := node.(sqlparser.ColIdent); ok && colIdent.String() != "json" {
			node = jsonExtractExpr(colIdent)
		}
		node.Format(buf)
	})
	stmt.Format(tb)
	return tb.String(), nil
}

func jsonExtractExpr(colIdent sqlparser.ColIdent) *sqlparser.FuncExpr {
	return &sqlparser.FuncExpr{
		Name: sqlparser.NewColIdent("json_extract"),
		Exprs: sqlparser.SelectExprs{
			&sqlparser.AliasedExpr{
				Expr: &sqlparser.ColName{Name: sqlparser.NewColIdent("json")},
			},
			&sqlparser.AliasedExpr{
				Expr: sqlparser.NewStrVal([]byte("$." + colIdent.String())),
			},
		},
	}
}
