package read_usecase

import (
	"net/http"
	"strconv"
	"test/internal/read"
	"test/pkg/logger"
	"test/pkg/postgres"

	"github.com/gin-gonic/gin"
)

type Search struct {
	DB     *postgres.Postgres
	Logger logger.Logger
}

const SearchEnd = " order by ID "
const SearchBegin = "select Name, Surname, Patronymic, Gender, Age, Nation from Users "

func (S *Search) Read(c *gin.Context) {
	var Resp read.SearchResponse
	var QueryStr = SearchBegin
	var Values []interface{}
	var iter = 0
	var isRows = 0

	S.Logger.Info("Got Search Request")
	if c.Query("page_num") == "" || c.Query("page_size") == "" {
		c.JSON(http.StatusBadRequest, "No page_num or page_size")
		return
	}
	Offset, err := strconv.Atoi(c.Query("page_num"))
	if err != nil || Offset < 0 {
		c.JSON(http.StatusBadRequest, "Bad page num")
		return
	}
	PageSize, err := strconv.Atoi(c.Query("page_size"))
	if err != nil || PageSize < 0 {
		S.Logger.Error(err.Error())
		c.JSON(http.StatusBadRequest, "Bad page size")
		return
	}

	if Age, ok := c.GetQuery("age"); ok {
		if iter == 0 {
			QueryStr += " WHERE "
		}
		iter++
		QueryStr += "Age = $" + strconv.Itoa(iter)
		AgeNum, err := strconv.Atoi(Age)
		if err != nil || AgeNum < 0 {
			c.JSON(http.StatusBadRequest, "Bad Age")
		}
		Values = append(Values, AgeNum)

	}
	if Name, ok := c.GetQuery("name"); ok {
		QueryStr, Values, iter = addStringQuery(QueryStr, Values, "Name", Name, iter)
	}
	if Surname, ok := c.GetQuery("name"); ok {
		QueryStr, Values, iter = addStringQuery(QueryStr, Values, "Surname", Surname, iter)
	}
	if Gender, ok := c.GetQuery("gender"); ok {
		QueryStr, Values, iter = addStringQuery(QueryStr, Values, "Gender", Gender, iter)
	}
	if Patronymic, ok := c.GetQuery("patronymic"); ok {
		QueryStr, Values, iter = addStringQuery(QueryStr, Values, "Patronymic", Patronymic, iter)
	}
	if Nation, ok := c.GetQuery("nation"); ok {
		QueryStr, Values, iter = addStringQuery(QueryStr, Values, "Nation", Nation, iter)
	}
	QueryStr += SearchEnd
	QueryStr += " LIMIT $" + strconv.Itoa(iter+1) + " OFFSET $" + strconv.Itoa(iter+2)
	Values = append(Values, PageSize)
	Values = append(Values, PageSize*Offset)
	db := S.DB.Get()
	S.Logger.Debug("Making request with " + QueryStr)
	rows, err := db.Query(QueryStr, Values...)
	rows.Scan(&isRows)
	if err != nil {
		S.Logger.Error(err)
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	for rows.Next() {
		isRows++
		val := read.UserInfo{}
		err = rows.Scan(&val.Name, &val.Surname, &val.Patronymic, &val.Gender, &val.Age, &val.Nation)
		if err != nil {
			S.Logger.Error(err)
			c.JSON(http.StatusInternalServerError, "")
			return
		}
		Resp.UserInfo = append(Resp.UserInfo, val)
	}
	if isRows == 0 {
		S.Logger.Debug("no rows found")
		c.JSON(http.StatusNoContent, "no rows found")
		return
	}
	S.Logger.Info("Search successfull")
	c.JSON(http.StatusOK, Resp)
}

func addStringQuery(QueryStr string, Values []interface{}, Param string, Value string, iter int) (string, []interface{}, int) {
	if iter == 0 {
		QueryStr += " where "
	}
	if iter < 0 {
		QueryStr += " and "
	}
	iter++
	QueryStr += Param + " = $" + strconv.Itoa(iter)
	Values = append(Values, Value)
	return QueryStr, Values, iter
}

func NewSearchUC(DB *postgres.Postgres, Logger logger.Logger) read.Reader {
	Logger.Info("Search service ready")
	return &Search{DB: DB, Logger: Logger}
}
