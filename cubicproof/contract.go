// Code generated
// This file is a generated precompile contract config with stubbed abstract functions.
// The file is generated by a template. Please inspect every code and comment in this file before use.

package cubicproof

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/ava-labs/subnet-evm/accounts/abi"
	"github.com/ava-labs/subnet-evm/precompile/contract"
	"github.com/ava-labs/subnet-evm/vmerrs"

	_ "embed"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
	"github.com/ethereum/go-ethereum/common"
)

const (
	// Gas costs for each function. These are set to 1 by default.
	// You should set a gas cost for each function in your contract.
	// Generally, you should not set gas costs very low as this may cause your network to be vulnerable to DoS attacks.
	// There are some predefined gas costs in contract/utils.go that you can use.
	VerifyGasCost uint64 = 1 /* SET A GAS COST HERE */
)

// CUSTOM CODE STARTS HERE
// Reference imports to suppress errors from unused imports. This code and any unnecessary imports can be removed.
var (
	_ = abi.JSON
	_ = errors.New
	_ = big.NewInt
	_ = vmerrs.ErrOutOfGas
	_ = common.Big0
)

// Singleton StatefulPrecompiledContract and signatures.
var (

	// CubicproofRawABI contains the raw ABI of Cubicproof contract.
	//go:embed contract.abi
	CubicproofRawABI string

	CubicproofABI = contract.ParseABI(CubicproofRawABI)

	CubicproofPrecompile = createCubicproofPrecompile()
)

type VerifyInput struct {
	X *big.Int
	Y *big.Int
}

// UnpackVerifyInput attempts to unpack [input] as VerifyInput
// assumes that [input] does not include selector (omits first 4 func signature bytes)
func UnpackVerifyInput(input []byte) (VerifyInput, error) {
	inputStruct := VerifyInput{}
	err := CubicproofABI.UnpackInputIntoInterface(&inputStruct, "verify", input)

	return inputStruct, err
}

// PackVerify packs [inputStruct] of type VerifyInput into the appropriate arguments for verify.
func PackVerify(inputStruct VerifyInput) ([]byte, error) {
	return CubicproofABI.Pack("verify", inputStruct.X, inputStruct.Y)
}

// PackVerifyOutput attempts to pack given result of type bool
// to conform the ABI outputs.
func PackVerifyOutput(result bool) ([]byte, error) {
	return CubicproofABI.PackOutput("verify", result)
}

// UnpackVerifyOutput attempts to unpack given [output] into the bool type output
// assumes that [output] does not include selector (omits first 4 func signature bytes)
func UnpackVerifyOutput(output []byte) (bool, error) {
	res, err := CubicproofABI.Unpack("verify", output)
	if err != nil {
		return false, err
	}
	unpacked := *abi.ConvertType(res[0], new(bool)).(*bool)
	return unpacked, nil
}

type CubicCircuit struct {
	// struct tags on a variable is optional
	// default uses variable name and secret visibility.
	X frontend.Variable `gnark:"x"`
	Y frontend.Variable `gnark:",public"`
}

func (circuit *CubicCircuit) Define(api frontend.API) error {
	x3 := api.Mul(circuit.X, circuit.X, circuit.X)
	api.AssertIsEqual(circuit.Y, api.Add(x3, circuit.X, 5))
	return nil
}

func verify(accessibleState contract.AccessibleState, caller common.Address, addr common.Address, input []byte, suppliedGas uint64, readOnly bool) (ret []byte, remainingGas uint64, err error) {
	if remainingGas, err = contract.DeductGas(suppliedGas, VerifyGasCost); err != nil {
		return nil, 0, err
	}
	// attempts to unpack [input] into the arguments to the VerifyInput.
	// Assumes that [input] does not include selector
	// You can use unpacked [inputStruct] variable in your code
	inputStruct, err := UnpackVerifyInput(input)
	if err != nil {
		return nil, remainingGas, err
	}

	// CUSTOM CODE STARTS HERE
	_ = inputStruct // CUSTOM CODE OPERATES ON INPUT

	var output bool // CUSTOM CODE FOR AN OUTPUT

	var circuit CubicCircuit
	ccs, _ := frontend.Compile(ecc.BN254, r1cs.NewBuilder, &circuit)

	// groth16 zkSNARK: Setup
	pk, vk, _ := groth16.Setup(ccs)

	// witness definition
	assignment := CubicCircuit{X: inputStruct.X, Y: inputStruct.Y}
	witness, _ := frontend.NewWitness(&assignment, ecc.BN254)
	publicWitness, _ := witness.Public()

	// groth16: Prove & Verify
	proof, _ := groth16.Prove(ccs, pk, witness)
	isInvalid := groth16.Verify(proof, vk, publicWitness)
	if isInvalid != nil {
		output = false
		return
	}
	output = true
	packedOutput, err := PackVerifyOutput(output)
	if err != nil {
		return nil, remainingGas, err
	}

	// Return the packed output and the remaining gas
	return packedOutput, remainingGas, nil
}

// createCubicproofPrecompile returns a StatefulPrecompiledContract with getters and setters for the precompile.

func createCubicproofPrecompile() contract.StatefulPrecompiledContract {
	var functions []*contract.StatefulPrecompileFunction

	abiFunctionMap := map[string]contract.RunStatefulPrecompileFunc{
		"verify": verify,
	}

	for name, function := range abiFunctionMap {
		method, ok := CubicproofABI.Methods[name]
		if !ok {
			panic(fmt.Errorf("given method (%s) does notexist in the ABI", name))
		}
		functions = append(functions, contract.NewStatefulPrecompileFunction(method.ID, function))
	}
	// Construct the contract with no fallback function.
	statefulContract, err := contract.NewStatefulPrecompileContract(nil, functions)
	if err != nil {
		panic(err)
	}
	return statefulContract
}

