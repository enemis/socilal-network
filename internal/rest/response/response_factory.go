package response

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"social-network-otus/internal/app_error"
	"social-network-otus/internal/logger"
)

type F = map[string]interface{}

type ResponseFactory struct {
	Logger logger.LoggerInterface
}

func NewResponseFactory(log logger.LoggerInterface) *ResponseFactory {
	return &ResponseFactory{Logger: log}
}

func (res *ResponseFactory) Created(c *gin.Context, new interface{}) {
	c.JSON(http.StatusCreated, F{"object": new})
}

func (res *ResponseFactory) Ok(c *gin.Context, result interface{}) {
	c.JSON(http.StatusOK, result)
}

func (res *ResponseFactory) OkWithMessage(c *gin.Context, message string) {
	c.JSON(http.StatusOK, F{"message": message})
}

func (res *ResponseFactory) Unauthorised(c *gin.Context) {
	c.AbortWithStatusJSON(http.StatusUnauthorized, F{"error": "unauthorized"})
}

func (res *ResponseFactory) Expired(c *gin.Context, err error) {
	c.AbortWithStatusJSON(http.StatusUnauthorized, F{"error": "token expired"})
}

func (res *ResponseFactory) Forbidden(c *gin.Context, err error) {
	res.Logger.Error(err.Error(), err, nil)
	c.AbortWithStatusJSON(http.StatusForbidden, F{"error": "access denied"})
}

func (res *ResponseFactory) NotFound(c *gin.Context, err error) {
	res.Logger.Error(err.Error(), err, nil)
	c.AbortWithStatusJSON(http.StatusNotFound, F{"error": "object not found"})
}

func (res *ResponseFactory) BadRequest(c *gin.Context, errors F) {
	c.AbortWithStatusJSON(http.StatusBadRequest, F{"errors": errors})
}

func (res *ResponseFactory) FromAppError(c *gin.Context, appError *app_error.AppError, field *string) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in f", r)
		}
	}()
	switch appError.Status() {
	case http.StatusNotFound:
		res.NotFound(c, errors.New(appError.Error()))
	case http.StatusUnauthorized:
		res.Unauthorised(c)
	case http.StatusForbidden:
		res.Forbidden(c, appError.OriginalError())
	case http.StatusBadRequest:
		res.BadRequest(c, F{*field: appError.Error()})
	default:
		res.InternalServerError(c, appError.OriginalError())
	}
}

//
//func (res *ResponseFactory) ErrorResponse(statusCode uint, c *gin.Context, errors interface{}, err error) {
//	if err != nil {
//		res.Logger.Error("bad request", err, map[string]interface{}{})
//	}
//
//	w.WriteHeader(int(statusCode))
//	_ = json.NewEncoder(w).Encode(map[string]interface{}{"errors": errors})
//}
//
//func (res *ResponseFactory) WrappedError(c *gin.Context, err error) {
//	processedError, status := res.errorsResolver.ProcessErrorAndStatus(err)
//
//	if status == 0 {
//		res.InternalServerError(w, r, err)
//		return
//	}
//
//	w.WriteHeader(int(status))
//	_ = json.NewEncoder(w).Encode(map[string]interface{}{"error": processedError.Error()})
//}

func (res *ResponseFactory) InternalServerError(c *gin.Context, err error) {
	res.Logger.Error("internal server error", err, nil)
	userError := F{"error": "internal server error"}
	c.AbortWithStatusJSON(http.StatusInternalServerError, userError)
}
