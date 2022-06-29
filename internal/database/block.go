package database

import (
	"context"
	"errors"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/yktseng/portto-assignment/internal/block"
	"gorm.io/gorm/clause"
)

type BlockQuery struct {
	Limit int    `json:"limit" form:"limit" binding:"omitempty,gte=1,lte=10000"`
	ID    string `json:"id" uri:"id"`
}

func (db *Database) SaveBlock(ctx context.Context, b *block.Block) error {
	// result := db.conn.Debug().CreateInBatches(logs, len(logs))
	result := db.conn.Clauses(
		clause.OnConflict{
			Columns:   []clause.Column{{Name: "block_hash"}},
			UpdateAll: true,
		},
	).Create(b)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (db *Database) SetBlockDone(ctx context.Context, hash common.Hash) error {
	var b block.Block
	result := db.conn.WithContext(ctx).Table("blocks").Where("block_hash = ?", hash.String()).Find(&b)
	if result.Error != nil {
		log.Panicln(result.Error)
		return result.Error
	}
	result = db.conn.WithContext(ctx).Model(&block.Block{}).Where("block_hash = ?", hash.String()).Update("done", true)
	if result.Error != nil {
		log.Panicln(result.Error)
		return result.Error
	}
	return nil
}

func (db *Database) SetBlockStable(ctx context.Context, hash common.Hash) error {
	var b block.Block
	result := db.conn.WithContext(ctx).Table("blocks").Where("block_hash = ?", hash.String()).Find(&b)
	if result.Error != nil {
		log.Panicln(result.Error)
		return result.Error
	}
	result = db.conn.WithContext(ctx).Model(&block.Block{}).Where("block_hash = ?", hash.String()).Update("stable", true)
	if result.Error != nil {
		log.Panicln(result.Error)
		return result.Error
	}
	return nil
}

func (db *Database) GetUnfinishedBlocks(ctx context.Context) ([]*big.Int, error) {
	var blocks []block.Block
	result := db.conn.WithContext(ctx).Table("blocks").Where("done = ?", false).Find(&blocks)
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

func (db *Database) GetMissingBlocks(ctx context.Context) ([]*big.Int, error) {
	var blocks []block.Block
	result := db.conn.Debug().WithContext(ctx).Raw("SELECT * DISTINCT num+1 FROM blocks WHERE num+1 NOT IN(SELECT DISTINCT num FROM blocks)")
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

func (db *Database) GetLastRecordedBlock(ctx context.Context) (*big.Int, error) {
	var b block.Block
	result := db.conn.WithContext(ctx).Table("blocks").Limit(1).Order("num DESC").Find(&b)
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

func (db *Database) Getblocks(ctx context.Context, q BlockQuery) ([]block.Block, error) {
	var blocks []block.Block
	query := db.conn.WithContext(ctx).Table("blocks")
	if q.Limit > 0 {
		query = query.Limit(q.Limit).Order("num desc")
	}
	result := query.Find(&blocks)
	if result.Error != nil {
		log.Panicln(result.Error)
		return nil, result.Error
	}
	return blocks, nil
}

func (db *Database) GetblockDetail(ctx context.Context, q BlockQuery) (*block.Block, error) {
	var block block.Block
	query := db.conn.WithContext(ctx).Table("blocks")
	if q.ID != "" {
		query = query.Where("block_hash = ?", q.ID)
	}
	result := query.Find(&block)
	if result.Error != nil {
		log.Panicln(result.Error)
		return nil, result.Error
	}
	if q.ID != block.Hash {
		return nil, errors.New("block not found")
	}
	var txs []string
	result = db.conn.WithContext(ctx).Table("txs").Select("tx_hash").Where("block_hash = ?", q.ID).Find(&txs)
	if result.Error != nil {
		log.Panicln(result.Error)
		return nil, result.Error
	}
	var hashes []common.Hash
	for _, tx := range txs {
		hashes = append(hashes, common.HexToHash(tx))
	}
	// fmt.Println(hashes)
	block.TXS = hashes
	return &block, nil
}
