package tran

import "test.com/project-project/internal/database"

// Transaction 事务操作，一定和数据库有关；
type Transaction interface {
	Action(func(conn database.DbConn) error) error
}
