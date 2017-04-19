package admin

type Auther interface {
	Auth(accessKey string, secretKey string) (userID int64, err error)
}

var DefaultAuther Auther

func init() {
	DefaultAuther = NewGithubAuther()
}
