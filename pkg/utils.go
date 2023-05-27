package pkg

import "strings"

type Statements struct {
	instructions []string
}

type Pair struct {
	first string
	second string
}

type TrackedCommand struct {
	command string
	idx int
}

func (statement *Statements) AddStatement(chunks ...string) {
	statement.instructions = append(statement.instructions, strings.Join(chunks, " "))
}

func wrapInQuotes(s string) string {
	return "\"" + s +"\""
}

func map2(data []string, f func(string) string) []string {

    mapped := make([]string, len(data))

    for i, e := range data {
        mapped[i] = f(e)
    }

    return mapped
}