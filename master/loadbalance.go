package master

import "github.com/Focinfi/sqs/models"

type InfoSlice []models.NodeInfo

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
