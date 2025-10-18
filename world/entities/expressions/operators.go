package expressions

type BinaryOp uint8

const (
	OpEq BinaryOp = iota
	OpNe
	OpGt
	OpGe
	OpLt
	OpLe
	OpAdd
	OpSub
	OpMul
	OpDiv
)

type UnaryOp uint8

const (
	UNot UnaryOp = iota
	UNeg
)
