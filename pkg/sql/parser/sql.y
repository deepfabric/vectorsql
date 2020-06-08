%{

package parser

import (
    "github.com/deepfabric/vectorsql/pkg/sql/tree"
	"github.com/deepfabric/vectorsql/pkg/vm/value"
)

%}

%{

type sqlSymUnion struct {
    val interface{}
}

func (u *sqlSymUnion) bool() bool {
    return u.val.(bool)
}

func (u *sqlSymUnion) isNull() bool {
    return u.val == nil
}

func (u *sqlSymUnion) setNegative() *tree.Value{
    v, ok := u.val.(*tree.Value)
    if !ok {
        return nil
    }
    i, ok := value.AsInt(v.E)
    if !ok {
        return nil
    }
    return &tree.Value{value.NewInt(int64(i)*-1)}
}

func (u *sqlSymUnion) valueStatement() *tree.Value {
    return u.val.(*tree.Value)
}

func (u *sqlSymUnion) selectStatement() *tree.Select {
    return u.val.(*tree.Select)
}

func (u *sqlSymUnion) limitStatement() *tree.Limit {
    if u.val == nil {
        return nil
    }
    return u.val.(*tree.Limit)
}

func (u *sqlSymUnion) orderTopStatement() tree.OrderStatement{
    if u.val == nil{
        return nil
    }
    return u.val.(tree.OrderStatement)
}

func (u *sqlSymUnion) orderByStatement() tree.OrderBy {
    if u.val == nil {
        return nil
    }
    return u.val.(tree.OrderBy)
}

func (u *sqlSymUnion) relationStatement() tree.RelationStatement {
    if u.val == nil {
        return nil
    }
    return u.val.(tree.RelationStatement)
}

func (u *sqlSymUnion) joinStatement() *tree.JoinClause {
    return u.val.(*tree.JoinClause)
}

func (u *sqlSymUnion) unionStatement() *tree.UnionClause {
    return u.val.(*tree.UnionClause)
}

func (u *sqlSymUnion) simpleSelectStatement() *tree.SelectClause {
    return u.val.(*tree.SelectClause)
}

func (u *sqlSymUnion) fromStatement() *tree.From {
    if u.val == nil {
        return nil
    }
    return u.val.(*tree.From)
}

func (u *sqlSymUnion) groupByStatement() *tree.GroupBy {
    if u.val == nil {
        return nil
    }
    return u.val.(*tree.GroupBy)
}

func (u *sqlSymUnion) whereStatement() *tree.Where {
    if u.val == nil {
        return nil
    }
    return u.val.(*tree.Where)
}

func (u *sqlSymUnion) subqueryStatement() *tree.Subquery {
    return u.val.(*tree.Subquery)
}

func (u *sqlSymUnion) funcStatement() *tree.FuncExpr {
    return u.val.(*tree.FuncExpr)
}

func (u *sqlSymUnion) selectExprs() tree.SelectExprs {
    return u.val.(tree.SelectExprs)
}

func (u *sqlSymUnion) selectExpr() *tree.SelectExpr {
    if u.val == nil {
        return nil
    }
    return u.val.(*tree.SelectExpr)
}

func (u *sqlSymUnion) tableStatements() tree.TableStatements {
    return u.val.(tree.TableStatements)
}

func (u *sqlSymUnion) tableStatement() tree.TableStatement {
    return u.val.(tree.TableStatement)
}

func (u *sqlSymUnion) limit() *tree.Limit {
    return u.val.(*tree.Limit)
}

func (u *sqlSymUnion) orderStatement() *tree.Order{
    return u.val.(*tree.Order)
}

func (u *sqlSymUnion) direction() tree.Direction{
    return u.val.(tree.Direction)
}

func (u *sqlSymUnion) joinType() tree.JoinType {
    return u.val.(tree.JoinType)
}

func (u *sqlSymUnion) joinCond() tree.JoinCond {
    return u.val.(tree.JoinCond)
}

func (u *sqlSymUnion) tableName() *tree.TableName {
    return u.val.(*tree.TableName)
}

func (u *sqlSymUnion) nameList() tree.NameList {
    return u.val.(tree.NameList)
}

func (u *sqlSymUnion) exprStatement() tree.ExprStatement {
    return u.val.(tree.ExprStatement)
}

func (u *sqlSymUnion) exprStatements() tree.ExprStatements {
    return u.val.(tree.ExprStatements)
}

func (u *sqlSymUnion) colunmNameList() tree.ColunmNameList {
    return u.val.(tree.ColunmNameList)
}

func (u *sqlSymUnion) aliasClause() *tree.AliasClause {
    if u.val == nil{
        return nil
    }
    return u.val.(*tree.AliasClause)
}

%}

