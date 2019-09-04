package ast

import (
	"Interpreter/token"
	"testing"
)

func TestProgram_String(t *testing.T) {

	program := &Program{
		Statements:[]Statement{
			&LetStatement{
				//let myVar = anotherVar;
				Token:token.Token{Type:token.LET, Literal:"let"},
				Name:&Identifier{Token:token.Token{Type:token.IDENT,Literal:"myVar"},Value:"myVar"},
				Value:&Identifier{Token:token.Token{Type:token.IDENT,Literal:"anotherVar"},Value:"anotherVar"},
			},
			&LetStatement{
				//let myVar = anotherVar;
				Token:token.Token{Type:token.LET, Literal:"let"},
				Name:&Identifier{Token:token.Token{Type:token.IDENT,Literal:"myVar2"},Value:"myVar2"},
				Value:&Identifier{Token:token.Token{Type:token.INT,Literal:"123"},Value:"123"},
			},
		},
	}
	if program.String() != "let myVar = anotherVar;let myVar2 = 123;" {
		t.Errorf("program.String() wrong. got=%q", program.String())
	}
}
