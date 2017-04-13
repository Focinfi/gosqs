package models

const (
	AuthCodeKey = "authCode"
)

type NodeInfo struct {
	Addr     string `json:"addr"`
	CPU      int    `json:"cpu"`
	Memory   int    `json:"memory"`
	Resource int    `json:"resource"`
}

type UserAuth struct {
	//AccessKey for sqs basic key
	AccessKey string `json:"access_key"`
	// Secret for user auth
	SecretKey string `json:"secret_key"`
}
