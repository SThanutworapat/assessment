package expenses

import (
	"database/sql"
	"net/http"

	expense "github.com/SThanutworapat/assessment/models"
	echo "github.com/labstack/echo/v4"
	"github.com/lib/pq"
)

type handler struct {
	Database expense.DatabaseModelImpl
}

const (
	_ = iota
	GetAll
	GetById
	Create
	UpdateById
)

func NewHandler(e expense.DatabaseModelImpl) *handler {
	return &handler{e}
}

func (h *handler) GetExpensesHandler(c echo.Context) error {
	isAuth, returnError := CheckAuth(c)
	if isAuth {
		return returnError
	}
	stmt, isError := h.Prepare(GetAll, c)
	if isError {
		return c.JSON(http.StatusInternalServerError, expense.Err{Message: "Prepare Fail"})
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, expense.Err{Message: "Query Fail"})
	}

	expenseses := []expense.Expenses{}
	for rows.Next() {
		var e expense.Expenses
		err = rows.Scan(&e.ID, &e.Title, &e.Amount, &e.Note, pq.Array(&e.Tags))
		if err != nil {
			return c.JSON(http.StatusInternalServerError, expense.Err{Message: err.Error()})
		}
		expenseses = append(expenseses, e)
	}

	return c.JSON(http.StatusOK, expenseses)
}

func (h *handler) Prepare(cases int, c echo.Context) (*sql.Stmt, bool) {
	switch cases {
	case GetAll:
		stmt, err := h.Database.FindAll()
		if err != nil {
			return stmt, true
		}
		return stmt, false
	case GetById:
		stmt, err := h.Database.FindByID()
		if err != nil {
			return stmt, true
		}
		return stmt, false
	case Create:
		stmt, err := h.Database.CreateByID()
		if err != nil {
			return stmt, true
		}
		return stmt, false

	case UpdateById:
		stmt, err := h.Database.UpdateByID()
		if err != nil {
			return stmt, true
		}
		return stmt, false
	default:
		return nil, false
	}
}

func (h *handler) GetExpensesByIdHandler(c echo.Context) error {
	isAuth, returnError := CheckAuth(c)
	if isAuth {
		return returnError
	}
	id := c.Param("id")
	stmt, isError := h.Prepare(GetById, c)
	if isError {
		return c.JSON(http.StatusInternalServerError, expense.Err{Message: "Prepare Fail"})
	}
	defer stmt.Close()
	row := stmt.QueryRow(id)
	e := expense.Expenses{}
	err := row.Scan(&e.ID, &e.Title, &e.Amount, &e.Note, pq.Array(&e.Tags))
	switch err {
	case sql.ErrNoRows:
		return c.JSON(http.StatusNotFound, expense.Err{Message: "id Expenses not found"})
	case nil:
		return c.JSON(http.StatusOK, e)
	default:
		return c.JSON(http.StatusInternalServerError, expense.Err{Message: "Query Fail"})
	}
}
func (h *handler) PutExpensesByIdHandler(c echo.Context) error {
	isAuth, returnError := CheckAuth(c)
	if isAuth {
		return returnError
	}
	id := c.Param("id")
	var e expense.Expenses
	err := c.Bind(&e)
	if err != nil {
		return c.JSON(http.StatusBadRequest, expense.Err{Message: err.Error()})
	}
	stmt, isError := h.Prepare(UpdateById, c)
	if isError {
		return c.JSON(http.StatusInternalServerError, expense.Err{Message: "Prepare Fail"})
	}
	defer stmt.Close()
	rows := stmt.QueryRow(e.Title, e.Amount, e.Note, pq.Array(e.Tags), id)
	err = rows.Scan(&e.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, expense.Err{Message: "can't scan id"})
	}
	return c.JSON(http.StatusOK, e)
}

func (h *handler) CreateExpensesHandler(c echo.Context) error {
	isAuth, returnError := CheckAuth(c)
	if isAuth {
		return returnError
	}
	var e expense.Expenses
	err := c.Bind(&e)
	if err != nil {
		return c.JSON(http.StatusBadRequest, expense.Err{Message: err.Error()})
	}
	stmt, isError := h.Prepare(Create, c)
	if isError {
		return c.JSON(http.StatusInternalServerError, expense.Err{Message: "Prepare Fail"})
	}
	defer stmt.Close()
	rows := stmt.QueryRow(e.Title, e.Amount, e.Note, pq.Array(e.Tags))
	err = rows.Scan(&e.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, expense.Err{Message: "can't scan id"})
	}
	return c.JSON(http.StatusCreated, e)
}

func CheckAuth(c echo.Context) (bool, error) {
	if c.Request().Header.Get("Authorization") != "November 10, 2009" {
		return true, c.JSON(http.StatusUnauthorized, expense.Err{Message: "Unauthorized"})
	}
	return false, nil
}
