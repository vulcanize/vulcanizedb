package ipfs

type DagPutter interface {
	DagPut(raw interface{}) ([]string, error)
}
