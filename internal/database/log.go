package database

import (
	"context"

	"github.com/yktseng/portto-assignment/internal/logs"
	"gorm.io/gorm/clause"
)

func (db *Database) SaveLogs(ctx context.Context, logs []*logs.TXLog) error {
	// result := db.conn.Debug().CreateInBatches(logs, len(logs))
	if len(logs) == 0 {
		return nil
	}
	// for _, log := range logs {
	// 	fmt.Println(log)
	// }
	result := db.conn.WithContext(ctx).Clauses(
		clause.OnConflict{
			Columns:   []clause.Column{{Name: "tx_hash"}, {Name: "log_id"}},
			DoUpdates: clause.AssignmentColumns([]string{"data"}),
		},
	).Create(logs)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
