package repositories

const (
	blocksFromHeadBeforeFinal = 20
)

type Repository interface {
	BlockRepository
	ContractRepository
	LogsRepository
	ReceiptRepository
	FilterRepository
}
