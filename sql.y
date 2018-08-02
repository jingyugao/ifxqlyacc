%{
package ifxqlyacc

import (
    "time"
)

func setParseTree(yylex interface{},stmt Statement){
    yylex.(*Tokenizer).query.Statements = append(yylex.(*Tokenizer).query.Statements,stmt)
}

%}

%union{
    stmt                Statement
    stmts               Statements
    selStmt             *SelectStatement
    sdbStmt             *ShowDatabasesStatement
    cdbStmt             *CreateDatabaseStatement
    smmStmt             *ShowMeasurementsStatement
    str                 string
    query               Query
    field               *Field
    fields              Fields
    sources             Sources
    sortfs              SortFields
    sortf               *SortField
    ment                *Measurement
    dimens              Dimensions
    dimen               *Dimension
    int                 int
    int64               int64
    float64             float64
    expr                Expr
    tdur                time.Duration
    bool                bool
}

%token <str>    SELECT FROM WHERE AS GROUP BY ORDER LIMIT SHOW CREATE
%token <str>    DATABASES DATABASE MEASUREMENTS
%token <str>    COMMA SEMICOLON
%token <int>    MUL
%token <int>    EQ NEQ LT LTE GT GTE
%token <str>    IDENT
%token <int64>  INTEGER
%token <tdur>   DURATIONVAL
%token <str>    STRING
%token <bool>   DESC ASC
%token <float64> NUMBER
%left <int> AND OR

%type <stmt>                        STATEMENT
%type <sdbStmt>                     SHOW_DATABASES_STATEMENT
%type <cdbStmt>                     CREATE_DATABASE_STATEMENT
%type <selStmt>                     SELECT_STATEMENT
%type <smmStmt>                     SHOW_MEASUREMENTS_STATEMENT
%type <fields>                      COLUMN_NAMES
%type <field>                       COLUMN_NAME
%type <stmts>                       ALL_QUERIES
%type <sources>                     FROM_CLAUSE TABLE_NAMES
%type <ment>                        TABLE_NAME
%type <dimens>                      DIMENSION_NAMES GROUP_BY_CLAUSE
%type <dimen>                       DIMENSION_NAME
%type <expr>                        WHERE_CLAUSE CONDITION CONDITION_VAR OPERATION_EQUAL
%type <int>                         OPER LIMIT_INT
%type <sortfs>                      SORTFIELDS ORDER_CLAUSES
%type <sortf>                       SORTFIELD
%%
ALL_QUERIES:
        STATEMENT
        {
            setParseTree(yylex, $1)
        }
        | STATEMENT SEMICOLON
        {
            setParseTree(yylex, $1)
        }
        | STATEMENT SEMICOLON ALL_QUERIES
        {
            setParseTree(yylex, $1)
        }
STATEMENT:
    SELECT_STATEMENT
    {
        $$ = $1
    }
    |SHOW_DATABASES_STATEMENT
    {
        $$ = $1
    }
    |CREATE_DATABASE_STATEMENT
    {
        $$ = $1
    }
    |SHOW_MEASUREMENTS_STATEMENT
    {
        $$ = $1
    }
SELECT_STATEMENT:
    //SELECT COLUMN_NAMES
    //SELECT COLUMN_NAMES FROM_CLAUSE GROUP_BY_CLAUSE WHERE_CLAUSE ORDER_CLAUSES INTO_CLAUSE
    SELECT COLUMN_NAMES FROM_CLAUSE GROUP_BY_CLAUSE WHERE_CLAUSE ORDER_CLAUSES LIMIT_INT
    {
        sel := &SelectStatement{}
        sel.Fields = $2
        //sel.Target = $7
        sel.Sources = $3
        sel.Dimensions = $4
        sel.Condition = $5
        sel.SortFields = $6
        sel.Limit = $7
        $$ = sel
    }
COLUMN_NAMES:
    COLUMN_NAME
    {
        $$ = []*Field{$1}
    }
    |COLUMN_NAME COMMA COLUMN_NAMES
    {
        $$ = append($3,$1)
    }
