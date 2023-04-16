package statement

type BlockStmt struct {
	Statements []Stmt
}

func NewBlockStmt(statements []Stmt) *BlockStmt {
	return &BlockStmt{
		Statements: statements,
	}
}

func (ps *BlockStmt) Stmt() {}
