package client

type InfoJsonAge struct {
	Count int64 `json:"count"`
	Age   int64 `json:"age"`
}

type InfoJsonGender struct {
	Count  int64  `json:"count"`
	Gender string `json:"gender"`
}

type InfoJsonNationality struct {
	CountryProbs []Country `json:"country"`
}

type Country struct {
	CountryId string `json:"country_id"`
}
