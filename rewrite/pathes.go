/*
Copyright (c) 2018 Simon Schmidt

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/


package rewrite

import "gopkg.in/src-d/go-vitess.v1/vt/sqlparser"

// Select, Update, Delete
type sqlSelect struct{
	from  sqlparser.TableExprs
	exprs sqlparser.Exprs
	xsq   []*sqlparser.Subquery
}

func getSelect(node sqlparser.SQLNode) (sel sqlSelect) {
	var lvjoin sqlparser.Visit
	lvjoin = func(node sqlparser.SQLNode) (kontinue bool, err error) {
		// Fetch the interesting stuff.
		switch v := node.(type) {
		case *sqlparser.JoinTableExpr:
			e := v.Condition.On
			sel.exprs = append(sel.exprs,e)
		}
		// do processing on TableExprs and TableExpr only.
		switch node.(type) {
		case sqlparser.TableExpr,sqlparser.TableExprs:
			kontinue = true
		}
		return
	}
	sqlparser.Walk(func(node sqlparser.SQLNode) (kontinue bool, err error) {
		switch v := node.(type) {
		case sqlparser.TableExprs:
			sel.from = v
			sqlparser.Walk(lvjoin,v)
		case sqlparser.Expr:
			sel.exprs = append(sel.exprs,v)
		default:
			kontinue = true
		}
		return
	},node)
	sqlparser.Walk(func(node sqlparser.SQLNode) (kontinue bool, err error) {
		switch v := node.(type) {
		case *sqlparser.Subquery:
			sel.xsq = append(sel.xsq,v)
		default: kontinue = true
		}
		return
	},sel.exprs)
	
	return
}

//  TableExprs-> ({ParenTableExpr,JoinTableExpr}->*) AliasedTableExpr-> Subquery
func subqueries(vis sqlparser.Visit) sqlparser.Visit {
	var step sqlparser.Visit
	step = func(node sqlparser.SQLNode) (kontinue bool, err error) {
		switch v := node.(type) {
		case *sqlparser.JoinTableExpr:
			return false,sqlparser.Walk(step,v.LeftExpr,v.RightExpr)
		case *sqlparser.Subquery:
			return vis(v.Select)
		}
		switch node.(type) {
		case sqlparser.TableExprs,sqlparser.TableExpr: return true,nil
		}
		return
	}
	return step
}
func findDefinition(i sqlparser.SelectStatement) *sqlparser.Select {
restart:
	switch v := i.(type) {
	case *sqlparser.Select: return v
	case *sqlparser.Union:
		i = v.Left
		goto restart
	case *sqlparser.ParenSelect:
		i = v.Select
		goto restart
	}
	// Unlikely!
	return nil
}
func getName(i sqlparser.Expr) string {
restart:
	switch v := i.(type) {
	case *sqlparser.Default: return v.ColName
	case *sqlparser.ColName: return v.Name.String()
	case *sqlparser.ParenExpr: i = v.Expr; goto restart
	}
	return "?column?"
}

