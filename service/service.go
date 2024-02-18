package service

type Database interface {
	Backup()
}
type callback func(data string)
