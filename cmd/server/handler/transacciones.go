package handler

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/SebasPu/Transacciones/internal/transacciones"
	"github.com/SebasPu/Transacciones/pkg/web"
	"github.com/gin-gonic/gin"
	"github.com/ignaciofalco/go-dynamic-getter/dynamicgetter"
)

type request struct {
	CodigoTransaccion string `json:"codigo" `
	Moneda            string `json:"moneda" `
	Monto             string `json:"monto" `
	Emisor            string `json:"emisor" `
	Receptor          string `json:"receptor" `
	Fecha             string `json:"fecha" `
}

type Transaccion struct {
	service transacciones.Service
}

const UNPROCESSABLE_ENTITY = 422
const NOT_FOUND = 404

func validateRequiredField(req request, requireFields []string) error {
	for _, field := range requireFields {
		_, err := dynamicgetter.GetField(&req, field, false)
		if err != nil {
			return err
		}
	}
	return nil
}

// GetAll godoc
// @Summary List transacciones
// @Tags Transacciones
// @Description get transacciones
// @Accept  json
// @Produce  json
// @Param token header string true "token"
// @Router /transacciones/getAll [get]
func (t *Transaccion) GetAll() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		p, err := t.service.GetAll()
		if err != nil {
			ctx.JSON(500, web.NewResponse(500, nil, fmt.Errorf("mensaje de error")))
			return
		}
		ctx.JSON(200, web.NewResponse(200, p, nil))
	}
}

func (t *Transaccion) LastId() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		p, err := t.service.LastId()
		if err != nil {
			ctx.JSON(404, gin.H{
				"error": err.Error(),
			})
			return
		}
		ctx.JSON(200, p)
	}
}

// Store godoc
// @Summary Store transacciones
// @Tags Transacciones
// @Description store transacciones
// @Accept  json
// @Produce  json
// @Param token header string true "token"
// @Param transaction body request true "Transaction to store"
// @Router /transacciones/store [post]
func (t *Transaccion) Store() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.GetHeader("token")
		if token != os.Getenv("TOKEN") {
			ctx.JSON(401, gin.H{"error": "unauthorized"})
			return
		}
		var r request
		if err := ctx.ShouldBindJSON(&r); err != nil {
			ctx.JSON(404, gin.H{
				"error": err.Error(),
			})
			return
		}
		id, err := t.service.Store(
			r.CodigoTransaccion,
			r.Moneda,
			r.Monto,
			r.Emisor,
			r.Receptor,
			r.Fecha)
		if err != nil {
			ctx.JSON(500, gin.H{
				"error": err.Error(),
			})
			return
		}
		ctx.JSON(200, id)
	}
}

func (t *Transaccion) Update() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.GetHeader("token")
		if token != os.Getenv("TOKEN") {
			ctx.JSON(401, gin.H{"error": "unauthorized"})
			return
		}
		id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
		if err != nil {
			ctx.JSON(400, gin.H{"error": "invalid ID"})
			return
		}
		var req request
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(400, gin.H{"error": err.Error()})
			return
		}
		requiredFields := []string{"codigo", "monto", "emisor", "receptor", "fecha"}
		err = validateRequiredField(req, requiredFields)
		if err != nil {
			switch {
			case errors.Is(err, dynamicgetter.ErrZeroValue):
				ctx.AbortWithStatusJSON(422, web.NewResponse(UNPROCESSABLE_ENTITY, nil, fmt.Errorf("todos los campos son requeridos")))
				return
			}
		}
		p, err := t.service.Update(
			int(id),
			req.CodigoTransaccion,
			req.Moneda,
			req.Monto,
			req.Emisor,
			req.Receptor,
			req.Fecha)
		if err != nil {
			ctx.JSON(400, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(200, p)
	}
}

func (t *Transaccion) UpdateCod() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.GetHeader("token")
		if token != os.Getenv("TOKEN") {
			ctx.JSON(401, gin.H{"error": "unauthorized"})
			return
		}
		id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
		if err != nil {
			ctx.JSON(400, gin.H{"error": "invalid ID"})
			return
		}
		var req request
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(400, gin.H{"error": err.Error()})
			return
		}
		requiredFields := []string{"codigo", "monto"}
		err = validateRequiredField(req, requiredFields)
		if err != nil {
			switch {
			case errors.Is(err, dynamicgetter.ErrZeroValue):
				ctx.AbortWithStatusJSON(422, web.NewResponse(UNPROCESSABLE_ENTITY, nil, fmt.Errorf("todos los campos son requeridos")))
				return
			}
		}
		p, err := t.service.UpdateCod(int(id), req.CodigoTransaccion, req.Monto)
		if err != nil {
			ctx.JSON(400, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(200, p)
	}
}

func (t *Transaccion) Delete() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.GetHeader("token")
		if token != os.Getenv("TOKEN") {
			ctx.JSON(401, gin.H{"error": "unauthorized"})
			return
		}
		id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
		if err != nil {
			ctx.JSON(400, gin.H{"error": "invalid ID"})
			return
		}
		err = t.service.Delete(int(id))
		if err != nil {
			ctx.JSON(400, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(200, gin.H{"data": fmt.Sprintf("El producto %d ha sido eliminado", id)})
	}
}

func NewTransaccion(s transacciones.Service) *Transaccion {
	return &Transaccion{service: s}
}

//go install github.com/swaggo/swag/cmd/swag@latest
