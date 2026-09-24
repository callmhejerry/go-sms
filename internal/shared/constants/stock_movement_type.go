package constants

type StockMovementType string

const (
	In     StockMovementType = "in"
	Out    StockMovementType = "out"
	Adjust StockMovementType = "adjust"
)
