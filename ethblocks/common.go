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
