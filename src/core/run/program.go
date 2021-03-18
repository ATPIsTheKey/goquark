package run

import (
	"goquark/src/core/ast"
	"goquark/src/core/lexer"
	"goquark/src/core/parser"
	"goquark/src/core/runtime"
	"io/ioutil"
	"os"
)

type Loader struct {
	Lexer        *lexer.Lexer
	Parser       *parser.Parser
	ProgramStmts []ast.Stmt
}

func newLoader() *Loader {
	return &Loader{
		Lexer:  &lexer.Lexer{},
		Parser: &parser.Parser{},
	}
}

func (loader *Loader) LoadProgram(fpath string) error {
	fp, err := os.Open(fpath)
	if err != nil {
		return nil
	} // todo
	src, err := ioutil.ReadAll(fp)
	if err != nil {
		return nil
	} // todo
	loader.Lexer.Init(src, lexer.BaseErrHandler, lexer.IgnoreSkippables)
	loader.Parser.Init(loader.Lexer.GetTokens(), parser.BaseErrHandler, parser.NoOpts)
	loader.ProgramStmts = loader.Parser.GetProgramStmts()
	return nil // todo
}

func (loader *Loader) Execute() []runtime.Object {
	return nil
}
