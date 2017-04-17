package models

const (
	// NodesKey for the nodes info storage key
	NodesKey = "sqs.nodes"
)

// NodeInfo contains the basic stats information of one node
type NodeInfo struct {
	Addr   string `json:"addr"`
	CPU    int    `json:"cpu"`
	Memory int    `json:"memory"`
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

// NodeSlice contains a slice of nodes
type NodeSlice []NodeInfo

func (ss NodeSlice) Len() int {
	return len(ss)
}

func (ss NodeSlice) Less(i, j int) bool {
	pi := ss[i].CPU + ss[i].Memory
	pj := ss[j].CPU + ss[j].Memory
	return pi < pj
}

func (ss NodeSlice) Swap(i, j int) {
	ss[i], ss[j] = ss[j], ss[i]
}
