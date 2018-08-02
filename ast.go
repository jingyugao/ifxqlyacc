package ifxqlyacc

import (
	"regexp"
	"time"
)
// DataType represents the primitive data types available in InfluxQL.
type DataType int

const (
	// Unknown primitive data type.
	Unknown DataType = 0
	// Float means the data type is a float.
	Float DataType = 1
	// Integer means the data type is an integer.
	Integer DataType = 2
	// String means the data type is a string of text.
	String DataType = 3
	// Boolean means the data type is a boolean.
	Boolean DataType = 4
	// Time means the data type is a time.
	Time DataType = 5
	// Duration means the data type is a duration of time.
	Duration DataType = 6
	// Tag means the data type is a tag.
	Tag DataType = 7
	// AnyField means the data type is any field.
	AnyField DataType = 8
	// Unsigned means the data type is an unsigned integer.
	Unsigned DataType = 9
)
type Sources []Source
type Source interface {
	Node
	source()
}

//source的实现,有可能是表，有可能是子查询
func (*Measurement)	source(){}
func (*SubQuery)	source(){}

type Statement interface{
	Node
	stmt()
}
func (*SelectStatement) 			stmt() 		{}
func (*ShowDatabasesStatement)		stmt()		{}
func (*CreateDatabaseStatement)		stmt()		{}
func (*ShowMeasurementsStatement) 	stmt()		{}

type Node interface {
	node()
	String() string
}
// 实现Node
func (*SelectStatement) 			node()		{}
func (*ShowDatabasesStatement)		node()		{}
func (*CreateDatabaseStatement)		node()		{}
func (*ShowMeasurementsStatement) 	node()	{}
func (Fields)						node()		{}
func (*Target)						node()		{}
func (Sources)						node()		{}
func (*SortField)					node()      {}
func (SortFields)					node() 		{}
func (*Field)						node() 		{}
func (Dimensions)					node()      {}
func (*Measurement)					node() 		{}
func (*Dimension)					node() 		{}
func (*SubQuery)					node() 		{}

func (*StringLiteral)				node() 		{}
func (*VarRef)						node()   	{}
func (*Wildcard)					node() 		{}
func (*NumberLiteral)				node()		{}
func (*IntegerLiteral)				node()		{}
func (*DurationLiteral)				node()		{}
func (*BinaryExpr)					node()		{}


type Expr interface{
	Node
	expr()
}
func (*Wildcard) 				expr() 		{}
func (*StringLiteral)			expr()		{}
func (*VarRef)					expr()		{}
func (*NumberLiteral)			expr()		{}
func (*IntegerLiteral)			expr()		{}
func (*DurationLiteral) 		expr()		{}
func (*BinaryExpr)				expr()		{}


type Query struct {
	Statements Statements
}

type Statements []Statement

// ShowDatabasesStatement represents a command for listing all databases in the cluster.
type ShowDatabasesStatement struct{}


// ShowMeasurementsStatement represents a command for listing measurements.
type ShowMeasurementsStatement struct {
	// Database to query. If blank, use the default database.
	Database string

	// Measurement name or regex.
	Source Source

	// An expression evaluated on data point.
	Condition Expr

	// Fields to sort results by
	SortFields SortFields

	// Maximum number of rows to be returned.
	// Unlimited if zero.
	Limit int

	// Returns rows starting at an offset from the first row.
	Offset int
}

// CreateDatabaseStatement represents a command for creating a new database
type CreateDatabaseStatement struct {
	// Name of the database to be created.
	Name string

	// RetentionPolicyCreate indicates whether the user explicitly wants to create a retention policy.
	RetentionPolicyCreate bool

	// RetentionPolicyDuration indicates retention duration for the new database.
	RetentionPolicyDuration *time.Duration

	// RetentionPolicyReplication indicates retention replication for the new database.
	RetentionPolicyReplication *int

	// RetentionPolicyName indicates retention name for the new database.
	RetentionPolicyName string

	// RetentionPolicyShardGroupDuration indicates shard group duration for the new database.
	RetentionPolicyShardGroupDuration time.Duration
}



type SelectStatement struct {
	//返回的列名，支持别名，*等
	Fields Fields

	//select into 使用的表
	Target *Target

	//group by 使用
	Dimensions Dimensions

	//要查询的表，一般由measurement实现，也可由子查询实现 select from select
	Sources Sources

	//where
	Condition Expr

	//sort
	SortFields SortFields

	//
	Limit int

	Offset int

	SLimit int

	SOffset int

	groupByInterval time.Duration

	IsRawQuery bool

	Fill FillOption

	FillValue interface{}

	Location *time.Location

	TimeAlias string

	OmitTime bool

	StripName bool

	EmitName string

	Dedupe bool

	JoinType int
}

type Measurements []*Measurement
type Measurement struct{
	Database string
	RetentionPolicy string
	Name string
	Regex *RegexLiteral

	SystemIterator string
}

//SubQuery is a source with a SelectStatement as the backing store.
type SubQuery struct {
	Statement *SelectStatement
}

type Fields []*Field
type Field struct{
	Expr Expr
	Alias string
}
type Target struct{
	Measurement *Measurement
}


type Dimensions []*Dimension
type Dimension struct {
	Expr Expr
}

type SortFields []*SortField
type SortField struct {
	Name string
	Ascending bool
}

type FillOption int

type RegexLiteral struct {
	Val *regexp.Regexp
}

// Wildcard represents a wild card expression.
type Wildcard struct {
	Type int
}

// NumberLiteral represents a numeric literal.
type NumberLiteral struct {
	Val float64
}

// BinaryExpr represents an operation between two expressions.
type BinaryExpr struct {
	Op  int
	LHS Expr
	RHS Expr
}

// StringLiteral represents a string literal.
type StringLiteral struct {
	Val string
}

// VarRef represents a reference to a variable.
type VarRef struct {
	Val  string
	Type DataType
}

// DurationLiteral represents a duration literal.
type DurationLiteral struct {
	Val time.Duration
}

// IntegerLiteral represents an integer literal.
type IntegerLiteral struct {
	Val int64
}
//深度优先节点遍历
func Walk(v Visitor,node Node){
	switch n:=node.(type) {
	case *SelectStatement:{
		Walk(v, n.Fields)
		Walk(v, n.Target)
		Walk(v, n.Dimensions)
		Walk(v, n.Sources)
		Walk(v, n.Condition)
		Walk(v, n.SortFields)
	}
	case Fields:{
		for _, c := range n{
			Walk(v,c)
		}
	}
	case *Target:{
		Walk(v,n.Measurement)
	}
	case Dimensions:{
		for _, c := range n{
			Walk(v,c)
		}
	}
	case Sources:{
		for _,s := range n{
			Walk(v,s)
		}
	}
	case *SubQuery:
		Walk(v,n.Statement)
	case *Dimension:{
		Walk(v,n.Expr)
	}
	case *Field:{
		Walk(v,n.Expr)
	}
	}
}

type Visitor interface{
	Visit(Node) Visitor
}