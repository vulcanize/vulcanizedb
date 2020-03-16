package core

type BatchElem struct {
	Method string
	Args   []interface{}
	Result interface{}
	Error  error
}
