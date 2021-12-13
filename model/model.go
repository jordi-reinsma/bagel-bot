package model

type User struct {
	ID   int
	UUID string
}

type Match struct {
	ID   int
	A, B User
	Freq int
}
