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
