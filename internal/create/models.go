package create

type UserInfo struct {
	Name       string `json:"name"`
	Surname    string `json:"surname"`
	Patronymic string `json:"patronymic,omitempty"`
	Gender     string `json:"gender,omitempty"`
	Age        int    `json:"age,omitempty"`
	Nation     string `json:"nation,omitempty"`
}

// Gender enum i guess
var GenderList [2]string = [2]string{"male", "female"}

type AgeResp struct {
	Age int `json:"age,omitempty"`
}
type Country struct {
	CountryID   string  `json:"country_id"`
	Probability float64 `json:"probability"`
}
type NationResp struct {
	Count   int       `json:"count,omitempty"`
	Country []Country `json:"country,omitempty"`
}

type GenderResp struct {
	Gender string `json:"gender,omitempty"`
}
