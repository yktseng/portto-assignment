package database

import (
	"context"
	"errors"
	"log"

	"github.com/yktseng/portto-assignment/internal/logs"
	transaction "github.com/yktseng/portto-assignment/internal/tx"
	"gorm.io/gorm/clause"
)

type TXQuery struct {
	TXHash string `json:"tx_hash" uri:"tx_hash" binding:"required"`
}

func (db *Database) SaveTxs(ctx context.Context, txs []*transaction.TX) error {
	// result := db.conn.Debug().CreateInBatches(logs, len(logs))
	result := db.conn.WithContext(ctx).Clauses(
		clause.OnConflict{
			Columns:   []clause.Column{{Name: "tx_hash"}},
			UpdateAll: true,
		},
	).Create(txs)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (db *Database) GetTXDetail(ctx context.Context, q TXQuery) (*transaction.TX, error) {
	var tx transaction.TX
	query := db.conn.WithContext(ctx).Table("txs")
	if q.TXHash != "" {
		query = query.Where("tx_hash = ?", q.TXHash)
	}
	result := query.Find(&tx)
	if result.Error != nil {
		log.Panicln(result.Error)
		return nil, result.Error
	}
	if tx.Hash != q.TXHash {
		// not found
		return nil, errors.New("tx not found")
	}
	var logs []*logs.TXLog
	result = db.conn.WithContext(ctx).Table("logs").Where("tx_hash = ?", q.TXHash).Find(&logs)
	if result.Error != nil {
		log.Panicln(result.Error)
		return nil, result.Error
	}
	tx.Logs = logs
	return &tx, nil
}
