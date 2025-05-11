package redactor_usecase

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	redactor "test/internal/update_delete"
	"test/pkg/logger"
	"test/pkg/postgres"

	"github.com/gin-gonic/gin"
)

type Redactor struct {
	Logger logger.Logger
	DB     postgres.Postgres
}

const QueryStrBegin = "update Users set "
const DelQuery = "delete from Users where ID = $1;"

func (R *Redactor) Update(c *gin.Context) {
	var Values map[string]interface{}
	var iter = 1
	var Query []interface{}
	QueryStr := QueryStrBegin
	R.Logger.Info("Got Change request")
	IDString := c.Query("ID")
	ID, err := strconv.Atoi(IDString)
	if err != nil {
		R.Logger.Error(err)
		c.JSON(http.StatusBadRequest, err)
		return
	}
	data, err := io.ReadAll(c.Request.Body)
	json.Unmarshal(data, &Values)
	if val, ok := Values["name"]; ok {
		Query = append(Query, val)
		QueryStr = QueryStr + "Name" + "=$" + strconv.Itoa(iter) + ","
		iter++
	}
	if val, ok := Values["surname"]; ok {
		Query = append(Query, val)
		QueryStr = QueryStr + "Surname" + "=$" + strconv.Itoa(iter) + ","
		iter++
	}
	if val, ok := Values["patronymic"]; ok {
		Query = append(Query, val)
		QueryStr = QueryStr + "Patronymic" + "=$" + strconv.Itoa(iter) + ","
		iter++
	}
	if val, ok := Values["nation"]; ok {
		Query = append(Query, val)
		QueryStr = QueryStr + "Nation" + "=$" + strconv.Itoa(iter) + ","
		iter++
	}
	if val, ok := Values["age"]; ok {
		Query = append(Query, val)
		QueryStr = QueryStr + "Age" + "=$" + strconv.Itoa(iter) + ","
		iter++
	}
	if val, ok := Values["gender"]; ok {
		Query = append(Query, val)
		QueryStr = QueryStr + "Gender" + "=$" + strconv.Itoa(iter) + " "
		iter++
	}
	QueryStr = QueryStr[:len(QueryStr)-1]
	QueryStr = QueryStr + " where ID=" + strconv.Itoa(ID) + " returning ID;"
	if len(Values) > iter-1 || iter == 1 {
		R.Logger.Debug("Bad params or no params")
		c.JSON(http.StatusBadRequest, "Bad params or no params")
		return
	}
	db := R.DB.Get()
	db.QueryRow(QueryStr, Query...).Scan(&ID)
	FinalID := strconv.Itoa(ID)
	R.Logger.Debug("Query str ready " + QueryStr)
	R.Logger.Info("Changed Entry")
	c.JSON(http.StatusOK, "Entry modifed ID:"+FinalID)
}
func (R *Redactor) Delete(c *gin.Context) {
	R.Logger.Info("Got Delete request")
	IDString := c.Query("ID")
	ID, err := strconv.Atoi(IDString)
	if err != nil {
		R.Logger.Error(err)
		c.JSON(http.StatusBadRequest, err)
		return
	}
	db := R.DB.Get()
	res, err := db.Exec(DelQuery, ID)
	if err != nil || res == nil {
		R.Logger.Error(err)
		c.JSON(http.StatusBadRequest, err)
		return
	}
	rows, err := res.RowsAffected()
	if err != nil {
		R.Logger.Error(err)
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	if rows == 0 {
		c.JSON(http.StatusNoContent, "No entry found")
		R.Logger.Error("no entry")
		return
	}
	R.Logger.Info("Entry Deleted")
	c.JSON(http.StatusOK, "Entry Deleted")
}
func NewRedactor(Logger logger.Logger, DB postgres.Postgres) redactor.Redactor {
	return &Redactor{Logger: Logger, DB: DB}
}
