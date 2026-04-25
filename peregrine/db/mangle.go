package db

import (
	"fmt"
	"strings"

	"github.com/google/mangle/ast"
	"github.com/google/mangle/parse"
)

// ReasonUsingMangle demonstrates basic Datalog-like logic evaluation.
// This allows the agent to run logical reasoning queries natively over its context data.
func ReasonUsingMangle(query string) string {
	clause, err := parse.Clause(query)
	if err != nil {
		return fmt.Sprintf("Mangle Parse Error: %v", err)
	}

	// Just a simple echo of the AST to show Mangle is successfully hooked up.
	// In production, we'd feed this into a mangle evaluator with local facts.
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Mangle Parsed Predicate: %v\n", clause.Head.Predicate))
	for _, arg := range clause.Head.Args {
		switch a := arg.(type) {
		case ast.Constant:
			sb.WriteString(fmt.Sprintf(" - Constant: %v\n", a.Symbol))
		case ast.Variable:
			sb.WriteString(fmt.Sprintf(" - Variable: %v\n", a.Symbol))
		}
	}

	return sb.String()
}
