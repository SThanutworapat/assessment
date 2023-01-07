package expense

import "database/sql"

type (
	DatabaseModelImpl interface {
		FindAll() (*sql.Stmt, error)
		FindByID() (*sql.Stmt, error)
		UpdateByID() (*sql.Stmt, error)
		CreateByID() (*sql.Stmt, error)
	}

	Expenses struct {
		ID     int      `json:"id"`
		Title  string   `json:"title"`
		Amount float64  `json:"amount"`
		Note   string   `json:"note"`
		Tags   []string `json:"tags"`
	}

	Err struct {
		Message string `json:"message"`
	}

	DatabaseModel struct {
		db *sql.DB
	}
)

func NewDatabaseModel(db *sql.DB) *DatabaseModel {
	return &DatabaseModel{
		db: db,
	}
}

func (d *DatabaseModel) FindAll() (*sql.Stmt, error) {
	return d.db.Prepare("SELECT id, title , amount , note , tags From expenses order by id asc")
}

func (d *DatabaseModel) FindByID() (*sql.Stmt, error) {
	return d.db.Prepare("SELECT id, title , amount , note , tags From expenses where id=$1")
}
func (d *DatabaseModel) UpdateByID() (*sql.Stmt, error) {
	return d.db.Prepare("update expenses set title = $1 , amount =  $2, note = $3, tags = $4 where id = $5 RETURNING id")
}
func (d *DatabaseModel) CreateByID() (*sql.Stmt, error) {
	return d.db.Prepare("INSERT INTO expenses ( title , amount , note , tags) values ($1, $2, $3, $4) RETURNING id")
}
