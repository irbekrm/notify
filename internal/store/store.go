package store

import "sync"

type Storer interface{}

var dbConnPool sync.Pool
