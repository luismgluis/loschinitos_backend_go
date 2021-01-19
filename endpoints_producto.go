package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi"

	"github.com/mholt/binding"
)

//-------------- API ENDPOINT ------------------//

func GetProductoIDfromOldID(olduid string, fn FunctionBackString) {
	q := fmt.Sprintf(`
	{
		productos(func: eq(PRODID,"%s")) {
		  uid
		  PRODID
		}
	}
	`, olduid)

	ConsultaDataBase(q, func(data []byte) {
		ccc := Productos{}
		err3312 := json.Unmarshal(data, &ccc)
		if err3312 == nil {
			prods := ccc.Productos
			micliente := prods[0]
			if micliente.UID != "" {
				fn(micliente.UID)
			} else {
				fn("")
			}
		} else {
			fn("")
		}
	})
}

// AllClientes returns todos los clientes en la DB
func AllProducto(w http.ResponseWriter, r *http.Request) {
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

func GetProductoByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id") //conseguimos el ID pasado por URL
	fmt.Printf("%s", id)
	query := fmt.Sprintf(`
	{
		cliente(func: uid(%s)){
			uid
			name
			age
		}
	}
	`, id)
	ConsultaDataBaseJson(query, func(data string) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(data))
	})
}

// PostCliente crea un nuevo cliente
func PostProducto(w http.ResponseWriter, r *http.Request) {
	cliente := new(Cliente)
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
					//fmt.Printf("%s", string(data))
					w.Write([]byte(`{"result":"ok"}`))
				} else {
					w.Write([]byte(`{"result":"error"}`))
				}
			})
		}
	})
}

// PutCliente actualiza un cliente
func PutProducto(w http.ResponseWriter, r *http.Request) {

	cliente := new(Cliente)
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
		if data != nil { //el cliente existe entonces lo actualizamos
			jsonbytes, err := json.Marshal(cliente)
			if err != nil {
				log.Fatal(err)
			}
			MutacionDataBase(jsonbytes, func(data []byte) {
				w.Header().Set("Content-Type", "application/json")
				if data == nil {
					//fmt.Printf("%s", string(data))
					w.Write([]byte(`{"result":"ok"}`))
				} else {
					w.Write([]byte(`{"result":"error"}`))
				}
			})
		} else { //no existe entonces si lo rechazamos
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"result":"error"}`))
		}
	})
}

// DeleteCliente remove a spesific post
func DeleteProducto(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id") //conseguimos el ID pasado por URL
	fmt.Printf("%s", id)
	query := fmt.Sprintf(`
	{
		"delete": [
			{
				"uid": "%s"
			}
		]
	}
	`, id)
	ConsultaDataBase(query, func(data []byte) {
		w.Header().Set("Content-Type", "application/json")
		if data != nil {
			//ok
			w.Write([]byte(`{"result":"ok"}`))
		} else {
			//bad
			w.Write([]byte(`{"result":"error"}`))
		}
	})
}
