package handlers

import (
	"blog/data/db"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Base(c *gin.Context){
	c.JSON(http.StatusOK,gin.H{"data":"ok!"})
}
func Get(c *gin.Context){
	id := c.Param("id")
	u,err := db.GetUser(id)
	if err != nil{
		c.AbortWithStatusJSON(http.StatusBadRequest,gin.H{"data":err.Error()})
	}
	c.JSON(http.StatusOK,gin.H{"data":u})
}
func Delete(c *gin.Context){
	id := c.Param("id")
	err := db.DeleteUser(id)
	if err != nil{
		c.AbortWithStatusJSON(http.StatusBadRequest,gin.H{"data":err.Error()})
	}
	c.JSON(http.StatusOK,gin.H{"data":"ok!"})
}
func Post(c *gin.Context){
	username := c.PostForm("username")
	password:= c.PostForm("password")
	u,err := db.PostUser(username,password)
	if err != nil{
		c.AbortWithStatusJSON(http.StatusBadRequest,gin.H{"data":err.Error()})
	}
	c.JSON(http.StatusOK,gin.H{"data":u})
}
func Put(c *gin.Context){
	id := c.PostForm("id")
	username := c.PostForm("username")
	password:= c.PostForm("password")
	err := db.PutUser(id,username,password)
	if err != nil{
		c.AbortWithStatusJSON(http.StatusBadRequest,gin.H{"data":err.Error()})
	}
	c.JSON(http.StatusOK,gin.H{"data":"ok!"})
}