%token <str> IDENT
%token <union> ICONST FCONST SCONST
%token <str> LESS_EQUALS GREATER_EQUALS NOT_EQUALS

%token <str> ALL AND AS ASC

%token <str> BETWEEN
%token <str> BOOL BY

%token <str> CAST
%token <str> CROSS

%token <str> DESC
%token <str> DISTINCT

%token <str> EXCEPT
%token <str> EXISTS

%token <str> FTOP
%token <str> FALSE FETCH
%token <str> FIRST FLOAT FROM FULL

%token <str> GROUP

%token <str> HAVING

%token <str> INNER INT
%token <str> INTERSECT IS

%token <str> JOIN

%token <str> NATURAL NEXT
%token <str> NOT NULL

%token <str> OFFSET ON ONLY OR
%token <str> ORDER OUTER

%token <str> RIGHT
%token <str> ROW ROWS

%token <str> SELECT
%token <str> STRING

%token <str> TOP
%token <str> TIME TRUE

%token <str> UNION

%token <str> WHERE

%token <str> NOT_LA

%union {
    id      int32
    pos     int32
    byt     byte
    str     string
    union   sqlSymUnion
}

%type <union> stmt_block
%type <union> stmt

%type <union> select_stmt

%type <union> relation

%type <union> join_clause union_clause select_clause

%type <union> simple_select

%type <union> subquery

%type <union> opt_asc_desc

%type <str> name
%type <str> func_name
%type <str> table_alias_name target_name

%type <union> table_name
%type <union> column_name

%type <union> distinct_clause
%type <union> opt_column_list
%type <union> order_clause opt_order_clause
%type <union> order_list
%type <union> name_list
%type <union> from_clause
%type <union> from_list
%type <union> expr_list
%type <union> target_list
%type <union> group_clause
%type <union> fetch_clause opt_fetch_clause

%type <union> all_or_distinct
%type <empty> join_outer
%type <union> join_qual
%type <union> join_type

%type <union> limit_clause offset_clause
%type <union> opt_select_fetch_first_value
%type <empty> row_or_rows
%type <empty> first_or_next

%type <union> where_clause opt_where_clause
%type <union> a_expr b_expr c_expr d_expr
%type <union> having_clause
%type <union> alias_clause opt_alias_clause
%type <union> order
%type <union> table_ref
%type <union> target_elem

%type <union> typename

%type <union> cast_target

%type <union> signed_iconst

%type <union> func_application func_expr_common_subexpr
%type <union> func_expr

%type <byt> '+' '-' '*' '/' '%' '<' '>' '=' '[' ']' '(' ')' '.'

%left      UNION EXCEPT
%left      INTERSECT
%left      OR
%left      AND
%right     NOT
%nonassoc  IS
%nonassoc  '<' '>' '=' LESS_EQUALS GREATER_EQUALS NOT_EQUALS
%nonassoc  BETWEEN NOT_LA

%nonassoc  IDENT NULL ROWS
%left      '+' '-'
%left      '*' '/' '%'

%left      AT
%right     UMINUS
%left      '[' ']'
%left      '(' ')'
%left      '.'

%left      JOIN CROSS LEFT FULL RIGHT INNER NATURAL

%%

stmt_block: stmt    { sqllex.(*lexer).SetStmt($1.selectStatement()) }

stmt: select_stmt   { $$.val = $1.selectStatement() }

