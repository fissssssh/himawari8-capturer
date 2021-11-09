package core

type Quality uint8

const (
	Low Quality = 1 << iota
	HD
	FHD
	QHD
	UHD
	UHDPlus Quality = 20
)
