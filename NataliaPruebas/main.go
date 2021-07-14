package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type infoVenta struct {
	Id        int64
	Nombre    string
	Contenido int64
	PrecioU   float32
}

type producto struct {
	Id        int64
	Marca     string
	Cantidad  int
	Nombre    string
	Contenido int64
	PrecioU   float32
	IdPlus    int64
}

var productos = []producto{
	{1, "Margarita", 100, "Papas Limón", 39, 1300, 101},
	{2, "Fritolay", 150, "Maní Salado", 50, 800, 102},
	{3, "Dux", 100, "Galletas Rellenas", 36, 1000, 103},
	{4, "Kythos", 80, "Yucas Fritas", 25, 1200, 104},
}

func main() {
	r := gin.Default()

	r.GET("/productos", productoGet)
	r.GET("/productos/:id", productoGetId)
	r.POST("/productos", productoPost)
	r.PUT("/productos/:id", productoPut)
	r.DELETE("productos/:id", productDelete)
	r.Run()
}

func productoGet(c *gin.Context) {
	contstr := c.DefaultQuery("contenido", "0")
	preciostr := c.DefaultQuery("precio", "0")

	if contstr == "" {
		contstr = "0"
	}
	if preciostr == "" {
		preciostr = "0"
	}

	cont, errCont := strconv.ParseInt(contstr, 10, 64)
	precio, errPre := strconv.ParseFloat(preciostr, 32)

	var errs = []string{}
	if errCont != nil {
		errs = append(errs, "El contenido debe ser un número")
	}
	if errPre != nil {
		errs = append(errs, "El Precio debe ser un número")
	}
	if precio < 0 {
		errs = append(errs, "El valor de precio deben ser positivo")
	}
	if cont < 0 {
		errs = append(errs, "El valor de contenido deben ser positivo")
	}

	if len(errs) > 0 {
		c.JSON(400, gin.H{
			"Errors": errs,
		})
		return
	}

	var arrInfo = []infoVenta{}
	for _, data := range productos {
		if data.Contenido >= cont && data.PrecioU >= float32(precio) {
			arrInfo = append(arrInfo, infoVenta{data.Id, data.Nombre, data.Contenido, data.PrecioU})
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"values": arrInfo,
	})

}

func productoGetId(c *gin.Context) {
	idstr := c.Param("id")
	id, err := strconv.ParseInt(idstr, 10, 64)

	if err != nil {
		c.JSON(400, gin.H{
			"Error": "El id debe ser un número",
		})
		return
	}

	for _, p := range productos {
		if id == p.Id {
			c.JSON(200, gin.H{
				"value": p,
			})
			return
		}
	}

	c.JSON(404, gin.H{})
}

func productoPost(c *gin.Context) {
	var p producto
	if err := c.ShouldBindJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	pp := productos[len(productos)-1]
	newId := pp.Id + 1
	newIdPlus := pp.Id + 1 + 100
	fmt.Println(newIdPlus)

	var errs = []string{}

	if p.Marca == "" {
		errs = append(errs, "Debe enviarse una Marca del producto")
	}
	if p.Cantidad == 0 {
		errs = append(errs, "Debe enviarse una Cantidad del producto")
	}
	if p.Nombre == "" {
		errs = append(errs, "Debe enviarse un Nombre del producto")
	}
	if p.Contenido == 0 {
		errs = append(errs, "Debe enviarse un Contenido del producto")
	}
	if p.PrecioU == 0 {
		errs = append(errs, "Debe enviarse un PrecioU del producto")
	}
	if p.Cantidad < 0 {
		errs = append(errs, "La cantidad debe ser un valor positivo mayor a 0")
	}
	if p.Contenido < 0 {
		errs = append(errs, "El contenido debe ser un valor positivo mayor a 0")
	}
	if p.PrecioU < 0 {
		errs = append(errs, "El precioU debe ser un valor positivo mayor a 0")
	}

	if len(errs) > 0 {
		c.JSON(400, gin.H{
			"Errors": errs,
		})
		return
	}

	productos = append(productos, producto{
		newId, p.Marca, p.Cantidad,
		p.Nombre, p.Contenido, p.PrecioU, newIdPlus},
	)

	c.JSON(http.StatusCreated, gin.H{
		"Id":        newId,
		"Marca":     p.Marca,
		"Cantidad":  p.Cantidad,
		"Nombre":    p.Nombre,
		"Contenido": p.Contenido,
		"PrecioU":   p.PrecioU,
		"IdPlus":    newIdPlus,
	})

}

func productoPut(c *gin.Context) {
	var newP producto
	idstr := c.Param("id")

	id, err := strconv.ParseInt(idstr, 10, 64)

	if err != nil {
		c.JSON(400, gin.H{
			"Error": "El id debe ser un número",
		})
		return
	}

	if err := c.ShouldBindJSON(&newP); err != nil { //Sintaxis if (variable creacion); validacion
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	for i, p := range productos {
		if id == p.Id {
			newP.Id = p.Id
			newP.IdPlus = p.IdPlus
			productos[i] = newP
			c.JSON(http.StatusOK, gin.H{
				"Id":        p.Id,
				"Marca":     newP.Marca,
				"Cantidad":  newP.Cantidad,
				"Nombre":    newP.Nombre,
				"Contenido": newP.Contenido,
				"PrecioU":   newP.PrecioU,
				"IdPlus":    p.IdPlus,
			})
			return
		}
	}

	c.JSON(404, gin.H{
		"error": "El id no existe",
	})
	return
}

func productDelete(c *gin.Context) {
	idstr := c.Param("id")
	id, err := strconv.ParseInt(idstr, 10, 64)
	if err != nil {
		c.JSON(400, gin.H{
			"Error": "El id debe ser un número",
		})
		return
	}

	var arrtemp []producto

	pos := -1
	for i, p := range productos {
		if id != p.Id {
			arrtemp = append(arrtemp, p)
		} else {
			pos = i
		}
	}
	if pos < 0 {
		c.JSON(404, gin.H{
			"error": "El id no existe",
		})
		return
	}

	productos = arrtemp

	c.Status(http.StatusOK)
}
