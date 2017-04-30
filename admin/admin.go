package admin

import (
	"fmt"

	"github.com/Focinfi/gosqs/config"
	"github.com/Focinfi/gosqs/errors"
	"github.com/Focinfi/gosqs/log"
	"github.com/Focinfi/gosqs/util/githubutil"
	"github.com/Focinfi/gosqs/util/httputil"
	"github.com/Focinfi/gosqs/util/token"
	"github.com/gin-gonic/gin"
	"github.com/go-gomail/gomail"
)

type githubLoginParam struct {
	Login string `json:"login"`
}

// Admin a server
type Admin struct {
	*gin.Engine
}

// NewAdmin returns a new Admin
func NewAdmin() *Admin {
	engine := gin.Default()
	engine.GET("/validateGithubLogin", ValidateGithubLogin)
	engine.GET("/sendGithubEmailSecretKey", SendGithubEmailSecretKey)
	return &Admin{Engine: engine}
}

// ValidateGithubLogin validates the github login
func ValidateGithubLogin(ctx *gin.Context) {
	login := ctx.Param("login")

	result := githubutil.DefaultValidator.ContainsLogin(login)
	httputil.ResponseOKData(ctx, gin.H{"isStargazer": result})
}

// SendGithubEmailSecretKey send the secret key to the email of the given github login
func SendGithubEmailSecretKey(ctx *gin.Context) {
	login := ctx.Param("login")
	validate := githubutil.DefaultValidator.ContainsLogin(login)
	if !validate {
		httputil.ResponseErr(ctx, errors.NotSQSStargazer)
		return
	}

	email, err := githubutil.EmailForUserLogin(login)
	paramsKey := config.Config.UserGithubLoginKey
	secretKey, err := token.Default.Make(config.Config.BaseSecret, map[string]interface{}{paramsKey: login}, -1)
	if err != nil {
		log.Internal.Error(err)
		httputil.ResponseErr(ctx, errors.NewInternalWrap(err))
		return
	}

	err = sendSecretKeyToEmail(email, login, secretKey)
	if err != nil {
		log.Internal.Error(err)
		httputil.ResponseErr(ctx, errors.NewInternalWrap(err))
		return
	}

	httputil.ResponseOKData(ctx, gin.H{"email": email})
}

func sendSecretKeyToEmail(email string, accessKey string, secretKey string) error {
	emailConfig := config.Config.Email

	body := fmt.Sprintf(`
	<h3>These keys are for gosqs service testing.</h3>
	<p>AcessKey: %s </p>
	<p>SecretKey: %s </p>
	`, accessKey, secretKey)
	m := gomail.NewMessage()
	m.SetHeader("From", emailConfig.From)
	m.SetHeader("To", email)
	m.SetHeader("Subject", "gosqs Tesing Keys")
	m.SetBody("text/html", body)

	d := gomail.NewDialer(emailConfig.SMTP, emailConfig.Port, emailConfig.User, emailConfig.Password)
	return d.DialAndSend(m)
}
