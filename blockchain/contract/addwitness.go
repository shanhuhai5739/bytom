package contract

import (
	"encoding/hex"
	"fmt"

	"github.com/bytom/blockchain/txbuilder"
	"github.com/bytom/crypto/ed25519/chainkd"
	chainjson "github.com/bytom/encoding/json"
	"github.com/bytom/errors"
)

// CommonInfo is the common elements for contract transaction
type CommonInfo struct {
	OutputID    string `json:"outputid"`
	AccountInfo string `json:"account"`
	AssetInfo   string `json:"asset"`
	Amount      string `json:"amount"`
	Alias       bool   `json:"alias"`
	BtmGas      string `json:"btm_gas"`
}

// PubKeyInfo is the elements of generating public key
type PubKeyInfo struct {
	RootPubKey string   `json:"rootpub"`
	Path       []string `json:"path"`
}

// PaymentInfo is the requirement of clause for contract
type PaymentInfo struct {
	InnerAccountInfo string `json:"inner_account"`
	InnerAssetInfo   string `json:"inner_asset"`
	InnerAmount      string `json:"inner_amount"`
	InnerProgram     string `json:"inner_prog"`
}

// ParamInfo is the entire struct of contract arguments
type ParamInfo struct {
	frontData   []string
	pubKeyInfos []PubKeyInfo
	lastData    []string
}

// CommonPubInfo is the elements of RawTxSigWitness
type CommonPubInfo struct {
	rootPubKeys []chainkd.XPub
	paths       [][]chainjson.HexBytes
	quorum      int
}

// NewPubKeyInfo creates a new PubKeyInfo
func NewPubKeyInfo(rootPub string, path []string) PubKeyInfo {
	return PubKeyInfo{
		RootPubKey: rootPub,
		Path:       path,
	}
}

// NewParamInfo creates a new ParamInfo
func NewParamInfo(front []string, pubKeys []PubKeyInfo, last []string) ParamInfo {
	return ParamInfo{
		frontData:   front,
		pubKeyInfos: pubKeys,
		lastData:    last,
	}
}

func reconstructTpl(tpl *txbuilder.Template, si *txbuilder.SigningInstruction) *txbuilder.Template {
	length := len(tpl.SigningInstructions)
	if length <= 0 {
		length = 1
		tpl.SigningInstructions = append(tpl.SigningInstructions, si)
		tpl.SigningInstructions[length-1].Position = 0
	} else {
		tpl.SigningInstructions[0] = si
	}

	return tpl
}

func convertPubInfo(pubKeyInfos []PubKeyInfo) (*CommonPubInfo, error) {
	rootPubKey := chainkd.XPub{}
	path := []chainjson.HexBytes{}
	commonPubInfo := CommonPubInfo{}

	for _, pubInfo := range pubKeyInfos {
		hexPubKey, err := hex.DecodeString(pubInfo.RootPubKey)
		if err != nil {
			return nil, err
		}
		copy(rootPubKey[:], hexPubKey[:])

		if len(pubInfo.Path) != 2 {
			buf := fmt.Sprintf("the length of path [%d] is not equal 2!", len(pubInfo.Path))
			err := errors.New(buf)
			return nil, err
		}

		for _, strPath := range pubInfo.Path {
			hexPath, err := hex.DecodeString(strPath)
			if err != nil {
				return nil, err
			}

			path = append(path, hexPath)
		}

		commonPubInfo.rootPubKeys = append(commonPubInfo.rootPubKeys, rootPubKey)
		commonPubInfo.paths = append(commonPubInfo.paths, path)
		commonPubInfo.quorum++
	}

	return &commonPubInfo, nil
}

func addPubKeyArgs(tpl *txbuilder.Template, pubKeyInfos []PubKeyInfo) (*txbuilder.Template, error) {
	si := txbuilder.SigningInstruction{}

	pubInfo, err := convertPubInfo(pubKeyInfos)
	if err != nil {
		return nil, err
	}

	err = si.AddRawTxSigWitness(pubInfo.rootPubKeys, pubInfo.paths, pubInfo.quorum)
	if err != nil {
		return nil, err
	}

	tpl = reconstructTpl(tpl, &si)
	return tpl, nil
}

func addDataArgs(tpl *txbuilder.Template, value []string) (*txbuilder.Template, error) {
	var dataWitness []chainjson.HexBytes
	for _, v := range value {
		data, err := hex.DecodeString(v)
		if err != nil {
			return nil, err
		}
		dataWitness = append(dataWitness, data)
	}

	si := txbuilder.SigningInstruction{}
	si.AddDataWitness(dataWitness)

	tpl = reconstructTpl(tpl, &si)
	return tpl, nil
}

func addParamArgs(tpl *txbuilder.Template, pubKeyValueInfo ParamInfo) (*txbuilder.Template, error) {
	si := txbuilder.SigningInstruction{}

	if pubKeyValueInfo.frontData != nil {
		var frontDataWitness []chainjson.HexBytes
		for _, data := range pubKeyValueInfo.frontData {
			front, err := hex.DecodeString(data)
			if err != nil {
				return nil, err
			}
			frontDataWitness = append(frontDataWitness, front)
		}

		si.AddDataWitness(frontDataWitness)
	}

	if pubKeyValueInfo.pubKeyInfos != nil {
		pubInfo, err := convertPubInfo(pubKeyValueInfo.pubKeyInfos)
		if err != nil {
			return nil, err
		}

		err = si.AddRawTxSigWitness(pubInfo.rootPubKeys, pubInfo.paths, pubInfo.quorum)
		if err != nil {
			return nil, err
		}
	}

	if pubKeyValueInfo.lastData != nil {
		var lastDataWitness []chainjson.HexBytes
		for _, data := range pubKeyValueInfo.lastData {
			front, err := hex.DecodeString(data)
			if err != nil {
				return nil, err
			}
			lastDataWitness = append(lastDataWitness, front)
		}

		si.AddDataWitness(lastDataWitness)
	}

	tpl = reconstructTpl(tpl, &si)
	return tpl, nil
}
