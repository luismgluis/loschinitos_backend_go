package main

import (
	"context"
	"fmt"
	"net/http"
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
	oganizarDB() //establecemos configuraciones en la bd para que funcionen los datos
	httpReques("https://kqxty15mpg.execute-api.us-east-1.amazonaws.com/transactions", func(data3 []byte) {
		transacciones := string(data3)
		analisisTransaccionesX(transacciones)
	})

	/*httpReques("https://kqxty15mpg.execute-api.us-east-1.amazonaws.com/buyers", func(data1 []byte) {
		clientes := string(data1)
		analisisClientesJson(clientes)
		httpReques("https://kqxty15mpg.execute-api.us-east-1.amazonaws.com/products", func(data2 []byte) {
			productos := string(data2)
			analisisProductosComillas(productos, true)
			httpReques("https://kqxty15mpg.execute-api.us-east-1.amazonaws.com/transactions", func(data3 []byte) {
				transacciones := string(data3)
				analisisTransaccionesX(transacciones)
			})
		})
	})*/

	/*cliente := new(Cliente)
	//esto nos ayuda a asignar los datos enviados desde frontend a el struct
	if errs := binding.Bind(r, cliente); errs != nil {
		http.Error(w, errs.Error(), http.StatusBadRequest)
		return
	}
	query := fmt.Sprintf(`
	{
		cliente(func: uid(%s)){
			uid
			name
			age
		}
	}
	`, cliente.UID)

	ConsultaDataBase(query, func(data []byte) {
		if data != nil { //el cliente existe entonces decimos rechazamos la creacion
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"result":"error"}`))
		} else { //no existe entonces si lo creamos
			cliente.UID = "_:elid" //asignamos esto para que se ponga el ID automaticamente
			jsonbytes, err := json.Marshal(cliente)
			if err != nil {
				log.Fatal(err)
			}
			MutacionDataBase(jsonbytes, func(data []byte) {
				w.Header().Set("Content-Type", "application/json")
				if data == nil {
					w.Write([]byte(`{"result":"ok"}`))
				} else {
					w.Write([]byte(`{"result":"error"}`))
				}
			})
		}
	})*/

}
