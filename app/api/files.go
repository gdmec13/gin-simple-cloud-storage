package api

import (
	"github.com/gin-gonic/gin"
	"simple-cloud-storage/app/api/request"
	"simple-cloud-storage/pkg/response"
)

func Files(c *gin.Context) {
	var R request.FilesStruct
	if err := c.ShouldBind(&R); err != nil {
		response.FailWithError(err, -1, c)
	} else {
		response.SuccessDetailed(R, "success", c)
	}
}
