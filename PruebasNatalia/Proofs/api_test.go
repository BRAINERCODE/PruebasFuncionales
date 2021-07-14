package proofs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_productoGet(t *testing.T) {

	type getProducts struct {
		Contenido int    `json:"Contenido"`
		ID        int    `json:"Id"`
		Nombre    string `json:"Nombre"`
		PrecioU   int    `json:"PrecioU"`
	}
	type responseRequest struct {
		Values []getProducts `json:"values"`
	}

	getAllProductsGET := responseRequest{

		Values: []getProducts{
			{
				Contenido: 39,
				ID:        1,
				Nombre:    "Papas Limón",
				PrecioU:   1300,
			},
			{
				Contenido: 50,
				ID:        2,
				Nombre:    "Maní Salado",
				PrecioU:   800,
			},
			{
				Contenido: 36,
				ID:        3,
				Nombre:    "Galletas Rellenas",
				PrecioU:   1000,
			},
			{
				Contenido: 25,
				ID:        4,
				Nombre:    "Yucas Fritas",
				PrecioU:   1200,
			},
		},
	}
	response, _ := http.Get("http://localhost:8080/productos")

	ResponseProducts, _ := ioutil.ReadAll(response.Body)
	var StructGetProducts responseRequest

	err := json.Unmarshal(ResponseProducts, &StructGetProducts)

	assert.Len(t, StructGetProducts.Values, 4)
	assert.NoError(t, err)
	assert.Equal(t, 200, response.StatusCode)
	assert.Equal(t, StructGetProducts, getAllProductsGET)
}

func Test_productoGetId(t *testing.T) {
	url := "http://restapi3.apiary.io/notes"
	fmt.Println("URL:>", url)
	var jsonStr = []byte(`{"title":"Buy cheese and bread for breakfast."}`)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
}
