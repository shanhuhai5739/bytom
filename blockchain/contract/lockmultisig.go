package contract

import (
	"encoding/json"
	"fmt"
	"github.com/bytom/blockchain/txbuilder"
)

// LockMultiSig stores the information of LockWithMultiSig contract
type LockMultiSig struct {
	CommonInfo
	PubKeys []PubKeyInfo
}

// DecodeLockMultiSig unmarshal JSON-encoded data of contract action
func DecodeLockMultiSig(data []byte) (ContractAction, error) {
	a := new(LockPubKey)
	err := json.Unmarshal(data, a)
	return a, err
}

// BuildContractReq create new ContractReq which contain contract's name and arguments
func (a *LockMultiSig) BuildContractReq(contractName string) (*ContractReq, error) {
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
func (a *LockMultiSig) Build() (*string, error) {
	var buildReqStr string

	if a.Alias {
		buildReqStr = fmt.Sprintf(buildAcctRecvReqFmtByAlias, a.OutputID, a.AssetInfo, a.Amount, a.AccountInfo, a.BtmGas, a.AccountInfo)
	} else {
		buildReqStr = fmt.Sprintf(buildAcctRecvReqFmt, a.OutputID, a.AssetInfo, a.Amount, a.AccountInfo, a.BtmGas, a.AccountInfo)
	}

	return &buildReqStr, nil
}

// AddArgs add the parameters for contract
func (a *LockMultiSig) AddArgs(tpl *txbuilder.Template) (*txbuilder.Template, error) {
	var err error

	if tpl, err = addPubKeyArgs(tpl, a.PubKeys); err != nil {
		return nil, err
	}

	return tpl, nil
}
