package contract

import (
	"encoding/json"
	"fmt"
	"github.com/bytom/blockchain/txbuilder"
)

// LockPubKey stores the information of LockWithPublicKey contract
type LockPubKey struct {
	CommonInfo
	PubKeyInfo
}

// DecodeLockPubKey unmarshal JSON-encoded data of contract action
func DecodeLockPubKey(data []byte) (ContractAction, error) {
	a := new(LockPubKey)
	err := json.Unmarshal(data, a)
	return a, err
}

// BuildContractReq create new ContractReq which contain contract's name and arguments
func (a *LockPubKey) BuildContractReq(contractName string) (*ContractReq, error) {
	arguments, err := json.Marshal(a)
	if err != nil {
		return nil, err
	}

	return &ContractReq{
		ContractName: contractName,
		ContractArgs: arguments,
	}, nil
}

// Build create a transaction request
func (a *LockPubKey) Build() (*string, error) {
	var buildReqStr string

	if a.Alias {
		buildReqStr = fmt.Sprintf(buildAcctRecvReqFmtByAlias, a.OutputID, a.AssetInfo, a.Amount, a.AccountInfo, a.BtmGas, a.AccountInfo)
	} else {
		buildReqStr = fmt.Sprintf(buildAcctRecvReqFmt, a.OutputID, a.AssetInfo, a.Amount, a.AccountInfo, a.BtmGas, a.AccountInfo)
	}

	return &buildReqStr, nil
}

// AddArgs add the parameters for contract
func (a *LockPubKey) AddArgs(tpl *txbuilder.Template) (*txbuilder.Template, error) {
	var err error
	pubInfo := NewPubKeyInfo(a.RootPubKey, a.Path)

	if tpl, err = addPubKeyArgs(tpl, []PubKeyInfo{pubInfo}); err != nil {
		return nil, err
	}

	return tpl, nil
}
