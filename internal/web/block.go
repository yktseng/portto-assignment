package web

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yktseng/portto-assignment/internal/database"
)

func (s *Server) getBlocks(c *gin.Context) {
	var query database.BlockQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(http.StatusBadRequest,
			gin.H{
				"error":   "invalid query",
				"message": "limit should be between 1 ~ 10000",
			})
		return
	}
	if query.Limit == 0 {
		query.Limit = 1
	}
	blocks, err := s.db.Getblocks(c.Request.Context(), query)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK,
			gin.H{
				"error":   "failed to fetch for blocks",
				"message": err.Error(),
			})
		return
	}
	c.JSON(http.StatusOK, blocks)
}

func (s *Server) getBlockByID(c *gin.Context) {
	var query database.BlockQuery
	if err := c.ShouldBindUri(&query); err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(http.StatusBadRequest,
			gin.H{
				"error":   "invalid query",
				"message": err,
			})
		return
	}
	// fmt.Println(query)
	block, err := s.db.GetblockDetail(c.Request.Context(), query)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK,
			gin.H{
				"error":   "failed to fetch for blocks",
				"message": err.Error(),
			})
		return
	}
	// fmt.Println(block)
	c.JSON(http.StatusOK, block)
}

func (s *Server) getTXByHash(c *gin.Context) {
	var query database.TXQuery
	if err := c.ShouldBindUri(&query); err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(http.StatusBadRequest,
			gin.H{
				"error":   "invalid query",
				"message": err,
			})
		return
	}
	// fmt.Println(query)
	tx, err := s.db.GetTXDetail(c.Request.Context(), query)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK,
			gin.H{
				"error":   "failed to fetch for blocks",
				"message": err.Error(),
			})
		return
	}
	c.JSON(http.StatusOK, tx)
}
