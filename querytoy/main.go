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


package main

import "github.com/xwb1989/sqlparser"
import "fmt"
import "reflect"
import _ "github.com/a-mail-group/querytoy/catalog"
//import "github.com/a-mail-group/querytoy/rewrite"
import "encoding/json"

func Type(i interface{}) interface{} {
	if i==nil { return nil }
	return reflect.ValueOf(i).Type()
}

func main(){
	queries := []string{
		"SELECT * FROM MyUsers JOIN OtherUser WHERE a = 'abc'",
		"SELECT * FROM MyUsers JOIN OtherUser ON MyUsers.x = OtherUser.y WHERE a = 'abc'",
		"SELECT * FROM MyUsers, OtherUser WHERE a = 'abc'",
		"SELECT * WHERE a = 'abc'",
		"SELECT 1",
		"DELETE FROM MyUsers WHERE a = 'abc'",
	}
	for _,query := range queries {
		ntb := sqlparser.NewTrackedBuffer(nil)
		stm,err := sqlparser.Parse(query)
		if err!=nil { fmt.Println(err); continue }
		ntb.WriteNode(stm)
		fmt.Println(ntb)
	}
	x,_ := sqlparser.Parse("SELECT test,wax,1,'1'")
	b,_ := json.MarshalIndent(x,""," ")
	fmt.Println(string(b))
}

