package handler

import (
	"context"
	"database/sql"
	"math/big"
	"net/http"

	blockchain "goledger-challenge/contract"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	Blockchain *blockchain.Client
	DB         *sql.DB
}

type SetRequest struct {
	Value string `json:"value" binding:"required"`
}

func (h *Handler) SetValue(c *gin.Context) {
	var req SetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	val, ok := new(big.Int).SetString(req.Value, 10)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid value"})
		return
	}
	txHash, err := h.Blockchain.SetValue(context.Background(), val)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"tx_hash": txHash})
}

func (h *Handler) GetValue(c *gin.Context) {
	val, err := h.Blockchain.GetValue(context.Background())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"value": val.String()})
}

func (h *Handler) SyncValue(c *gin.Context) {
	val, err := h.Blockchain.GetValue(context.Background())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	_, err = h.DB.Exec("UPDATE storage SET value = $1 WHERE id = 1", val.String())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DB update failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"synced_value": val.String()})
}

func (h *Handler) CheckValue(c *gin.Context) {
	var dbValue string
	err := h.DB.QueryRow("SELECT value FROM storage WHERE id = 1").Scan(&dbValue)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DB read failed"})
		return
	}
	chainValue, err := h.Blockchain.GetValue(context.Background())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	isEqual := (dbValue == chainValue.String())
	c.JSON(http.StatusOK, gin.H{"equal": isEqual})
}
