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
import "strings"
import uuid "github.com/nu7hatch/gouuid"
import "fmt"

func MakeColName(schema,name,col string) *sqlparser.ColName {
	r := new(sqlparser.ColName)
	if schema!="" { r.Qualifier.Qualifier = sqlparser.NewTableIdent(schema) }
	if name  !="" { r.Qualifier.Name      = sqlparser.NewTableIdent(name  ) }
	if col   !="" { r.Name                = sqlparser.NewColIdent  (col   ) }
	return r
}
func MakeTableName(schema,name string) *sqlparser.TableName {
	r := new(sqlparser.TableName)
	if schema!="" { r.Qualifier = sqlparser.NewTableIdent(schema) }
	if name  !="" { r.Name      = sqlparser.NewTableIdent(schema) }
	return r
}

type aliasedSrc struct{
	schema,name string
	colOrder []string
	colMap map[string]*sqlparser.ColName
}
func aliasedSrc_mk(schema,name string, altt *sqlparser.TableName, cols []string) (src aliasedSrc) {
	src.schema = schema
	src.name = name
	src.colOrder = cols
	src.colMap = make(map[string]*sqlparser.ColName)
	for _,col := range cols {
		cn := MakeColName(schema,name,col)
		if altt!=nil { cn.Qualifier = *altt }
		src.colMap[col] = cn
	}
	return
}
func (a *aliasedSrc) match(name sqlparser.TableName) bool {
	q := strings.ToLower(name.Qualifier.String())
	n := strings.ToLower(name.Name.String())
	return (q=="" || q==a.schema) && n==a.name
}

type queryContext struct{
	src []aliasedSrc
	lim int
}
func (q *queryContext) replace(node sqlparser.SQLNode) (kontinue bool, err error) {
	// Exclude unwanted recursion
	switch node.(type) {
	case *sqlparser.Subquery: return
	}
	
	switch v := node.(type) {
	case *sqlparser.ColName:
		if v.Qualifier.Name.String()=="" {
			for _,a := range q.src {
				coln := a.colMap[strings.ToLower(v.Name.String())]
				if coln!=nil { *v = *coln; return }
			}
		} else {
			for _,a := range q.src {
				if !a.match(v.Qualifier) { continue }
				coln := a.colMap[strings.ToLower(v.Name.String())]
				if coln!=nil { *v = *coln; return }
			}
		}
		if len(q.src)>1 {
			err = fmt.Errorf("unkown column: %s",sqlparser.String(v))
		}
		return
	}
	return true,nil
}

func (q *queryContext) collect(c *Conn,expr sqlparser.TableExpr) error {
	var xerr error
	ste := func(name string, sexpr sqlparser.SimpleTableExpr) (altt *sqlparser.TableName) {
		switch v := sexpr.(type) {
		case sqlparser.TableName:
			Q := strings.ToLower(v.Qualifier.String())
			N := strings.ToLower(v.Name.String())
			if Q=="" { Q = c.CurrentSchema }
			sch := c.Schemas[Q]
			if sch==nil { return }
			tab := sch.Tables[N]
			if tab==nil { return }
			if name!="" { Q,N = "",name } else if v.Qualifier.String()=="" { Q = "" }
			q.src = append(q.src,aliasedSrc_mk(Q,N,altt,tab.Columns))
		case *sqlparser.Subquery:
			sel := findDefinition(v.Select)
			if sel==nil { return }
			ncol := make([]string,0,len(sel.SelectExprs))
			for _,sx := range sel.SelectExprs {
				switch v2 := sx.(type) {
				case *sqlparser.AliasedExpr:
					if as := v2.As.String(); as!="" {
						ncol = append(ncol,as)
					} else {
						ncol = append(ncol,getName(v2.Expr))
					}
				case sqlparser.Nextval:
					ncol = append(ncol,getName(v2.Expr))
				}
			}
			if name=="" {
				u,e := uuid.NewV4()
				if e!=nil { xerr = e; return }
				altt = MakeTableName("",u.String())
			}
			q.src = append(q.src,aliasedSrc_mk("",name,altt,ncol))
		}
		return
	}
	return sqlparser.Walk(func(node sqlparser.SQLNode) (kontinue bool, err error) {
		switch node.(type){
		case sqlparser.TableExpr,sqlparser.TableExprs:
		default: return
		}
		
		switch v := node.(type) {
		case *sqlparser.AliasedTableExpr:
			altt := ste(v.As.String(),v.Expr)
			if altt!=nil {
				v.As = altt.Name
			}
			err = xerr
		default:
			kontinue = true
		}
		return
	},expr)
}

func (c *Conn) star(ctx *queryContext) sqlparser.Visit {
	perform := func(sxp *sqlparser.SelectExprs) {
		var nsx sqlparser.SelectExprs
		for i,e := range (*sxp) {
			if se,ok := e.(*sqlparser.StarExpr); ok {
				if nsx==nil {
					nsx = make(sqlparser.SelectExprs,i)
					copy(nsx,*sxp)
				}
				if se.TableName.Name.String()=="" {
					for _,a := range ctx.src[:ctx.lim] {
						for _,col := range a.colOrder {
							nsx = append(nsx,&sqlparser.AliasedExpr{Expr:a.colMap[col]})
						}
					}
				} else {
					for _,a := range ctx.src[:ctx.lim] {
						if !a.match(se.TableName) { continue }
						for _,col := range a.colOrder {
							nsx = append(nsx,&sqlparser.AliasedExpr{Expr:a.colMap[col]})
						}
						break
					}
				}
			} else if nsx!=nil { nsx = append(nsx,e) }
			if nsx!=nil { *sxp = nsx }
		}
	}
	return func(node sqlparser.SQLNode) (kontinue bool, err error) {
		switch v := node.(type) {
		case *sqlparser.Subquery: return
		case *sqlparser.FuncExpr: perform(&v.Exprs)
		case *sqlparser.GroupConcatExpr: perform(&v.Exprs)
		case *sqlparser.MatchExpr: perform(&v.Columns)
		case *sqlparser.Select:
			perform(&v.SelectExprs)
			kontinue = true
		default:
			kontinue = true
		}
		return
	}
}
func (c *Conn) visitor(par *queryContext) sqlparser.Visit {
	var mv,recurse sqlparser.Visit
	
	mv = func(node sqlparser.SQLNode) (kontinue bool, err error) {
		// Exclude...
		switch node.(type){
		// (SELECT ... UNION SELECT ...)
		case *sqlparser.ParenSelect,
		// SELECT ... UNION SELECT ...
			*sqlparser.Union:
			return true,nil
		}
		
		sel := getSelect(node)
		sqlparser.Walk(recurse,sel.from)
		ctx := new(queryContext)
		for _,f := range sel.from { ctx.collect(c,f) }
		ctx.lim = len(ctx.src)
		if par!=nil { ctx.src = append(ctx.src,par.src...) }
		sqlparser.Walk(ctx.replace,sel.exprs)
		for _,sq := range sel.xsq { sqlparser.Walk(c.visitor(ctx),sq.Select) }
		sqlparser.Walk(c.star(ctx),node)
		return
	}
	recurse = subqueries(mv)
	return mv
}
func (c *Conn) Rewrite(s sqlparser.Statement) {
	sqlparser.Walk(c.visitor(nil),s)
}
func ProcessStatement(s sqlparser.Statement) {
	switch s.(type) {
	case *sqlparser.Select:
	}
}

