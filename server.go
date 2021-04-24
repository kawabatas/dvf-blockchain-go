package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type Server struct {
	server         *http.Server
	blockchain     *Blockchain
	nodeIdentifier string
}

func NewServer(addr string, bc *Blockchain, nodeIdentifier string) *Server {
	return &Server{
		server:         &http.Server{Addr: addr},
		blockchain:     bc,
		nodeIdentifier: nodeIdentifier,
	}
}

func (s *Server) Start() error {
	s.initHandlers()
	err := s.server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	err := s.server.Shutdown(ctx)
	if err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

func (s *Server) initHandlers() {
	mux := http.NewServeMux()
	s.server.Handler = mux

	// メソッドはPOSTで/transactions/newエンドポイントを作る。メソッドはPOSTなのでデータを送信する
	mux.HandleFunc("/transactions/new", s.HandleNewTransactions)
	// メソッドはGETで/mineエンドポイントを作る
	mux.HandleFunc("/mine", s.HandleMine)
	// メソッドはGETで、フルのブロックチェーンをリターンする/chainエンドポイントを作る
	mux.HandleFunc("/chain", s.HandleFullChain)
}

type NewTransactionsResponse struct {
	Message string `json:"message"`
}

func (s *Server) HandleNewTransactions(w http.ResponseWriter, r *http.Request) {
	var transaction Transaction
	if err := json.NewDecoder(r.Body).Decode(&transaction); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// POSTされたデータに必要なデータがあるかを確認
	if transaction.Sender == "" || transaction.Recipient == "" || transaction.Amount == 0 {
		http.Error(w, "Missing values", http.StatusBadRequest)
		return
	}

	// 新しいトランザクションを作る
	index := s.blockchain.NewTransaction(transaction.Sender, transaction.Recipient, transaction.Amount)

	if err := json.NewEncoder(w).Encode(&NewTransactionsResponse{
		Message: fmt.Sprintf("トランザクションはブロック %d に追加されました", index),
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

type MineResponse struct {
	Message      string         `json:"message"`
	Index        int            `json:"index"`
	Transactions []*Transaction `json:"transactions"`
	Proof        int            `json:"proof"`
	PreviousHash string         `json:"previous_hash"`
}

func (s *Server) HandleMine(w http.ResponseWriter, r *http.Request) {
	// 次のプルーフを見つけるためプルーフ・オブ・ワークアルゴリズムを使用する
	lastBlock := s.blockchain.LastBlock()
	proof := ProofOfWork(lastBlock.Proof)

	// プルーフを見つけたことに対する報酬を得る
	// 送信者は、採掘者が新しいコインを採掘したことを表すために"0"とする
	s.blockchain.NewTransaction(
		"0",
		s.nodeIdentifier,
		1,
	)

	// チェーンに新しいブロックを加えることで、新しいブロックを採掘する
	block, err := s.blockchain.NewBlock(proof, "")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(&MineResponse{
		Message:      "新しいブロックを採掘しました",
		Index:        block.Index,
		Transactions: block.Transactions,
		Proof:        block.Proof,
		PreviousHash: block.PreviousHash,
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

type FullChainResponse struct {
	Chain  []*Block `json:"chain"`
	Length int      `json:"length"`
}

func (s *Server) HandleFullChain(w http.ResponseWriter, r *http.Request) {
	if err := json.NewEncoder(w).Encode(&FullChainResponse{
		Chain:  s.blockchain.Chain,
		Length: len(s.blockchain.Chain),
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
