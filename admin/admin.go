package admin

import (
	"fmt"

	"github.com/Focinfi/sqs/config"
	"github.com/Focinfi/sqs/util/githubutil"
	"github.com/Focinfi/sqs/util/httputil"
	"github.com/gin-gonic/gin"
	"gopkg.in/gomail.v2"
)

func validateGithubLogin(ctx *gin.Context) {
	params := struct {
		Login string `json:"login"`
	}{}

	if err := ctx.BindJSON(&params); err != nil {
		return
	}

	result := githubutil.DefaultValidator.ContainsLogin(params.Login)
	httputil.ResponseOKData(ctx, gin.H{"isStargazer": result})
}

func sendSecretKeyToEmail(email string, accessKey string, secretKey string) error {
	emailConfig := config.Config.Email

	body := fmt.Sprintf(`
	This keys are for sqs service testing.

	AcessKey: %s
	SecretKey: %s
	`, accessKey, secretKey)
	m := gomail.NewMessage()
	m.SetHeader("From", emailConfig.From)
	m.SetHeader("To", email)
	m.SetHeader("Subject", "SQS Tesing Keys")
	m.SetBody("text/html", body)

	d := gomail.NewDialer(emailConfig.SMTP, 587, emailConfig.User, emailConfig.Password)
	return d.DialAndSend(m)
}
