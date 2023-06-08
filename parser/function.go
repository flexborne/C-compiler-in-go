package parser

type ReturnType int

const (
	INT ReturnType = iota
	VOID
)

type StatementType int

const (
	FUNC_CALL StatementType = iota
	RETURN
)

type FuncCallStatement struct {
	Name string
	Args []string
}

type ReturnStatement struct {
	Expr string
}

type Function struct {
	Name string
	Body []any
}
