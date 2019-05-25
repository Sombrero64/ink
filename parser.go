package main

import (
	"fmt"
	"strconv"
)

func Parse(tokenStream <-chan Tok, nodes chan<- interface{}, done chan<- bool) {
	tokens := make([]Tok, 0)
	for tok := range tokenStream {
		fmt.Println("Token: " + tokKindToName(tok.kind) + " | " + tok.stringVal())
		tokens = append(tokens, tok)
	}

	idx, length := 0, len(tokens)
	fmt.Println("tokens length", strconv.Itoa(length))
	for idx < length && idx != -1 {
		var expr interface{}
		expr, idx = parseExpression(tokens[idx:])
		nodes <- expr
	}
	close(nodes)

	done <- true
}

type UnaryExprNode struct {
	operator Tok
	operand  interface{}
}

type BinaryExprNode struct {
	operator     Tok
	leftOperand  interface{}
	rightOperand interface{}
}

type FunctionCallExprNode struct {
	function  interface{}
	arguments []interface{}
}

type MatchClauseNode struct {
	target      interface{}
	expressions []interface{}
}

type MatchExprNode struct {
	condition interface{}
	clauses   []MatchClauseNode
}

func parseExpression(tokens []Tok) (interface{}, int) {
	tok := tokens[0]
	if tok.kind == NegationOp {
		atom, idx := parseAtom(tokens[1:])
		return UnaryExprNode{
			tok,
			atom,
		}, idx
	}

	atom, idx := parseAtom(tokens)

	next := tokens[idx]
	switch next.kind {
	case Separator:
		return atom, idx
	case AddOp, SubtractOp, MultiplyOp, DivideOp, ModulusOp,
		GreaterThanOp, LessThanOp, EqualOp, IsOp, DefineOp, AccessorOp:
		var rightOperand interface{}
		rightOperand, idx = parseAtom(tokens[idx:])
		return BinaryExprNode{
			next,
			atom,
			rightOperand,
		}, idx
	case LeftParen:
		arguments := make([]interface{}, 0)
		for tokens[idx].kind != RightParen {
			var exprNode interface{}
			exprNode, idx = parseExpression(tokens[idx:])
			arguments = append(arguments, exprNode)
		}
		return FunctionCallExprNode{
			atom,
			arguments,
		}, idx
	case MatchColon:
		idx++ // assume next token is RightBrace for now
		clauses := make([]MatchClauseNode, 0)
		for tokens[idx].kind != RightBrace {
			var clauseNode MatchClauseNode
			clauseNode, idx = parseMatchClause(tokens[idx:])
			clauses = append(clauses, clauseNode)
		}
		return MatchExprNode{
			atom,
			clauses,
		}, idx
	default:
		// error -- should not happen
		return []interface{}{}, -1
	}

	return []interface{}{}, -1
}

type IdentifierNode struct {
	val string
}

type NumberLiteralNode struct {
	val float64
}

type StringLiteralNode struct {
	val string
}

type ObjectLiteralNode struct {
}

type ListLiteralNode struct {
}

func parseAtom(tokens []Tok) (interface{}, int) {
	tok := tokens[0]
	switch tok.kind {
	case Identifier:
		return IdentifierNode{tok.stringVal()}, len(tok.stringVal())
	case NumberLiteral:
		return NumberLiteralNode{tok.numberVal()}, len(tok.stringVal())
	case StringLiteral:
		return StringLiteralNode{tok.stringVal()}, len(tok.stringVal())
	case TrueLiteral, FalseLiteral, NullLiteral:
		return tok, 0
	case LeftParen:
		// grouped expression
	case LeftBrace:
		// object literal
	case LeftBracket:
		// array literal
	default:
		// may be function literal
	}
	return IdentifierNode{}, 0
}

func parseMatchClause(tokens []Tok) (MatchClauseNode, int) {
	return MatchClauseNode{}, 0
}

func parseObjectLiteral(tokens []Tok) (interface{}, int) {
	return ObjectLiteralNode{}, 0
}

func parseListLiteral(tokens []Tok) (interface{}, int) {
	return ListLiteralNode{}, 0
}

func parseFunctionLiteral(tokens []Tok) (FunctionCallExprNode, int) {
	return FunctionCallExprNode{}, 0
}
