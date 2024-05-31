package fhevm

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/zama-ai/fhevm-go/fhevm/tfhe"
	"pgregory.net/rapid"
)

func TestTeeDecryptRun(t *testing.T) {
	signature := "teeDecrypt(uint256)"
	rapid.Check(t, func(t *rapid.T) {
		testcases := []struct {
			typ      tfhe.FheUintType
			expected uint64
		}{
			{tfhe.FheUint4, uint64(rapid.Uint8Range(0, 15).Draw(t, "expected"))},
			{tfhe.FheUint8, uint64(rapid.Uint8().Draw(t, "expected"))},
			{tfhe.FheUint16, uint64(rapid.Uint16().Draw(t, "expected"))},
			{tfhe.FheUint32, uint64(rapid.Uint32().Draw(t, "expected"))},
			{tfhe.FheUint64, rapid.Uint64().Draw(t, "expected")},
		}
		for _, tc := range testcases {
			depth := 1
			environment := newTestEVMEnvironment()
			environment.depth = depth
			addr := common.Address{}
			readOnly := false
			ct, err := importTeePlaintextToEVM(environment, depth, tc.expected, tc.typ)
			if err != nil {
				t.Fatalf(err.Error())
			}

			input := toLibPrecompileInput(signature, false, ct.GetHash())
			out, err := TeeLibRun(environment, addr, addr, input, readOnly)
			if err != nil {
				t.Fatalf(err.Error())
			}

			result := new(big.Int).SetBytes(out).Uint64()
			if result != tc.expected {
				t.Fatalf("incorrect result, expected=%d, got=%d", tc.expected, result)
			}
		}
	})
}
