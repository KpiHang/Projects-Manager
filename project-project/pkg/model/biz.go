package model

var (
	Normal         = 1 // 不是1说明，表中的一条数据被禁用了
	Personal int32 = 1 // 1就表示个人项目 自己本来就属于一个组织；
)

var AESKey = "sdfgyrhgbxcdgryfhgywertd"

const (
	NoDeleted = iota
	Deleted
)

const (
	NoArchive = iota
	Archive
)

const (
	Open = iota
	Private
)

const (
	Default = "default"
	Simple  = "simple"
)

const (
	NoCollected = iota
	Collected
)
