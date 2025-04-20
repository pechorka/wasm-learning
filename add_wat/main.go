package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/bytecodealliance/wasmtime-go/v31"
	"github.com/pechorka/gostdlib/pkg/errs"
)

func main() {
	if err := run(); err != nil {
		slog.Error("error running program", "error", err)
		os.Exit(1)
	}
}

func run() error {
	store := wasmtime.NewStore(wasmtime.NewEngine())

	wasm, err := wasmtime.Wat2Wasm(`
	(module
		(func $add (param $lhs i32) (param $rhs i32) (result i32)
			local.get $lhs
			local.get $rhs
			i32.add
		)
		(export "add" (func $add))
	)
		`)
	if err != nil {
		return errs.Wrap(err, "failed to compile wasm")
	}

	module, err := wasmtime.NewModule(store.Engine, wasm)
	if err != nil {
		return errs.Wrap(err, "failed to init module")
	}

	instance, err := wasmtime.NewInstance(store, module, nil)
	if err != nil {
		return errs.Wrap(err, "failed to init module instance")
	}

	add := instance.GetFunc(store, "add")
	if add == nil {
		return errs.Newf("failed to get add function")
	}

	res, err := add.Call(store, 13, 56)
	if err != nil {
		return errs.Wrap(err, "failed to run add")
	}

	fmt.Println(res)

	return nil
}
