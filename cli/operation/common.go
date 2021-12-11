package operation

type OperationType int

const (
	SET      = OperationType(0)
	INCREASE = OperationType(1)
	DECREASE = OperationType(-1)
)
