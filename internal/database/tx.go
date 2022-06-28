package database

import (
	"context"

	"github.com/yktseng/portto-assignment/internal/tx"
	"gorm.io/gorm/clause"
)

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
