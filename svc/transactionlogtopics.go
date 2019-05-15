package svc

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// GetTopics - Get Transaction Log Topic by log
func GetTopics(lg *types.Log) []common.Hash {
	topics := lg.Topics
	return topics
}
