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

	wasm, err := os.ReadFile("add.wasm")
	if err != nil {
		return errs.Wrap(err, "failed to compile wasm")
	}

	module, err := wasmtime.NewModule(store.Engine, wasm)
	if err != nil {
		return errs.Wrap(err, "failed to init module")
	}

	linker := wasmtime.NewLinker(store.Engine)
	err = linker.DefineWasi()
	if err != nil {
		return errs.Wrap(err, "failed to define wasi")
	}
	store.SetWasi(wasmtime.NewWasiConfig())

	instance, err := linker.Instantiate(store, module)
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
