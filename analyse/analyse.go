package analyse

type analyse interface {
	getValue(rule string) (string, error)
}
