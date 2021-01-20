package main

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi"
)

//-------------- API ENDPOINT ------------------//

// AllData returns todoslos datos en la DB
func AllData(w http.ResponseWriter, r *http.Request) {
	dg, cancel := getDgraphClient()
	defer cancel()

	ctx := context.Background()

	txn := dg.NewTxn()
	defer txn.Discard(ctx)
	q := `
	  {
		clientes(func: has(name)) {
		  name
		  age
		  uid
		  followers {
			uid
			name
		  }
		}
	  } 
	  `
	//clientes := []Cliente{}
	//res, err := txn.QueryWithVars(ctx, q, map[string]string{"$a": "Alice"})
	res, err := txn.Query(ctx, q)

	//s := string(`{"operation": "get", "key": "example"}`)

	if err == nil {
		fmt.Printf("%s\n", res.Json)
		w.Header().Set("Content-Type", "application/json")
		w.Write(res.Json)
		//respondwithJSON(w, http.StatusOK, res.Json)
	} else {
		respondwithJSON(w, http.StatusOK, `{"result":"error"}`)
	}
}

// PostCliente crea un nuevo cliente
func dataFromInternet(w http.ResponseWriter, r *http.Request) {
	// fecha: 1/1/2021  = unix = 1611153894
	fecha := chi.URLParam(r, "date") //conseguimos el ID pasado por URL
	fechaInt := 0
	fmt.Printf("%s", fecha)
	if fecha == "" {
		t := time.Now()
		t1 := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, t.Nanosecond(), t.Location()) //normalizamos a dia exacto
		s := t1.Unix()
		fecha = strconv.FormatInt(s, 10)
		fechaInt = int(s)
	}
	fmt.Println("La fecha a buscar es =>", fecha)
	oganizarDB() //establecemos configuraciones en la bd para que funcionen los datos
	// http://localhost:3000/importx/1611153894
	// 18 - 1610946000
	// 19 - 1611032400
	// 20 - 1611118800
	// 21 - 1611205200
	httpReques("https://kqxty15mpg.execute-api.us-east-1.amazonaws.com/buyers?date="+fecha, func(data1 []byte) {
		clientes := string(data1)
		analisisClientesJson(clientes, fechaInt)
		httpReques("https://kqxty15mpg.execute-api.us-east-1.amazonaws.com/products?date="+fecha, func(data2 []byte) {
			productos := string(data2)
			analisisProductosComillas(productos, fechaInt, true)
			httpReques("https://kqxty15mpg.execute-api.us-east-1.amazonaws.com/transactions?date="+fecha, func(data3 []byte) {
				transacciones := string(data3)
				analisisTransaccionesX(transacciones, fechaInt, true)
			})
		})
	})

	/**/
}
