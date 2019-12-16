package token

import (
	"strconv"

	"github.com/QuoineFinancial/liquid-chain/abi"
	"github.com/QuoineFinancial/liquid-chain/crypto"
	"github.com/QuoineFinancial/liquid-chain/engine"
	"github.com/QuoineFinancial/liquid-chain/event"
	"github.com/QuoineFinancial/liquid-chain/gas"
	"github.com/QuoineFinancial/liquid-chain/storage"
)

// Token contract
type Token struct {
	state    *storage.State
	contract *storage.Account
}

func (token *Token) invokeContract(caller crypto.Address, method string, values []string) (uint64, []event.Event, error) {
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
func (token *Token) Transfer(caller crypto.Address, addr crypto.Address, amount uint64) ([]event.Event, error) {
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
