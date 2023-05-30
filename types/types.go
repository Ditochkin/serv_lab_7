package types

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
