package client

type ConnStatus int

const (
	ConnConnecting = iota
	ConnReconnecting
	ConnOnline
)
