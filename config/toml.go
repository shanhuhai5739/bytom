package config

import (
	"path"

	cmn "github.com/tendermint/tmlibs/common"
)

/****** these are for production settings ***********/
func EnsureRoot(rootDir string, network string) {
	cmn.EnsureDir(rootDir, 0700)
	cmn.EnsureDir(rootDir+"/data", 0700)

	configFilePath := path.Join(rootDir, "config.toml")

	// Write default config file if missing.
	if !cmn.FileExists(configFilePath) {
		cmn.MustWriteFile(configFilePath, []byte(selectNetwork(network)), 0644)
	}
}

var defaultConfigTmpl = `# This is a TOML config file.
# For more information, see https://github.com/toml-lang/toml
fast_sync = true
db_backend = "leveldb"
api_addr = "0.0.0.0:9888"
`

var mainNetConfigTmpl = `chain_id = "mainnet"
[p2p]
laddr = "tcp://0.0.0.0:46657"
seeds = "52.83.204.67:46657,52.83.161.31:46657,52.83.113.196:46657,52.83.127.78:46657,52.83.177.255:46657,52.83.168.61:46657,52.83.232.126:46657,52.83.196.210:46657,52.83.106.68:46657,52.83.139.63:46657,52.83.234.161:46657,52.83.189.31:46657,52.83.218.101:46657,47.100.214.154:46657,47.100.109.199:46657,47.100.105.165:46657,47.100.247.186:46657,101.132.176.116:46657,47.100.246.237:46657,47.100.247.164:46657,198.74.61.131:46657,45.79.213.28:46657,212.111.41.245:46657,139.198.177.243:46657,139.198.177.164:46657,139.198.177.190:46657,139.198.177.231:46657,101.37.164.153:46657"
`

var testNetConfigTmpl = `chain_id = "wisdom"
[p2p]
laddr = "tcp://0.0.0.0:46656"
seeds = "52.83.107.224:46656,52.83.107.224:46656,52.83.251.197:46656"
`

var soloNetConfigTmpl = `chain_id = "solonet"
[p2p]
laddr = "tcp://0.0.0.0:46658"
seeds = ""
`

// Select network seeds to merge a new string.
func selectNetwork(network string) string {
	if network == "testnet" {
		return defaultConfigTmpl + testNetConfigTmpl
	} else if network == "mainnet" {
		return defaultConfigTmpl + mainNetConfigTmpl
	} else {
		return defaultConfigTmpl + soloNetConfigTmpl
	}
}