COLUMN_NAME:
    MUL
    {
        $$ = &Field{Expr:&Wildcard{Type:$1}}
    }
    |IDENT
    {
        $$ = &Field{Expr:&VarRef{Val:$1}}
    }
    |IDENT AS IDENT
    {
        $$ = &Field{Expr:&VarRef{Val:$1},Alias:$3}
    }
FROM_CLAUSE:
    FROM TABLE_NAMES
    {
        $$ = $2
    }
    |
    {
        $$ = nil
    }
TABLE_NAMES:
    TABLE_NAME
    {
        $$ = []Source{$1}
    }
    |TABLE_NAME COMMA TABLE_NAMES
    {
        $$ = append($3,$1)
    }
TABLE_NAME:
    IDENT
    {
        $$ = &Measurement{Name:$1}

    }
GROUP_BY_CLAUSE:
    GROUP BY DIMENSION_NAMES
    {
        $$ = $3
    }
    |
    {
        $$ = nil
    }
DIMENSION_NAMES:
    DIMENSION_NAME
    {
        $$ = []*Dimension{$1}
    }
    |DIMENSION_NAME COMMA DIMENSION_NAMES
    {
        $$ = append($3,$1)
    }
DIMENSION_NAME:
    IDENT
    {
        $$ = &Dimension{Expr:&VarRef{Val:$1}}

    }
WHERE_CLAUSE:
    WHERE CONDITION
    {
        $$ = $2
    }
    |
    {
        $$ = nil
    }
CONDITION:
    OPERATION_EQUAL
    {
        $$ = $1
    }
    |CONDITION AND CONDITION
    {
        $$ = &BinaryExpr{Op:$2,LHS:$1,RHS:$3}
    }
    |CONDITION OR CONDITION
    {
        $$ = &BinaryExpr{Op:$2,LHS:$1,RHS:$3}

    }
OPERATION_EQUAL:
    CONDITION_VAR OPER CONDITION_VAR
    {
        $$ = &BinaryExpr{Op:$2,LHS:$1,RHS:$3}
    }
OPER:
    EQ
    {
        $$ = $1
    }
    |NEQ
    {
        $$ = $1
    }
    |LT
    {
        $$ =$1
    }
    |LTE
    {
        $$ = $1
    }
    |GT
    {
        $$ = $1
    }
    |GTE
    {
        $$ = $1
    }
CONDITION_VAR:
    IDENT
    {
        $$ = &VarRef{Val:$1}
    }
    |NUMBER
    {
        $$ = &NumberLiteral{Val:$1}
    }
    |INTEGER
    {
        $$ = &IntegerLiteral{Val:$1}
    }
    |DURATIONVAL
    {
        $$ = &DurationLiteral{Val:$1}
    }
    |STRING
    {
        $$ = &StringLiteral{Val:$1}
    }
ORDER_CLAUSES:
    ORDER BY SORTFIELDS
    {
        $$ = $3
    }
    |
    {
        $$ = nil
    }
SORTFIELDS:
    SORTFIELD
    {
        $$ = []*SortField{$1}
    }
    |SORTFIELD COMMA SORTFIELDS
    {
        $$ = append($3,$1)
    }
SORTFIELD:
    IDENT
    {
        $$ = &SortField{Name:$1}
    }
    |IDENT DESC
    {
        $$ = &SortField{Name:$1,Ascending:$2}
    }
    |IDENT ASC
    {
        $$ = &SortField{Name:$1,Ascending:$2}
    }
LIMIT_INT:
    LIMIT INTEGER
    {
        $$ = int($2)
    }
    |
    {
        $$ = 0
    }
SHOW_DATABASES_STATEMENT:
    SHOW DATABASES
    {
        $$ = &ShowDatabasesStatement{}
    }
CREATE_DATABASE_STATEMENT:
    CREATE DATABASE IDENT
    {
        $$ = &CreateDatabaseStatement{Name:$3}
    }
SHOW_MEASUREMENTS_STATEMENT:
    SHOW MEASUREMENTS WHERE_CLAUSE ORDER_CLAUSES LIMIT_INT
    {
        sms := &ShowMeasurementsStatement{}
        sms.Condition = $3
        sms.SortFields = $4
        sms.Limit = $5
        $$ = sms
    }
%%
