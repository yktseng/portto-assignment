package logs

import (
	"encoding/hex"

	"github.com/ethereum/go-ethereum/core/types"
)

// CREATE TABLE IF NOT EXISTS log (
//   tx_hash CHAR(64),
//   id int,
//   data TEXT,
//   CONSTRAINT fk_tx_hash
//     FOREIGN KEY(tx_hash)
//       REFERENCES tx(tx_hash),
//   PRIMARY KEY (tx_hash, id)
// );

type TXLog struct {
	TXHash string `gorm:"column:tx_hash"`
	Index  int    `gorm:"column:log_id"`
	Data   string `gorm:"type:TEXT"`
}

func FromGethLog(l *types.Log) (*TXLog, error) {
	return &TXLog{
		TXHash: l.TxHash.String(),
		Index: int(l.Index),
		Data: hex.EncodeToString(l.Data),
	}, nil
}

func (TXLog) TableName() string {
	return "log"
}
