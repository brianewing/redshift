package opc

type AppType int

const (
	RedshiftApp AppType = iota
)

type ClientInfo struct {
	AppType
	AppVersionMajor int
	AppVersionMinor int
	DeviceName      string
}