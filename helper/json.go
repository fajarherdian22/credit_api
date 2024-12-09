package helper

import (
	"encoding/json"
	"fmt"

	"github.com/gin-gonic/gin"
)

// Encode Writer to JSON
func HandleEncodeWriteJson(c *gin.Context, WebResponse any) {
	c.Writer.Header().Add("Content-Type", "application/json")
	encoder := json.NewEncoder(c.Writer)
	err := encoder.Encode(WebResponse)
	IsError(err)
}

// Decode Request Body JSON
func HandleDecodeReqJson(c *gin.Context, dataStruct any) {
	fmt.Println(c.Request.Body)
	decoder := json.NewDecoder(c.Request.Body)
	err := decoder.Decode(dataStruct)
	IsError(err)
	defer c.Request.Body.Close()
}
