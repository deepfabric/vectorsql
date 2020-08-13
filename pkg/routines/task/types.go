package task

type Task interface {
	Stop(TaskResult)
	Execute() TaskResult
}

type TaskResult interface {
	Error() error
	Result() interface{}
}
