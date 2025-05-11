package create_usecase

import (
	"encoding/json"
	"net/http"
	"sort"
	"test/internal/create"
	"test/pkg/logger"
	"test/pkg/postgres"
	"test/pkg/simple_http"

	"github.com/gin-gonic/gin"
)

type Creator struct {
	DB        *postgres.Postgres
	Logger    logger.Logger
	GenderApi string
	AgeApi    string
	NationApi string
}

const CreateQuery = "insert into Users (Name, Surname, Patronymic, Nation, Gender, Age) values ($1, $2, $3, $4, $5, $6) returning ID"

func NewCreator(DB *postgres.Postgres, Logger logger.Logger, GenderApi string, AgeApi string, NationApi string) create.Creator {
	Logger.Info("Create, Update, Delete service ready")
	return &Creator{DB: DB, Logger: Logger, AgeApi: AgeApi, GenderApi: GenderApi, NationApi: NationApi}
}
func (Creator *Creator) Create(c *gin.Context) {
	var ID int
	var GeneralInfo create.UserInfo
	var AgeInfo create.AgeResp
	var GenderInfo create.GenderResp
	var NationInfo create.NationResp
	db := Creator.DB.Get()

	Creator.Logger.Info("Got Add request")

	if err := c.ShouldBindBodyWithJSON(&GeneralInfo); err != nil {
		Creator.Logger.Error(err)
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	if GeneralInfo.Name == "" || GeneralInfo.Surname == "" {
		Creator.Logger.Error("Name or surname absent")
		c.JSON(http.StatusBadRequest, "Name or surname absent")
		return
	}
	data, err := simple_http.MakeRequest(Creator.NationApi, GeneralInfo.Name)
	if err != nil {
		Creator.Logger.Error(err)
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	err = json.Unmarshal(data, &NationInfo)
	if err != nil {
		Creator.Logger.Error(err)
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	if NationInfo.Country == nil || len(NationInfo.Country) == 0 || NationInfo.Count == 0 {
		c.JSON(http.StatusBadRequest, "Nation for "+GeneralInfo.Name+" not found")
		Creator.Logger.Error("Nation for " + GeneralInfo.Name + " not found")
		return
	}
	data, err = simple_http.MakeRequest(Creator.AgeApi, GeneralInfo.Name)
	if err != nil {
		Creator.Logger.Error(err)
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	err = json.Unmarshal(data, &AgeInfo)
	if err != nil {
		Creator.Logger.Error(err)
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	if AgeInfo.Age == 0 {
		Creator.Logger.Error("Age for " + GeneralInfo.Name + " not found")
		c.JSON(http.StatusBadRequest, "Age for "+GeneralInfo.Name+" not found")
		return
	}
	data, err = simple_http.MakeRequest(Creator.GenderApi, GeneralInfo.Name)
	if err != nil {
		Creator.Logger.Error(err)
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	err = json.Unmarshal(data, &GenderInfo)
	if err != nil {
		Creator.Logger.Error(err)
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	if GenderInfo.Gender == "" || !checkGender(GenderInfo) {
		Creator.Logger.Error("Gender for " + GeneralInfo.Name + " not found")
		c.JSON(http.StatusBadRequest, "Gender for "+GeneralInfo.Name+"not found")
		return
	}
	GeneralInfo.Age = AgeInfo.Age
	GeneralInfo.Gender = GenderInfo.Gender
	GeneralInfo.Nation = sortNation(NationInfo)
	Creator.Logger.Debug("I sorted list, most probable nation is " + sortNation(NationInfo))
	db.QueryRow(CreateQuery, GeneralInfo.Name, GeneralInfo.Surname, GeneralInfo.Patronymic, GeneralInfo.Nation, GeneralInfo.Gender, GeneralInfo.Age).Scan(&ID)
	Creator.Logger.Info("User added")
	c.JSON(http.StatusOK, ID)
}

func sortNation(Nation create.NationResp) string {
	sort.Slice(Nation.Country, func(i, j int) bool {
		return Nation.Country[i].Probability > Nation.Country[j].Probability
	})
	return Nation.Country[0].CountryID
}

func checkGender(Gender create.GenderResp) bool {
	for _, v := range create.GenderList {
		if v == Gender.Gender {
			return true
		}
	}
	return false
}
