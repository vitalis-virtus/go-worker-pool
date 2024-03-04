package job

type Result struct {
	Value    interface{}
	Error    error
	Metadata map[string]interface{}
}
