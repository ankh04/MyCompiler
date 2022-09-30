package object

import "fmt"

// 值类型对象
// 包含三种数据类型：空值 布尔值 整数

type ObjectType string

const (
	INTEGER_OBJ = "INTEGER"
	BOOLEAN_OBJ = "BOOLEAN"
	NULL_OBJ    = "NULL"
)

type Object interface {
	Type() ObjectType
	Inspect() string
}

// region Integer

type Integer struct {
	Value int64
}

func (i *Integer) Type() ObjectType { return INTEGER_OBJ }

func (i *Integer) Inspect() string { return fmt.Sprintf("%d", i.Value) }

// endregion

// region Boolean

type Boolean struct {
	Value bool
}

func (b *Boolean) Type() ObjectType { return BOOLEAN_OBJ }

func (b *Boolean) Inspect() string { return fmt.Sprintf("%t", b.Value) }

// endregion

// region Null

type Null struct{}

func (n Null) Type() ObjectType { return NULL_OBJ }

func (n Null) Inspect() string { return "null" }

// endregion
