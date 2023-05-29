package types

type Country struct {
	Id          int
	CountryName string
}

type University struct {
	Id             int
	CountryName    string
	UniversityName string
}

type RankingCriteria struct {
	Id           int
	SystemID     int
	CriteriaName string
}

type ChangeStudentStaffRatio struct {
	UniversityName string
	Year           int
	NewStaffRatio  int
}

type AddUniversityRankingYear struct {
	UniversityName string
	CriteriaName   string
	Year           int
	Score          int
}

type Publisher struct {
	Id            int
	PublisherName string
}

type ChangePublisher struct {
	NewPublisherName string
	PublisherName    string
}

type GamePublisher struct {
	GameName      string
	PublisherName string
}

type PlatformYear struct {
	Year int
}

type User struct {
	Id       int64
	Name     string `json:"name"`
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
}
