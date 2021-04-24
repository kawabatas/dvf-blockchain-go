package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type Blockchain struct {
	// Goでは配列は順序づけされている
	Chain               []*Block
	CurrentTransactions []*Transaction
}

type Block struct {
	Index        int            `json:"index"`
	Timestamp    time.Time      `json:"timestamp"`
	Transactions []*Transaction `json:"transactions"`
	Proof        int            `json:"proof"`
	PreviousHash string         `json:"previous_hash"`
}

type Transaction struct {
	Sender    string `json:"sender"`
	Recipient string `json:"recipient"`
	Amount    int    `json:"amount"`
}

func InitBlockchain() (*Blockchain, error) {
	bc := &Blockchain{}
	// ジェネシスブロックを作る
	if _, err := bc.NewBlock(100, "1"); err != nil {
		return nil, err
	}
	return bc, nil
}

// 新しいブロックを作り、チェーンに加える
//
// ブロックチェーンに新しいブロックを作る
// @param proof: プルーフ・オブ・ワークアルゴリズムから得られるプルーフ
// @param previousHash: 前のブロックのハッシュ
// @return: 新しいブロック
// @return: エラー
func (bc *Blockchain) NewBlock(proof int, previousHash string) (*Block, error) {
	var prevHash string
	var err error
	if len(previousHash) > 0 {
		prevHash = previousHash
	} else {
		prevHash, err = Hash(bc.LastBlock())
		if err != nil {
			return nil, err
		}
	}
	b := &Block{
		Index:        len(bc.Chain) + 1,
		Timestamp:    time.Now(),
		Transactions: bc.CurrentTransactions,
		Proof:        proof,
		PreviousHash: prevHash,
	}
	// 現在のトランザクションリストをリセット
	bc.CurrentTransactions = []*Transaction{}
	bc.Chain = append(bc.Chain, b)
	return b, nil
}

// 新しいトランザクションをリストに加える
//
// 次に採掘されるブロックに加える新しいトランザクションを作る
// @param sender: 送信者のアドレス
// @param recipient: 受信者のアドレス
// @param amount: 量
// @return: このトランザクションを含むブロックのアドレス
func (bc *Blockchain) NewTransaction(sender, recipient string, amount int) int {
	bc.CurrentTransactions = append(bc.CurrentTransactions, &Transaction{
		Sender:    sender,
		Recipient: recipient,
		Amount:    amount,
	})
	return bc.LastBlock().Index + 1
}

// ブロックをハッシュ化する
//
// ブロックの SHA-256 ハッシュを作る
// @param b: ブロック
// @return: ハッシュ
// @return: エラー
func Hash(b *Block) (string, error) {
	// Goでは配列は順序づけされている
	block_bytes, err := json.Marshal(b)
	if err != nil {
		return "", err
	}
	hash := sha256.Sum256(block_bytes)
	return hex.EncodeToString(hash[:]), nil
}

// チェーンの最後のブロックをリターンする
func (bc *Blockchain) LastBlock() *Block {
	return bc.Chain[len(bc.Chain)-1]
}

// シンプルなプルーフ・オブ・ワークのアルゴリズム
// - hash(pp')の最初の4つが0となるような p' を探す
// - p は1つ前のブロックのプルーフ、p' は新しいブロックのプルーフ
// @param last_proof
// @return
func ProofOfWork(lastProof int) int {
	proof := 0
	for !ValidProof(lastProof, proof) {
		proof += 1
	}
	return proof
}

// プルーフが正しいかを確認する。hash(last_proof, proof)の最初の4つが0となっているか
// @param last_proof: 前のプルーフ
// @param proof: 現在のプルーフ
// @return: 正しければ true、そうでなければ false
func ValidProof(lastProof, proof int) bool {
	guess := fmt.Sprintf("%d%d", lastProof, proof)
	guessHash := sha256.Sum256([]byte(guess))
	guessHashString := hex.EncodeToString(guessHash[:])
	return strings.HasPrefix(guessHashString, "0000")
}