select_stmt: relation opt_order_clause opt_fetch_clause
             {
                $$.val = &tree.Select{
                    Limit:      $3.limitStatement(),
                    Order:      $2.orderTopStatement(),
                    Relation:   $1.relationStatement(),
                }
             }



opt_order_clause: order_clause { $$.val = $1.orderTopStatement() }
                |              { $$.val = nil }

order_clause: ORDER BY order_list   { $$.val = $3.orderByStatement() }
            | TOP a_expr            { $$.val = &tree.Top{
                                        N: $2.exprStatement(),
                                      }
                                    }
            | FTOP a_expr           { $$.val = &tree.Ftop{
                                        N: $2.exprStatement(),
                                      }
                                    }

order_list: order { $$.val = tree.OrderBy{$1.orderStatement()}}
            | order_list ',' order  { $$.val = append($1.orderByStatement(), $3.orderStatement())}

order: a_expr opt_asc_desc
       {
            $$.val = &tree.Order{
                E:          $1.exprStatement(),
                Type:       $2.direction(),
            }
       }

opt_asc_desc: ASC   { $$.val = tree.Ascending }
            | DESC  { $$.val = tree.Descending }
            |       { $$.val = tree.DefaultDirection }


opt_fetch_clause: fetch_clause  { $$.val = $1.limitStatement() }
                |               { $$.val = nil }

fetch_clause: limit_clause offset_clause
              {
                if $1.limitStatement() == nil {
                    $$.val = $2.limitStatement()
                }else{
                    $$.val = $1.limitStatement()
                    $$.val.(*tree.Limit).Offset = $2.limitStatement().Offset
                }
              }
            | offset_clause limit_clause
              {
                $$.val = $1.limitStatement()
                if $2.limitStatement() != nil {
                    $$.val.(*tree.Limit).Count = $2.limitStatement().Count
                }
              }
            | limit_clause
              {
                $$.val = $1.limitStatement()
              }
            | offset_clause
              {
                $$.val = $1.limitStatement()
              }

limit_clause: FETCH first_or_next opt_select_fetch_first_value row_or_rows ONLY
              {
                $$.val = &tree.Limit{Count: $3.exprStatement()}
              }

offset_clause: OFFSET a_expr                { $$.val = &tree.Limit{Offset: $2.exprStatement()} }
             | OFFSET d_expr row_or_rows    { $$.val = &tree.Limit{Offset: $2.exprStatement()} }

opt_select_fetch_first_value: signed_iconst     { $$.val = $1.exprStatement() }
                            | '(' a_expr ')'    { $$.val = $2.exprStatement() }
                            |                   { $$.val = &tree.Value{value.NewInt(1)} }

row_or_rows: ROW    {}
           | ROWS   {}

first_or_next: FIRST    {}
             | NEXT     {}



relation: table_name  opt_alias_clause              { $$.val = &tree.AliasedTable{
                                                                    As: $2.aliasClause(),
                                                                    Tbl: $1.tableName(),
                                                               }
                                                    }
        | join_clause                               { $$.val = $1.joinStatement() }
        | union_clause                              { $$.val = $1.unionStatement() }
        | simple_select                             { $$.val = $1.simpleSelectStatement() }
        | '(' select_stmt ')' opt_alias_clause      { $$.val = &tree.AliasedSelect{
                                                                    As:     $4.aliasClause(),
                                                                    Sel:    $2.selectStatement(),
                                                                }
                                                    }


simple_select: SELECT target_list from_clause opt_where_clause group_clause having_clause
               {
                    $$.val = &tree.SelectClause{
                                Distinct: false,
                                Sel:     $2.selectExprs(),
                                From:    $3.fromStatement(),
                                Where:   $4.whereStatement(),
                                Having:  $6.whereStatement(),
                                GroupBy: $5.groupByStatement(),
                            }
               }
             | SELECT distinct_clause target_list from_clause opt_where_clause group_clause having_clause
               {
                    $$.val = &tree.SelectClause{
                                Distinct: $2.bool(),
                                Sel:      $3.selectExprs(),
                                From:     $4.fromStatement(),
                                Where:    $5.whereStatement(),
                                Having:   $7.whereStatement(),
                                GroupBy:  $6.groupByStatement(),
                            }
               }



