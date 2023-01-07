package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"
	"testing"

	expense "github.com/SThanutworapat/assessment/models"
	"github.com/stretchr/testify/assert"
)

func TestGetAllExpenses(t *testing.T) {
	done := make(chan bool)
	go func() {
		main()
		done <- true
	}()

	seedUser(t)
	var us []expense.Expenses

	res := request(http.MethodGet, uri("expenses"), nil)
	err := res.Decode(&us)

	assert.Nil(t, err)
	assert.EqualValues(t, http.StatusOK, res.StatusCode)
	assert.Greater(t, len(us), 0)

	request(http.MethodPost, uri("quit"), nil)
}

func TestGetExpensesById(t *testing.T) {
	done := make(chan bool)
	go func() {
		main()
		done <- true
	}()

	c := seedUser(t)

	var latest expense.Expenses
	res := request(http.MethodGet, uri("expenses", strconv.Itoa(c.ID)), nil)
	err := res.Decode(&latest)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, c.ID, latest.ID)
	assert.NotEmpty(t, latest.ID)
	assert.NotEmpty(t, latest.Title)
	assert.NotEmpty(t, latest.Tags)
	assert.NotEmpty(t, latest.Amount)
	assert.NotEmpty(t, latest.Note)

	request(http.MethodPost, uri("quit"), nil)
}

func TestCreateExpenses(t *testing.T) {
	done := make(chan bool)
	go func() {
		main()
		done <- true
	}()

	body := bytes.NewBufferString(`{
		"title": "strawberry smoothie",
		"amount": 79,
		"note": "night market promotion discount 10 bath",
		"tags": ["food", "beverage"]
	}`)
	var e expense.Expenses

	res := request(http.MethodPost, uri("expenses"), body)
	err := res.Decode(&e)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusCreated, res.StatusCode)
	assert.NotEqual(t, 0, e.ID)
	assert.Equal(t, "strawberry smoothie", e.Title)
	assert.Equal(t, 79.0, e.Amount)
	assert.Equal(t, "night market promotion discount 10 bath", e.Note)
	assert.Equal(t, []string{"food", "beverage"}, e.Tags)

	request(http.MethodPost, uri("quit"), nil)
}

func TestUpdateExpenses(t *testing.T) {
	done := make(chan bool)
	go func() {
		main()
		done <- true
	}()

	c := seedUser(t)
	body := bytes.NewBufferString(`{
		"title": "apple smoothie",
		"amount": 89,
		"note": "no discount",
		"tags": ["beverage"]
	}`)
	var e expense.Expenses

	res := request(http.MethodPut, uri("expenses", strconv.Itoa(c.ID)), body)
	err := res.Decode(&e)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.NotEqual(t, 0, e.ID)
	assert.Equal(t, c.ID, e.ID)
	assert.Equal(t, "apple smoothie", e.Title)
	assert.Equal(t, 89.0, e.Amount)
	assert.Equal(t, "no discount", e.Note)
	assert.Equal(t, []string{"beverage"}, e.Tags)

	request(http.MethodPost, uri("quit"), nil)
}

func seedUser(t *testing.T) expense.Expenses {
	var c expense.Expenses
	body := bytes.NewBufferString(`{
		"title": "strawberry smoothie",
		"amount": 79,
		"note": "night market promotion discount 10 bath",
		"tags": ["food", "beverage"]
	}`)
	err := request(http.MethodPost, uri("expenses"), body).Decode(&c)
	if err != nil {
		t.Fatal("can't create uomer:", err)
	}
	return c
}

func uri(paths ...string) string {
	host := "http://localhost:2565"
	if paths == nil {
		return host
	}
	url := append([]string{host}, paths...)
	return strings.Join(url, "/")
}

type Response struct {
	*http.Response
	err error
}

func (r *Response) Decode(v interface{}) error {
	if r.err != nil {
		return r.err
	}
	return json.NewDecoder(r.Body).Decode(v)
}

func request(method, url string, body io.Reader) *Response {
	req, _ := http.NewRequest(method, url, body)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "November 10, 2009")
	client := http.Client{}
	res, err := client.Do(req)
	return &Response{res, err}
}
