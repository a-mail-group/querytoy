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


package catalog

type ObjectImpl struct{
	FName,FSchema,FInternalID string
}
func (o *ObjectImpl) Name() string { return o.FName }
func (o *ObjectImpl) Schema() string { return o.FSchema }
func (o *ObjectImpl) InternalID() string { return o.FInternalID }

var _ Object = (*ObjectImpl)(nil)

type ColumnImpl struct{
	FName,FType string
	FIndex interface{}
}
func (c *ColumnImpl) Name() string       { return c.FName  }
func (c *ColumnImpl) Type() string       { return c.FType  }
func (c *ColumnImpl) Index() interface{} { return c.FIndex }

var _ Column = (*ColumnImpl)(nil)

type RelationImpl struct{
	ObjectImpl
	FColumns []Column
}
func (c *RelationImpl) Columns() []Column { return c.FColumns }
func (c *RelationImpl) Column(name string) Column {
	for _,c := range c.FColumns {
		if c.Name()==name { return c }
	}
	return nil
}
var _ Relation = (*RelationImpl)(nil)

