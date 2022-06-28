package block

import (
	"encoding/json"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/jackc/pgtype"
)

// CREATE TABLE IF NOT EXISTS block (
// 	num bigint,
// 	block_hash CHAR(64) PRIMARY KEY,
// 	block_time TIMESTAMPTZ NOT NULL,
// 	parent_hash CHAR(64),
// 	created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
// 	updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
// );

type Block struct {
	Num        pgtype.Numeric `gorm:"type:NUMERIC"`
	Hash       string         `gorm:"column:block_hash"`
	Time       time.Time      `gorm:"column:block_time"`
	ParentHash string         `gorm:"column:parent_hash"`
	Done       bool
	TXS        []common.Hash `gorm:"-"`
}

func (b *Block) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Num        int64         `json:"num"`
		Hash       string        `json:"block_hash"`
		Time       time.Time     `json:"block_time"`
		ParentHash string        `json:"parent_hash"`
		TXS        []common.Hash `json:"transactions,omitempty"`
	}{
		Num:        b.Num.Int.Int64(),
		Hash:       b.Hash,
		Time:       b.Time,
		ParentHash: b.ParentHash,
		TXS: b.TXS,
	})
}

// func PrintBlock(b *types.Block) {
// 	fmt.Println("Block", b.Number())
// 	fmt.Println("-----------------")
// 	fmt.Println("Hash", b.Hash())
// 	fmt.Println("Parent Hash", b.ParentHash())
// 	fmt.Println("TX Hash", b.ParentHash())
// 	fmt.Println("Time", b.Time())
// 	fmt.Println("Uncles", b.Uncles())
// 	for i := 0; i < len(b.Transactions()); i++ {
// 		fmt.Println("TX", i, b.Transactions()[i].Hash())
// 	}
// 	fmt.Println("")
// }

func NewBlock(num *big.Int, hash, parentHash string,
	blockTime time.Time, txs []common.Hash) *Block {
	return &Block{
		Num:        pgtype.Numeric{Int: num, Status: pgtype.Present},
		Hash:       hash,
		Time:       blockTime,
		ParentHash: parentHash,
		TXS:        txs,
	}
}

func FromGethBlock(b *types.Block) (*Block, error) {
	unix := time.Unix(int64(b.Time()), 0)
	txs := make([]common.Hash, 0)
	for _, tx := range b.Transactions() {
		txs = append(txs, tx.Hash())
	}
	return NewBlock(b.Number(),
		b.Hash().String(), b.ParentHash().String(), unix, txs), nil
}

func (b *Block) Transactions() []common.Hash {
	return b.TXS
}

func (Block) TableName() string {
	return "block"
}
