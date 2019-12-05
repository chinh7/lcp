package token

import (
	"log"
	"strconv"

	"github.com/QuoineFinancial/vertex/abi"
	"github.com/QuoineFinancial/vertex/crypto"
	"github.com/QuoineFinancial/vertex/engine"
	"github.com/QuoineFinancial/vertex/gas"
	"github.com/QuoineFinancial/vertex/storage"
)

// Token contract
type Token struct {
	contract *storage.Account
}

func (token *Token) invokeContract(caller crypto.Address, method string, values []string) (uint64, error) {
	contract, err := token.contract.GetContract()
	if err != nil {
		return 0, err
	}

	function, err := contract.Header.GetFunction(method)
	if err != nil {
		return 0, err
	}
	methodArgs, err := abi.EncodeFromString(function.Parameters, values)
	if err != nil {
		return 0, err
	}

	engine := engine.NewEngine(token.contract, caller, &gas.FreePolicy{}, -1)
	ret, _, err := engine.Ignite(method, methodArgs)
	if err != nil {
		return 0, err
	}
	return *ret, err
}

// GetBalance retrieve token balance by address
func (token *Token) GetBalance(addr crypto.Address) (uint64, error) {
	ret, err := token.invokeContract(addr, "get_balance", []string{addr.String()})
	log.Printf("Get balance of %v: %v\n", addr.String(), ret)
	return ret, err
}

// Transfer transfer token from caller address to another address
func (token *Token) Transfer(caller crypto.Address, addr crypto.Address, amount uint64) error {
	ret, err := token.invokeContract(caller, "transfer", []string{addr.String(), strconv.FormatUint(amount, 10)})
	if int(ret) < 0 {
		panic("Burn gas error")
	}
	log.Printf("Transfer %v %v to %v\n", amount, caller.String(), addr.String())
	return err
}

// GetContract account
func (token *Token) GetContract() *storage.Account {
	return token.contract
}

// NewToken from contract
func NewToken(contract *storage.Account) *Token {
	return &Token{
		contract: contract,
	}
}
