package svc

import (
	"log"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// PrintTopic - Print Transaction Log Topic
func PrintTopic(s common.Hash) {
	log.Println("Topic  : ", s.Hex())
}

// GetTopics - Get Transaction Log Topic by log
func GetTopics(lg *types.Log) []common.Hash {
	topics := lg.Topics
	return topics
}