distinct_clause: DISTINCT { $$.val = true }



target_list: target_elem
             {
                if $1.isNull() {
                    $$.val = tree.SelectExprs{}
                }else{
                    $$.val = tree.SelectExprs{$1.selectExpr()}
                }
             }
           | target_list ',' target_elem
             {
                if $3.isNull() {
                    $$.val = $1.selectExprs()
                }else{
                    $$.val = append($1.selectExprs(), $3.selectExpr())
                }
             }

target_elem: a_expr
             {
                $$.val = &tree.SelectExpr{E: $1.exprStatement()}
             }
           | a_expr target_name
             {
                $$.val = &tree.SelectExpr{E: $1.exprStatement(), As: tree.Name($2)}
             }
           | a_expr AS target_name
             {
                $$.val = &tree.SelectExpr{E: $1.exprStatement(), As: tree.Name($3)}
             }
           | '*'
             {
                $$.val = nil
             }



from_clause: FROM from_list
             {
                $$.val = &tree.From{$2.tableStatements()}
             }
           | { $$.val = nil }

from_list: table_ref
           {
                $$.val = tree.TableStatements{$1.tableStatement()}
           }
         | from_list ',' table_ref
           {
                $$.val = append($1.tableStatements(), $3.tableStatement())
           }



opt_where_clause: where_clause
                 {
                    $$.val = &tree.Where{ Type: tree.AstWhere, E: $1.exprStatement() }
                 }
                | { $$.val = nil }

where_clause: WHERE a_expr { $$.val = $2.exprStatement() }



group_clause: GROUP BY expr_list    { $$.val = &tree.GroupBy{$3.exprStatements()} }
            |                       { $$.val = nil }



having_clause: HAVING a_expr
               {
                    $$.val = &tree.Where{ Type: tree.AstHaving, E: $2.exprStatement() }
               }
             | { $$.val = nil }



expr_list: a_expr               { $$.val = tree.ExprStatements{$1.exprStatement()} }
         | expr_list ',' a_expr { $$.val = append($1.exprStatements(), $3.exprStatement()) }

a_expr: c_expr                  { $$.val = $1.exprStatement() }
      | NOT a_expr              { $$.val = &tree.NotExpr{E: $2.exprStatement()} }
      | a_expr OR a_expr        { $$.val = &tree.OrExpr{Left: $1.exprStatement(), Right: $3.exprStatement()} }
      | a_expr AND a_expr       { $$.val = &tree.AndExpr{Left: $1.exprStatement(), Right: $3.exprStatement()} }
      | a_expr IS NULL          { $$.val = &tree.IsNullExpr{E: $1.exprStatement()} }
      | a_expr IS NOT NULL      { $$.val = &tree.IsNotNullExpr{E: $1.exprStatement()} }
      | b_expr                  { $$.val = $1.exprStatement() }

b_expr: d_expr                  { $$.val = $1.exprStatement() }
      | column_name             { $$.val = $1.colunmNameList() }
      | '+' b_expr              { $$.val = $2.exprStatement() }
      | '-' b_expr              { $$.val = &tree.UnaryMinusExpr{E: $2.exprStatement()} }
      | b_expr '+' b_expr       { $$.val = &tree.PlusExpr{Left: $1.exprStatement(), Right: $3.exprStatement()} }
      | b_expr '-' b_expr       { $$.val = &tree.MinusExpr{Left: $1.exprStatement(), Right: $3.exprStatement()} }
      | b_expr '*' b_expr       { $$.val = &tree.MultExpr{Left: $1.exprStatement(), Right: $3.exprStatement()} }
      | b_expr '/' b_expr       { $$.val = &tree.DivExpr{Left: $1.exprStatement(), Right: $3.exprStatement()} }
      | b_expr '%' b_expr       { $$.val = &tree.ModExpr{Left: $1.exprStatement(), Right: $3.exprStatement()} }
      | func_expr               { $$.val = $1.funcStatement() }

