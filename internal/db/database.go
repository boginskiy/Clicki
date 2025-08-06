package db

const CAP = 10

var Store = make(map[string]string, CAP)
