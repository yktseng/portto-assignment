package transaction

import (
	"encoding/hex"
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/jackc/pgtype"
	"github.com/yktseng/portto-assignment/internal/logs"
)

// CREATE TABLE IF NOT EXISTS tx (
//   block_hash CHAR(64) NOT NULL,
//   tx_hash CHAR(64) PRIMARY KEY,
//   sender CHAR(40),
//   receiver CHAR(40) NOT NULL,
//   nonce INT,
//   tx_data TEXT,
//   amount bigint,
//   CONSTRAINT fk_block_hash
//     FOREIGN KEY(block_hash)
//       REFERENCES block(block_hash)
// );

type TX struct {
	BlockHash string
	Hash      string `gorm:"column:tx_hash"`
	From      string `gorm:"column:sender"`
	To        string `gorm:"column:receiver"`
	Nonce     int
	Data      string         `gorm:"column:tx_data;type:TEXT"`
	Value     pgtype.Numeric `gorm:"column:amount"`
	logs      []*logs.TXLog
}

// func PrintTransactionReceipt(t *types.Receipt) {
// 	fmt.Println("TX Hash", t.TxHash)
// 	fmt.Println("Logs", t.Logs)
// }

func From(tx *types.Transaction, chainID *big.Int) (string, error) {
	if msg, err := tx.AsMessage(types.LatestSignerForChainID(chainID), tx.GasPrice()); err != nil {
		return "", err
	} else {
		return msg.From().String(), nil
	}
}

func (t *TX) Logs() []*logs.TXLog {
	return t.logs
}

func FromGethTX(tx *types.Transaction, receipt *types.Receipt) (*TX, error) {
	from, err := From(tx, big.NewInt(97))
	if err != nil {
		return nil, err
	}
	txLogs := make([]*logs.TXLog, 0)
	for _, l := range receipt.Logs {
		txLog, err := logs.FromGethLog(l)
		if err != nil {
			return nil, err
		}
		txLogs = append(txLogs, txLog)
	}
	var toString string
	if tx.To() != nil {
		toString = tx.To().String()
	} else {
		toString = ""
	}
	t := &TX{
		BlockHash: receipt.BlockHash.String(),
		Hash:      tx.Hash().String(),
		To:        toString,
		From:      from,
		Nonce:     int(tx.Nonce()),
		Data:      hex.EncodeToString(tx.Data()),
		Value:     pgtype.Numeric{Int: tx.Value(), Status: pgtype.Present},
		logs:      txLogs,
	}
	return t, nil
}

func (TX) TableName() string {
	return "tx"
}
