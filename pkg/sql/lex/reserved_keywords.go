package lex

// GetKeywordID returns the lex id of the SQL keyword k or IDENT if k is
// not a keyword.
func GetKeywordID(k string) int32 {
	// The previous implementation generated a map that did a string ->
	// id lookup. Various ideas were benchmarked and the implementation below
	// was the fastest of those, between 3% and 10% faster (at parsing, so the
	// scanning speedup is even more) than the map implementation.
	switch k {
	case "all":
		return ALL
	case "and":
		return AND
	case "as":
		return AS
	case "asc":
		return ASC
	case "between":
		return BETWEEN
	case "bool":
		return BOOL
	case "by":
		return BY
	case "cast":
		return CAST
	case "cross":
		return CROSS
	case "desc":
		return DESC
	case "distinct":
		return DISTINCT
	case "exists":
		return EXISTS
	case "ftop":
		return FTOP
	case "false":
		return FALSE
	case "fetch":
		return FETCH
	case "first":
		return FIRST
	case "float":
		return FLOAT
	case "from":
		return FROM
	case "full":
		return FULL
	case "group":
		return GROUP
	case "having":
		return HAVING
	case "inner":
		return INNER
	case "int":
		return INT
	case "intersect":
		return INTERSECT
	case "is":
		return IS
	case "except":
		return EXCEPT
	case "join":
		return JOIN
	case "natural":
		return NATURAL
	case "next":
		return NEXT
	case "not":
		return NOT
	case "null":
		return NULL
	case "offset":
		return OFFSET
	case "on":
		return ON
	case "only":
		return ONLY
	case "or":
		return OR
	case "order":
		return ORDER
	case "outer":
		return OUTER
	case "right":
		return RIGHT
	case "row":
		return ROW
	case "rows":
		return ROWS
	case "select":
		return SELECT
	case "string":
		return STRING
	case "top":
		return TOP
	case "time":
		return TIME
	case "true":
		return TRUE
	case "union":
		return UNION
	case "where":
		return WHERE
	default:
		return IDENT
	}
}
