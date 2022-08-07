package analyse

type analyse interface {
	GetValue(rule string) (string, error)
}