c_expr: b_expr '<' b_expr                           { $$.val = &tree.LtExpr{Left: $1.exprStatement(), Right: $3.exprStatement()} }
      | b_expr '>' b_expr                           { $$.val = &tree.GtExpr{Left: $1.exprStatement(), Right: $3.exprStatement()} }
      | b_expr '=' b_expr                           { $$.val = &tree.EqExpr{Left: $1.exprStatement(), Right: $3.exprStatement()} }
      | b_expr LESS_EQUALS b_expr                   { $$.val = &tree.LeExpr{Left: $1.exprStatement(), Right: $3.exprStatement()} }
      | b_expr GREATER_EQUALS b_expr                { $$.val = &tree.GeExpr{Left: $1.exprStatement(), Right: $3.exprStatement()} }
      | b_expr NOT_EQUALS b_expr                    { $$.val = &tree.NeExpr{Left: $1.exprStatement(), Right: $3.exprStatement()} }
      | b_expr BETWEEN b_expr AND b_expr            { $$.val = &tree.BetweenExpr{E: $1.exprStatement(), From: $3.exprStatement(), To: $5.exprStatement()} }
      | b_expr NOT_LA BETWEEN b_expr AND b_expr     { $$.val = &tree.NotBetweenExpr{E: $1.exprStatement(), From: $4.exprStatement(), To: $6.exprStatement()} }
      | EXISTS subquery                             {
                                                        $$.val = $2.subqueryStatement()
                                                        $$.val.(*tree.Subquery).Exists = true
                                                    }

d_expr: ICONST          { $$.val = $1.valueStatement() }
      | FCONST          { $$.val = $1.valueStatement() }
      | SCONST          { $$.val = $1.valueStatement() }
      | TRUE            { $$.val = &tree.Value{&value.ConstTrue} }
      | FALSE           { $$.val = &tree.Value{&value.ConstFalse} }
      | NULL            { $$.val = &tree.Value{&value.ConstNull} }
      | '(' a_expr ')'  { $$.val = &tree.ParenExpr{$2.exprStatement()} }

signed_iconst: ICONST       { $$.val = $1.valueStatement() }
             | '+' ICONST   { $$.val = $2.valueStatement() }
             | '-' ICONST   { $$.val = $2.setNegative() }



func_expr: func_application
           {
                $$.val = $1.funcStatement()
           }
         | func_expr_common_subexpr
           {
                $$.val = $1.funcStatement()
           }

func_application: func_name '(' ')'
                  {
                    $$.val = &tree.FuncExpr{Name: $1}
                  }
                | func_name '(' expr_list ')'
                  {
                    $$.val = &tree.FuncExpr{Name: $1, Es: $3.exprStatements() }
                  }

func_expr_common_subexpr: CAST '(' a_expr AS cast_target ')'
                          {
                            $$.val = &tree.FuncExpr{Name: "cast", Es: tree.ExprStatements{$3.exprStatement(), $5.exprStatement()} }
                          }

cast_target: typename { $$.val = $1.exprStatement() }

typename: INT       { $$.val = &tree.Value{value.NewString("int")} }
        | BOOL      { $$.val = &tree.Value{value.NewString("bool")} }
        | TIME      { $$.val = &tree.Value{value.NewString("time")} }
        | FLOAT     { $$.val = &tree.Value{value.NewString("float")} }
        | STRING    { $$.val = &tree.Value{value.NewString("string")} }



alias_clause: AS table_alias_name opt_column_list
              {
                $$.val = &tree.AliasClause{Alias: tree.Name($2), Cols: $3.nameList()}
              }
            | table_alias_name opt_column_list
              {
                $$.val = &tree.AliasClause{Alias: tree.Name($1), Cols: $2.nameList()}
              }

opt_alias_clause: alias_clause  { $$.val = $1.aliasClause() }
                |               { $$.val = nil }



subquery: '(' select_stmt ')'   { $$.val = &tree.Subquery{ Select: $2.selectStatement(), Exists: false }}


select_clause: relation         { $$.val = $1.relationStatement() }



