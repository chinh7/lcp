package token

import (
	"strconv"

	"github.com/QuoineFinancial/vertex/abi"
	"github.com/QuoineFinancial/vertex/crypto"
	"github.com/QuoineFinancial/vertex/engine"
	"github.com/QuoineFinancial/vertex/gas"
	"github.com/QuoineFinancial/vertex/storage"
	"github.com/tendermint/tendermint/abci/types"
)

// Token contract
type Token struct {
	state    *storage.State
	contract *storage.Account
}

func (token *Token) invokeContract(caller crypto.Address, method string, values []string) (uint64, []types.Event, error) {
	contract, err := token.contract.GetContract()
	if err != nil {
		return 0, nil, err
	}

	function, err := contract.Header.GetFunction(method)
	if err != nil {
		return 0, nil, err
	}
	methodArgs, err := abi.EncodeFromString(function.Parameters, values)
	if err != nil {
		return 0, nil, err
	}

	engine := engine.NewEngine(token.state, token.contract, caller, &gas.FreePolicy{}, 0)
	ret, err := engine.Ignite(method, methodArgs)
	if err != nil {
		return 0, nil, err
	}
	return ret, engine.GetEvents(), err
}

// GetBalance retrieve token balance by address
func (token *Token) GetBalance(addr crypto.Address) (uint64, error) {
	ret, _, err := token.invokeContract(addr, "get_balance", []string{addr.String()})
	return ret, err
}

// Transfer transfer token from caller address to another address
func (token *Token) Transfer(caller crypto.Address, addr crypto.Address, amount uint64) ([]types.Event, error) {
	ret, events, err := token.invokeContract(caller, "transfer", []string{addr.String(), strconv.FormatUint(amount, 10)})
	if int(ret) < 0 {
		panic("transfer token failed")
	}
	return events, err
}

// GetContract account
func (token *Token) GetContract() *storage.Account {
	return token.contract
}

// NewToken from contract
func NewToken(state *storage.State, contract *storage.Account) *Token {
	return &Token{
		state:    state,
		contract: contract,
	}
}
