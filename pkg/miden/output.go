package miden

import field "github.com/qredo/verifiable-oracles/pkg/goldilocks"

type Output struct {
	Stack         field.Vector `json:"stack"`
	OverflowAddrs []string     `json:"overflow_addrs"`
}
