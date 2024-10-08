package rpc

import (
	"context"
	"encoding/json"
	"os"
	"testing"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/utils"
	"github.com/stretchr/testify/require"
)

// TestTransactionTrace is a function that tests the TransactionTrace function.
//
// It sets up the necessary test configuration and expected response. Then it performs a series of test sets,
// each with a different transaction hash. For each test set, it calls the TransactionTrace function and compares
// the response with the expected response.
//
// Parameters:
// - t: the testing object for running the test cases
// Returns:
//
//	none
func TestTransactionTrace(t *testing.T) {
	testConfig := beforeEach(t)

	var expectedResp InvokeTxnTrace
	if testEnv == "mock" {
		var rawjson struct {
			Result InvokeTxnTrace `json:"result"`
		}
		expectedrespRaw, err := os.ReadFile("./tests/trace/sepoliaInvokeTrace_0x6a4a9c4f1a530f7d6dd7bba9b71f090a70d1e3bbde80998fde11a08aab8b282.json")
		require.NoError(t, err, "Error ReadFile for TestTraceTransaction")

		err = json.Unmarshal(expectedrespRaw, &rawjson)
		require.NoError(t, err, "Error unmarshalling testdata TestTraceTransaction")

		txnTrace, err := json.Marshal(rawjson.Result)
		require.NoError(t, err, "Error unmarshalling testdata TestTraceTransaction")
		require.NoError(t, json.Unmarshal(txnTrace, &expectedResp))
	}

	type testSetType struct {
		TransactionHash *felt.Felt
		ExpectedResp    *InvokeTxnTrace
		ExpectedError   *RPCError
	}
	testSet := map[string][]testSetType{
		"mock": {
			testSetType{
				TransactionHash: utils.TestHexToFelt(t, "0x6a4a9c4f1a530f7d6dd7bba9b71f090a70d1e3bbde80998fde11a08aab8b282"),
				ExpectedResp:    &expectedResp,
				ExpectedError:   nil,
			},
			testSetType{
				TransactionHash: utils.TestHexToFelt(t, "0xc0ffee"),
				ExpectedResp:    nil,
				ExpectedError:   ErrHashNotFound,
			},
			testSetType{
				TransactionHash: utils.TestHexToFelt(t, "0xf00d"),
				ExpectedResp:    nil,
				ExpectedError: &RPCError{
					Code:    10,
					Message: "No trace available for transaction",
					Data:    "REJECTED",
				},
			},
		},
		"devnet":  {},
		"mainnet": {},
	}[testEnv]

	for _, test := range testSet {
		resp, err := testConfig.provider.TraceTransaction(context.Background(), test.TransactionHash)
		if err != nil {
			require.Equal(t, test.ExpectedError, err)
		} else {
			invokeTrace := resp.(InvokeTxnTrace)
			require.Equal(t, invokeTrace, *test.ExpectedResp)
		}
	}
}

// TestSimulateTransaction is a function that tests the SimulateTransaction function in the codebase.
//
// It sets up the necessary test configuration and variables, and then performs a series of tests based on the test environment.
// The function reads input data from JSON files and performs JSON unmarshalling to set the values of the simulateTxIn and expectedResp variables.
// It then iterates over the testSet map, calling the SimulateTransactions function with the appropriate parameters and checking the response against the expected response.
// The function uses the testing.T type to report any errors or failures during the test execution.
//
// Parameters:
// - t: the testing object for running the test cases
// Returns:
//
//	none
func TestSimulateTransaction(t *testing.T) {
	testConfig := beforeEach(t)

	var simulateTxIn SimulateTransactionInput
	var expectedResp SimulateTransactionOutput
	if testEnv == "mainnet" {
		simulateTxnRaw, err := os.ReadFile("./tests/trace/simulateInvokeTx.json")
		require.NoError(t, err, "Error ReadFile simulateInvokeTx")

		err = json.Unmarshal(simulateTxnRaw, &simulateTxIn)
		require.NoError(t, err, "Error unmarshalling simulateInvokeTx")

		expectedrespRaw, err := os.ReadFile("./tests/trace/simulateInvokeTxResp.json")
		require.NoError(t, err, "Error ReadFile simulateInvokeTxResp")

		err = json.Unmarshal(expectedrespRaw, &expectedResp)
		require.NoError(t, err, "Error unmarshalling simulateInvokeTxResp")
	}

	type testSetType struct {
		SimulateTxnInput SimulateTransactionInput
		ExpectedResp     SimulateTransactionOutput
	}
	testSet := map[string][]testSetType{
		"devnet":  {},
		"mock":    {},
		"testnet": {},
		"mainnet": {testSetType{
			SimulateTxnInput: simulateTxIn,
			ExpectedResp:     expectedResp,
		}},
	}[testEnv]

	for _, test := range testSet {

		resp, err := testConfig.provider.SimulateTransactions(
			context.Background(),
			test.SimulateTxnInput.BlockID,
			test.SimulateTxnInput.Txns,
			test.SimulateTxnInput.SimulationFlags)
		require.NoError(t, err)
		require.Equal(t, test.ExpectedResp.Txns, resp)
	}
}

// TestTraceBlockTransactions tests the TraceBlockTransactions function.
//
// It sets up the test configuration and expected response. It then iterates
// through the test set, calling TraceBlockTransactions with the provided
// block hash. It checks if there is an error and compares the response with
// the expected response.
//
// Parameters:
// - t: the testing object for running the test cases
// Returns:
//
//	none
func TestTraceBlockTransactions(t *testing.T) {
	testConfig := beforeEach(t)

	var expectedResp []Trace
	if testEnv == "mock" {
		var rawjson struct {
			Result []Trace `json:"result"`
		}
		expectedrespRaw, err := os.ReadFile("./tests/trace/sepoliaBlockTrace_0x42a4c6a4c3dffee2cce78f04259b499437049b0084c3296da9fbbec7eda79b2.json")
		require.NoError(t, err, "Error ReadFile for TestTraceBlockTransactions")

		err = json.Unmarshal(expectedrespRaw, &rawjson)
		require.NoError(t, err, "Error unmarshalling testdata TestTraceBlockTransactions")
		expectedResp = rawjson.Result
	}

	type testSetType struct {
		BlockID      BlockID
		ExpectedResp []Trace
		ExpectedErr  *RPCError
	}
	testSet := map[string][]testSetType{
		"devnet":  {}, // devenet doesn't support TraceBlockTransactions https://0xspaceshard.github.io/starknet-devnet/docs/guide/json-rpc-api#trace-api
		"mainnet": {},
		"mock": {
			testSetType{
				BlockID:      BlockID{Hash: utils.TestHexToFelt(t, "0x42a4c6a4c3dffee2cce78f04259b499437049b0084c3296da9fbbec7eda79b2")},
				ExpectedResp: expectedResp,
				ExpectedErr:  nil,
			},
			testSetType{
				BlockID:      BlockID{Hash: utils.TestHexToFelt(t, "0x0")},
				ExpectedResp: nil,
				ExpectedErr:  ErrBlockNotFound,
			}},
	}[testEnv]

	for _, test := range testSet {
		resp, err := testConfig.provider.TraceBlockTransactions(context.Background(), test.BlockID)

		if err != nil {
			require.Equal(t, test.ExpectedErr, err)
		} else {
			require.Equal(t, test.ExpectedResp, resp)
		}

	}
}
