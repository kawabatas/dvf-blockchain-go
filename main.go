package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/google/uuid"
)

const (
	ExitOK    = 0
	ExitError = 1
)

func main() {
	// port5000でサーバーを起動する
	flagAddr := flag.String("addr", ":5000", "host:port")
	flag.Parse()
	os.Exit(run(*flagAddr))
}

func run(addr string) int {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	errCh := make(chan error)

	// このノードのグローバルにユニークなアドレスを作る
	nodeUuid, err := uuid.NewRandom()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return ExitError
	}
	nodeIdentifier := strings.Replace(nodeUuid.String(), "-", "", -1)

	// ブロックチェーンクラスをインスタンス化する
	blockchain, err := InitBlockchain()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return ExitError
	}

	s := NewServer(addr, blockchain, nodeIdentifier)

	go func() {
		errCh <- s.Start()
	}()

	select {
	case err := <-errCh:
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return ExitError
		}
	case <-sigCh:
		if err := s.Stop(context.Background()); err != nil {
			fmt.Fprintln(os.Stderr, err)
			return ExitError
		}
	}
	return ExitOK
}
