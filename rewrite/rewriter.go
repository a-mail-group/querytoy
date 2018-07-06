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


package rewriter

import "github.com/xwb1989/sqlparser"

func appendTableX(f *[]sqlparser.TableExprs, te sqlparser.TableExprs) {
	switch vte := te.(type) {
	case *sqlparser.JoinTableExpr:
		vte.LeftExpr
	}
}
func processSelect(s *sqlparser.Select) {
	{
		nfrom := make(sqlparser.TableExprs,0,len(s.From))
		for _,te := range s.From {
			switch vte := te.(type) {
			case *sqlparser.JoinTableExpr:
				vte.
			}
		}
	}
}

func ProcessStatement(s sqlparser.Statement) {
	switch v := s.(type) {
	case *sqlparser.Select:
		
	}
}

