package handlers

import (
	"net/http"

	"github.com/edwinbustillos/api-go-gin/database"
	"github.com/edwinbustillos/api-go-gin/models"
	"github.com/gin-gonic/gin"
)

func GetItems(c *gin.Context) {
	db := database.InitDB()
	defer db.Close()

	rows, err := db.Query(`SELECT * FROM items`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var items []models.Item
	for rows.Next() {
		var item models.Item
		err = rows.Scan(&item.ID, &item.Name, &item.Price)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		items = append(items, item)
	}

	if err = rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": items})
}

func GetItem(c *gin.Context) {
	db := database.InitDB()
	defer db.Close()

	var item models.Item
	query := `SELECT * FROM items WHERE id = $1`
	err := db.QueryRow(query, c.Param("id")).Scan(&item.ID, &item.Name, &item.Price)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": item})
}

func CreateItem(c *gin.Context) {
	db := database.InitDB()
	defer db.Close()

	var item models.Item
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	query := `INSERT INTO items (name, price) VALUES ($1, $2)`
	_, err := db.Query(query, item.Name, item.Price)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{"data": item})
	}

}

func UpdateItem(c *gin.Context) {
	db := database.InitDB()
	defer db.Close()

	var item models.Item
	query := `SELECT * FROM items WHERE id = $1`
	err := db.QueryRow(query, c.Param("id")).Scan(&item.ID, &item.Name, &item.Price)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updateQuery := `UPDATE items SET name = $1, price = $2 WHERE id = $3`
	_, err = db.Exec(updateQuery, item.Name, item.Price, item.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": item})
}

func DeleteItem(c *gin.Context) {
	db := database.InitDB()
	defer db.Close()

	var item models.Item
	query := `SELECT * FROM items WHERE id = $1`
	err := db.QueryRow(query, c.Param("id")).Scan(&item.ID, &item.Name, &item.Price)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	deleteQuery := `DELETE FROM items WHERE id = $1`
	_, err = db.Exec(deleteQuery, item.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": "Record deleted!"})
}
