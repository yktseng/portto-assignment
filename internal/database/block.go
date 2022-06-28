package database

import (
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/yktseng/portto-assignment/internal/block"
	"gorm.io/gorm/clause"
)

func (db *Database) SaveBlock(b *block.Block) error {
	// result := db.conn.Debug().CreateInBatches(logs, len(logs))
	result := db.conn.Clauses(
		clause.OnConflict{
			Columns:   []clause.Column{{Name: "block_hash"}},
			DoNothing: true,
		},
	).Create(b)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (db *Database) SetBlockDone(hash common.Hash) error {
	var b block.Block
	result := db.conn.Table("block").Where("block_hash = ?", hash.String()).Find(&b)
	if result.Error != nil {
		log.Panicln(result.Error)
		return result.Error
	}
	result = db.conn.Model(&block.Block{}).Where("block_hash = ?", hash.String()).Update("done", true)
	if result.Error != nil {
		log.Panicln(result.Error)
		return result.Error
	}
	return nil
}

func (db *Database) GetUnfinishedBlocks() ([]*big.Int, error) {
	var blocks []block.Block
	result := db.conn.Table("block").Where("done = ?", false).Find(&blocks)
	if result.Error != nil {
		log.Panicln(result.Error)
		return nil, result.Error
	}
	var unfinishedBlocks []*big.Int
	for _, b := range blocks {
		var n64 int64
		err := b.Num.AssignTo(&n64)
		if err != nil {
			return nil, err
		}
		unfinishedBlocks = append(unfinishedBlocks, big.NewInt(n64))
	}
	return unfinishedBlocks, nil
}

func (db *Database) GetLastRecordedBlock() (*big.Int, error) {
	var b block.Block
	result := db.conn.Table("block").Limit(1).Order("num DESC").Find(&b)
	if result.Error != nil {
		log.Panicln(result.Error)
		return nil, result.Error
	}
	var n64 int64
	err := b.Num.AssignTo(&n64)
	if err != nil {
		return nil, err
	}
	return big.NewInt(n64), nil
}
