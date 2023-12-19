package ethblocks

import (
	"github.com/spf13/viper"
)

// GetEthblocksClientAddr - used to get address from infura
func GetEthblocksClientAddr() string {
	v := viper.New()
	v.AutomaticEnv()
	clientAddr := v.GetString("ETHBLOCKS_CLIENT")
	return clientAddr
}

// GetEthblocksClient2Details - used to get address from sepolia.infura
func GetEthblocksClient2Details() (string, string, string) {
	v := viper.New()
	v.AutomaticEnv()
	clientAddr2 := v.GetString("ETHBLOCKS_CLIENT2")
	privateKey := v.GetString("ETHBLOCKS_PRIVATEKEY")
	toAddress := v.GetString("ETHBLOCKS_TOADDRESS")
	return clientAddr2, privateKey, toAddress
}
