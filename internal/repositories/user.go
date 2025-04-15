package repositories

import "database/sql"

type UserRepository struct {
	Conn *sql.DB
}

func (u *UserRepository) GetUserById(id int) {

}
