package app

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PromoController struct {
	db *gorm.DB
}

func (pc *PromoController) CreatePromo(c *gin.Context) {
	promo := new(PromoModel)

	if err := c.ShouldBindJSON(promo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := pc.db.Create(promo).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": promo})
}

func (pc *PromoController) GetPromos(c *gin.Context) {
	var promos []PromoModel

	if err := pc.db.Limit(10).Offset(0).Order("id desc").Find(&promos).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": promos})
}

func (pc *PromoController) GetPromo(c *gin.Context) {
	promo := new(PromoModel)

	if err := c.ShouldBindUri(promo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := pc.db.First(&promo).Error; err != nil {
		code := http.StatusInternalServerError
		if err == gorm.ErrRecordNotFound {
			code = http.StatusNotFound
		}
		c.JSON(code, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": promo})
}

func (pc *PromoController) UpdatePromo(c *gin.Context) {
	promo := new(PromoModel)

	if err := c.ShouldBindUri(promo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.ShouldBindJSON(promo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := pc.db.Save(promo).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.AbortWithStatus(http.StatusNoContent)
}

func (pc *PromoController) DeletePromo(c *gin.Context) {
	promo := new(PromoModel)

	if err := c.ShouldBindUri(promo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := pc.db.First(&promo).Error; err != nil {
		code := http.StatusInternalServerError
		if err == gorm.ErrRecordNotFound {
			code = http.StatusNotFound
		}
		c.JSON(code, gin.H{"error": err.Error()})
		return
	}

	if err := pc.db.Delete(promo).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.AbortWithStatus(http.StatusNoContent)
}