union_clause: select_clause UNION all_or_distinct select_clause
              {
                    $$.val = &tree.UnionClause{
                        Type:  tree.UnionOp,
                        Left:  $1.relationStatement(),
                        Right: $4.relationStatement(),
                        All:   $3.bool(),
                    }
              }
            | select_clause INTERSECT all_or_distinct select_clause
              {
                    $$.val = &tree.UnionClause{
                        Type:  tree.IntersectOp,
                        Left:  $1.relationStatement(),
                        Right: $4.relationStatement(),
                        All:   $3.bool(),
                    }
              }
            | select_clause EXCEPT all_or_distinct select_clause
              {
                    $$.val = &tree.UnionClause{
                        Type:  tree.ExceptOp,
                        Left:  $1.relationStatement(),
                        Right: $4.relationStatement(),
                        All:   $3.bool(),
                    }
              }

all_or_distinct: ALL        { $$.val = true }
               | DISTINCT   { $$.val = false }
               |            { $$.val = false }



join_clause: select_clause CROSS JOIN select_clause
             {
                $$.val = &tree.JoinClause{
                    Type: tree.CrossOp,
                    Cond: &tree.NonJoinCond{},
                    Left: $1.relationStatement(),
                    Right: $4.relationStatement(),
                }
             }
           | select_clause join_type JOIN select_clause join_qual
             {
                $$.val = &tree.JoinClause{
                    Type: $2.joinType(),
                    Cond: $5.joinCond(),
                    Left: $1.relationStatement(),
                    Right: $4.relationStatement(),
                }
             }
           | select_clause JOIN select_clause join_qual
             {
                $$.val = &tree.JoinClause{
                    Type: tree.InnerOp,
                    Cond: $4.joinCond(),
                    Left: $1.relationStatement(),
                    Right: $3.relationStatement(),
                }
             }
           | select_clause NATURAL JOIN select_clause
             {
                $$.val = &tree.JoinClause{
                    Type: tree.NaturalOp,
                    Cond: &tree.NonJoinCond{},
                    Left: $1.relationStatement(),
                    Right: $4.relationStatement(),
                }
             }

join_qual: ON a_expr        { $$.val = &tree.OnJoinCond{E: $2.exprStatement()} }

join_type: FULL join_outer  { $$.val = tree.FullOp }
         | LEFT join_outer  { $$.val = tree.LeftOp }
         | RIGHT join_outer { $$.val = tree.RightOp }
         | INNER            { $$.val = tree.InnerOp }

join_outer: OUTER {}
          |   {}



table_ref: table_name opt_alias_clause
           {
                $$.val = &tree.AliasedTable {
                    Tbl:    $1.tableName(),
                    As:     $2.aliasClause(),
                }
           }
         | subquery opt_alias_clause
           {
                $$.val = &tree.AliasedTable {
                    Tbl:    $1.subqueryStatement(),
                    As:     $2.aliasClause(),
                }
           }



table_name: column_name
            {
                $$.val = &tree.TableName{$1.colunmNameList()}
            }



column_name: name
            {
                $$.val = tree.ColunmNameList{tree.ColunmName{Path: tree.Name($1)}}
            }
           | name '[' a_expr ']'
            {
                $$.val = tree.ColunmNameList{tree.ColunmName{Path: tree.Name($1), Index: $3.exprStatement()}}
            }
           | column_name '.' name
            {
                $$.val = append($1.colunmNameList(), tree.ColunmName{Path: tree.Name($3)})
            }
           | column_name '.' name '[' a_expr ']'
            {
                $$.val = append($1.colunmNameList(), tree.ColunmName{Path: tree.Name($3), Index: $5.exprStatement()})
            }


opt_column_list: '(' name_list ')'  { $$.val = $2.nameList() }
               |                    { $$.val = tree.NameList(nil)}

name_list: name
           {
            $$.val = tree.NameList{tree.Name($1)}
           }
         | name_list ',' name
           {
            $$.val = append($1.nameList(), tree.Name($3))
           }



name: IDENT

func_name: name

target_name: name

table_alias_name: name

%%
