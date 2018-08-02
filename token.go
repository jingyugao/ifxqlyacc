package ifxqlyacc

import "strings"

// Token is a lexical token of the InfluxQL language.
//type Token int

// These are a comprehensive list of InfluxQL language tokens.
const (
//	// ILLEGAL Token, EOF, WS are Special InfluxQL tokens.
	ILLEGAL int = iota
	EOF
	WS
	COMMENT
//
	literalBeg
//	// IDENT and the following are InfluxQL literal tokens.
	//IDENT       // main
	BOUNDPARAM  // $param
	//NUMBER      // 12345.67
	//INTEGER     // 12345
	//DURATIONVAL // 13h
	//STRING      // "abc"
	BADSTRING   // "abc
	BADESCAPE   // \q
	TRUE        // true
	FALSE       // false
	REGEX       // Regular expressions
	BADREGEX    // `.*
	literalEnd
//
	operatorBeg
//	// ADD and the following are InfluxQL Operators
	ADD         // +
	SUB         // -
	//MUL         // *
	DIV         // /
	MOD         // %
	BITWISE_AND // &
	BITWISE_OR  // |
	BITWISE_XOR // ^

	//AND // AND
	//OR  // OR
	//
	//EQ       // =
	//NEQ      // !=
	EQREGEX  // =~
	NEQREGEX // !~
	//LT       // <
	//LTE      // <=
	//GT       // >
	//GTE      // >=
	operatorEnd

	LPAREN      // (
	RPAREN      // )
	//COMMA       // ,
	COLON       // :
	DOUBLECOLON // ::
	//SEMICOLON   // ;
	DOT         // .
//
//	keywordBeg
//	// ALL and the following are InfluxQL Keywords
//	ALL
//	ALTER
//	ANALYZE
//	ANY
//	//AS
//	ASC
//	BEGIN
//	BY
//	CARDINALITY
//	CREATE
//	CONTINUOUS
//	DATABASE
//	DATABASES
//	DEFAULT
//	DELETE
//	DESC
//	DESTINATIONS
//	DIAGNOSTICS
//	DISTINCT
//	DROP
//	DURATION
//	END
//	EVERY
//	EXACT
//	EXPLAIN
//	FIELD
//	FOR
//	FROM
//	GRANT
//	GRANTS
//	GROUP
//	GROUPS
//	IN
//	INF
//	INSERT
//	INTO
//	INNER_JOIN
//	LEFT_JOIN
//	RIGHT_JOIN
//	JOIN
//	KEY
//	KEYS
//	KILL
//	LIMIT
//	MEASUREMENT
//	MEASUREMENTS
//	NAME
//	OFFSET
//	ON
//	ORDER
//	PASSWORD
//	POLICY
//	POLICIES
//	PRIVILEGES
//	QUERIES
//	QUERY
//	READ
//	REPLICATION
//	RESAMPLE
//	RETENTION
//	REVOKE
//	//SELECT
//	SERIES
//	SET
//	SHOW
//	SHARD
//	SHARDS
//	SLIMIT
//	SOFFSET
//	STATS
//	SUBSCRIPTION
//	SUBSCRIPTIONS
//	TAG
//	TO
//	USER
//	USERS
//	VALUES
//	WHERE
//	WITH
//	WRITE
//	keywordEnd
)
//
//var tokens = [...]string{
//	ILLEGAL: "ILLEGAL",
//	EOF:     "EOF",
//	WS:      "WS",
//
//	IDENT:       "IDENT",
//	NUMBER:      "NUMBER",
//	DURATIONVAL: "DURATIONVAL",
//	STRING:      "STRING",
//	BADSTRING:   "BADSTRING",
//	BADESCAPE:   "BADESCAPE",
//	TRUE:        "TRUE",
//	FALSE:       "FALSE",
//	REGEX:       "REGEX",
//
//	ADD:         "+",
//	SUB:         "-",
//	MUL:         "*",
//	DIV:         "/",
//	MOD:         "%",
//	BITWISE_AND: "&",
//	BITWISE_OR:  "|",
//	BITWISE_XOR: "^",
//
//	AND: "AND",
//	OR:  "OR",
//
//	EQ:       "=",
//	NEQ:      "!=",
//	EQREGEX:  "=~",
//	NEQREGEX: "!~",
//	LT:       "<",
//	LTE:      "<=",
//	GT:       ">",
//	GTE:      ">=",
//
//	LPAREN:      "(",
//	RPAREN:      ")",
//	COMMA:       ",",
//	COLON:       ":",
//	DOUBLECOLON: "::",
//	SEMICOLON:   ";",
//	DOT:         ".",
//
//	ALL:           "ALL",
//	ALTER:         "ALTER",
//	ANALYZE:       "ANALYZE",
//	ANY:           "ANY",
//	//AS:            "AS",
//	ASC:           "ASC",
//	BEGIN:         "BEGIN",
//	BY:            "BY",
//	CARDINALITY:   "CARDINALITY",
//	CREATE:        "CREATE",
//	CONTINUOUS:    "CONTINUOUS",
//	DATABASE:      "DATABASE",
//	DATABASES:     "DATABASES",
//	DEFAULT:       "DEFAULT",
//	DELETE:        "DELETE",
//	DESC:          "DESC",
//	DESTINATIONS:  "DESTINATIONS",
//	DIAGNOSTICS:   "DIAGNOSTICS",
//	DISTINCT:      "DISTINCT",
//	DROP:          "DROP",
//	DURATION:      "DURATION",
//	END:           "END",
//	EVERY:         "EVERY",
//	EXACT:         "EXACT",
//	EXPLAIN:       "EXPLAIN",
//	FIELD:         "FIELD",
//	FOR:           "FOR",
//	FROM:          "FROM",
//	GRANT:         "GRANT",
//	GRANTS:        "GRANTS",
//	GROUP:         "GROUP",
//	GROUPS:        "GROUPS",
//	IN:            "IN",
//	INF:           "INF",
//	INSERT:        "INSERT",
//	INTO:          "INTO",
//	INNER_JOIN:	"INNER JOIN",
//	LEFT_JOIN:		"LEFT JOIN",
//	RIGHT_JOIN:	"RIGHT JOIN",
//	JOIN:			"JOIN",
//	KEY:           "KEY",
//	KEYS:          "KEYS",
//	KILL:          "KILL",
//	LIMIT:         "LIMIT",
//	MEASUREMENT:   "MEASUREMENT",
//	MEASUREMENTS:  "MEASUREMENTS",
//	NAME:          "NAME",
//	OFFSET:        "OFFSET",
//	ON:            "ON",
//	ORDER:         "ORDER",
//	PASSWORD:      "PASSWORD",
//	POLICY:        "POLICY",
//	POLICIES:      "POLICIES",
//	PRIVILEGES:    "PRIVILEGES",
//	QUERIES:       "QUERIES",
//	QUERY:         "QUERY",
//	READ:          "READ",
//	REPLICATION:   "REPLICATION",
//	RESAMPLE:      "RESAMPLE",
//	RETENTION:     "RETENTION",
//	REVOKE:        "REVOKE",
//	SELECT:        "SELECT",
//	SERIES:        "SERIES",
//	SET:           "SET",
//	SHOW:          "SHOW",
//	SHARD:         "SHARD",
//	SHARDS:        "SHARDS",
//	SLIMIT:        "SLIMIT",
//	SOFFSET:       "SOFFSET",
//	STATS:         "STATS",
//	SUBSCRIPTION:  "SUBSCRIPTION",
//	SUBSCRIPTIONS: "SUBSCRIPTIONS",
//	TAG:           "TAG",
//	TO:            "TO",
//	USER:          "USER",
//	USERS:         "USERS",
//	VALUES:        "VALUES",
//	WHERE:         "WHERE",
//	WITH:          "WITH",
//	WRITE:         "WRITE",
//}

var keywords = map[string]int{
	"select":	SELECT,
	"as":		AS,
	"from":		FROM,
	"group":	GROUP,
	"by":		BY,
	"where":	WHERE,
	"or":		OR,
	"and":		AND,
	"order":	ORDER,
	"desc":		DESC,
	"asc":		ASC,
	"limit":	LIMIT,
	"show":		SHOW,
	"databases":	DATABASES,
	"database":	DATABASE,
	"create":	CREATE,
	"measurements":	MEASUREMENTS,
}

// Lookup returns the token associated with a given string.
func Lookup(ident string) int {
	if tok, ok := keywords[strings.ToLower(ident)]; ok {
		return tok
	}
	return IDENT
}

// Pos specifies the line and character position of a token.
// The Char and Line are both zero-based indexes.
type Pos struct {
	Line int
	Char int
}
