package aptos

import (
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/aptos-labs/aptos-go-sdk"
)

func ConvertToRawTransaction(req *TransactionRequest) (*aptos.RawTransaction, error) {
	if req == nil {
		return nil, errors.New("transaction request cannot be nil")
	}

	// 1. convert sender address
	sender, err := AddressToAccountAddress(req.Sender)
	if err != nil {
		return nil, fmt.Errorf("invalid sender address: %w", err)
	}

	// 2. convert module address
	coinTypeAddress, err := CoinAddressToAccountAddress(req.Payload.Payload.Module.Address)
	if err != nil {
		return nil, err
	}
	moduleId := aptos.ModuleId{
		Address: coinTypeAddress,
		Name:    req.Payload.Payload.Module.Name,
	}

	// 3. deal TypeArgs
	//typeArgs := make([]aptos.TypeTag, len(req.Payload.Payload.ArgTypes))
	//for i, argType := range req.Payload.Payload.ArgTypes {
	//	typeTag, err := ParseTypeTag(argType)
	//	if err != nil {
	//		return nil, fmt.Errorf("failed to parse type arg %s: %w", argType, err)
	//	}
	//	typeArgs[i] = typeTag
	//}

	// 4. deal Args
	args, err := convertFunctionArgs(req.Payload.Payload.Args, req.Payload.Payload.ArgTypes)
	if err != nil {
		return nil, fmt.Errorf("failed to convert function arguments: %w", err)
	}

	// 5. create EntryFunction
	entryFunction := &aptos.EntryFunction{
		Module:   moduleId,
		Function: req.Payload.Payload.Function,
		ArgTypes: make([]aptos.TypeTag, 0),
		Args:     args,
	}

	// 6. create TransactionPayload
	payload := aptos.TransactionPayload{
		Payload: entryFunction,
	}

	// 7. create RawTransaction
	rawTx := &aptos.RawTransaction{
		Sender:                     sender,
		SequenceNumber:             req.SequenceNumber,
		Payload:                    payload,
		MaxGasAmount:               req.MaxGasAmount,
		GasUnitPrice:               req.GasUnitPrice,
		ExpirationTimestampSeconds: req.ExpirationTimestampSeconds,
		ChainId:                    req.ChainId,
	}

	return rawTx, nil
}

func convertFunctionArgs(args []string, argTypes []string) ([][]byte, error) {
	if len(args) != len(argTypes) {
		return nil, errors.New("argument count does not match type count")
	}

	result := make([][]byte, len(args))
	for i, arg := range args {
		convertedArg, err := convertSingleArg(arg, argTypes[i])
		if err != nil {
			return nil, fmt.Errorf("failed to convert argument %d: %w", i, err)
		}
		result[i] = convertedArg
	}
	return result, nil
}

func convertSingleArg(arg string, argType string) ([]byte, error) {
	switch argType {
	case "address":
		return convertAddressArg(arg)
	case "u64":
		return convertU64Arg(arg)
	default:
		return nil, fmt.Errorf("unsupported argument type: %s", argType)
	}
}

func convertAddressArg(arg string) ([]byte, error) {
	addr, err := hex.DecodeString(strings.TrimPrefix(arg, "0x"))
	if err != nil {
		return nil, fmt.Errorf("failed to decode address: %w", err)
	}

	if len(addr) < 32 {
		paddedAddr := make([]byte, 32)
		copy(paddedAddr[32-len(addr):], addr)
		addr = paddedAddr
	}
	return addr, nil
}

func convertU64Arg(arg string) ([]byte, error) {
	value, err := strconv.ParseUint(arg, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse u64: %w", err)
	}
	buf := make([]byte, 8)
	binary.LittleEndian.PutUint64(buf, value)
	return buf, nil
}

func ParseTypeTag(typeStr string) (aptos.TypeTag, error) {
	var inner aptos.TypeTagImpl
	switch typeStr {
	case "bool":
		inner = &aptos.BoolTag{}
	case "u8":
		inner = &aptos.U8Tag{}
	case "u16":
		inner = &aptos.U16Tag{}
	case "u32":
		inner = &aptos.U32Tag{}
	case "u64":
		inner = &aptos.U64Tag{}
	case "u128":
		inner = &aptos.U128Tag{}
	case "u256":
		inner = &aptos.U256Tag{}
	case "address":
		inner = &aptos.AddressTag{}
	case "signer":
		inner = &aptos.SignerTag{}
	default:
		if strings.HasPrefix(typeStr, "vector<") && strings.HasSuffix(typeStr, ">") {
			innerTypeStr := strings.TrimPrefix(strings.TrimSuffix(typeStr, ">"), "vector<")
			innerType, err := ParseTypeTag(innerTypeStr)
			if err != nil {
				return aptos.TypeTag{}, fmt.Errorf("invalid vector type: %w", err)
			}
			inner = &aptos.VectorTag{TypeParam: innerType}
		} else if strings.Contains(typeStr, "::") {
			structTag, err := parseStructTag(typeStr)
			if err != nil {
				return aptos.TypeTag{}, err
			}
			inner = structTag
		} else {
			return aptos.TypeTag{}, fmt.Errorf("unsupported type: %s", typeStr)
		}
	}

	return aptos.NewTypeTag(inner), nil
}

func parseStructTag(typeStr string) (*aptos.StructTag, error) {
	parts := strings.Split(typeStr, "::")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid struct type format: %s", typeStr)
	}

	addr, err := AddressToAccountAddress(parts[0])
	if err != nil {
		return nil, fmt.Errorf("invalid address in struct type: %w", err)
	}

	name := parts[2]
	var typeParams []aptos.TypeTag
	if idx := strings.Index(name, "<"); idx != -1 {
		if !strings.HasSuffix(name, ">") {
			return nil, fmt.Errorf("invalid type parameters format in: %s", name)
		}

		baseType := name[:idx]
		typeParamsStr := name[idx+1 : len(name)-1]

		if typeParamsStr != "" {
			paramStrs := strings.Split(typeParamsStr, ",")
			typeParams = make([]aptos.TypeTag, len(paramStrs))
			for i, param := range paramStrs {
				paramType, err := ParseTypeTag(strings.TrimSpace(param))
				if err != nil {
					return nil, fmt.Errorf("invalid type parameter: %w", err)
				}
				typeParams[i] = paramType
			}
		}
		name = baseType
	}

	return &aptos.StructTag{
		Address:    addr,
		Module:     parts[1],
		Name:       name,
		TypeParams: typeParams,
	}, nil
}
