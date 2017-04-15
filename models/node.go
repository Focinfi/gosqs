package models

const (
	// AuthCodeKey for key in kv storage
	AuthCodeKey = "authCode"
)

// NodeInfo contains the basic stats information of one node
type NodeInfo struct {
	Addr     string `json:"addr"`
	CPU      int    `json:"cpu"`
	Memory   int    `json:"memory"`
	Resource int    `json:"resource"`
}

// UserAuth for one user auth info
type UserAuth struct {
	//AccessKey for sqs basic key
	AccessKey string `json:"access_key"`
	// Secret for user auth
	SecretKey string `json:"secret_key"`
}

// NodeRequestParams contains the basic info of every request
type NodeRequestParams struct {
	Token     string `json:"token"`
	QueueName string `json:"queue_name"`
	SquadName string `json:"squad_name,omitempty"`
}

// InfoSlice contains a slice of nodes
type InfoSlice []NodeInfo

func (ss InfoSlice) Len() int {
	return len(ss)
}

func (ss InfoSlice) Less(i, j int) bool {
	pi := ss[i].CPU + ss[i].Memory + ss[i].Resource
	pj := ss[j].CPU + ss[j].Memory + ss[j].Resource
	return pi < pj
}

func (ss InfoSlice) Swap(i, j int) {
	ss[i], ss[j] = ss[j], ss[i]
}
