package model

//请求
type Request interface {
	Resource() Resource
	Action() Action
}
