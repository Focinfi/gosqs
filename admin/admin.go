package admin

import (
	"fmt"
	"time"

	"github.com/Focinfi/sqs/config"
	"github.com/Focinfi/sqs/errors"
	"github.com/Focinfi/sqs/util/githubutil"
	"github.com/Focinfi/sqs/util/httputil"
	"github.com/Focinfi/sqs/util/token"
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
	params := githubLoginParam{}
	if err := ctx.BindJSON(&params); err != nil {
		return
	}

	result := githubutil.DefaultValidator.ContainsLogin(params.Login)
	httputil.ResponseOKData(ctx, gin.H{"isStargazer": result})
}

// SendGithubEmailSecretKey send the secret key to the email of the given github login
func SendGithubEmailSecretKey(ctx *gin.Context) {
	params := githubLoginParam{}
	if err := ctx.BindJSON(&params); err != nil {
		return
	}

	validate := githubutil.DefaultValidator.ContainsLogin(params.Login)
	if !validate {
		httputil.ResponseErr(ctx, errors.NotSQSStargazer)
		return
	}

	email, err := githubutil.EmailForUserLogin(params.Login)
	paramsKey := config.Config.UserGithubLoginKey
	secretKey, err := token.Default.Make(config.Config.BaseSecret, map[string]interface{}{paramsKey: params.Login}, time.Hour)
	if err != nil {
		httputil.ResponseErr(ctx, errors.NewInternalWrap(err))
		return
	}

	err = sendSecretKeyToEmail(email, params.Login, secretKey)
	if err != nil {
		httputil.ResponseErr(ctx, errors.NewInternalWrap(err))
		return
	}

	httputil.ResponseOKData(ctx, gin.H{"email": email})
}

func sendSecretKeyToEmail(email string, accessKey string, secretKey string) error {
	emailConfig := config.Config.Email

	body := fmt.Sprintf(`
	<h3>These keys are for sqs service testing.</h3>
	<p>AcessKey: %s </p>
	<p>SecretKey: %s </p>
	`, accessKey, secretKey)
	m := gomail.NewMessage()
	m.SetHeader("From", emailConfig.From)
	m.SetHeader("To", email)
	m.SetHeader("Subject", "SQS Tesing Keys")
	m.SetBody("text/html", body)

	d := gomail.NewDialer(emailConfig.SMTP, emailConfig.Port, emailConfig.User, emailConfig.Password)
	return d.DialAndSend(m)
}
