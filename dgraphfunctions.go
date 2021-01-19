package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/dgraph-io/dgo/v200"
	"github.com/dgraph-io/dgo/v200/protos/api"
	"google.golang.org/grpc"
	//ayuda para httpreques
)

type CancelFunc func()
type FunctionBackBytes func(data []byte)
type FunctionBackString func(data string)
type FunctionBackArrStrings func(data []string)

func getDgraphClient() (*dgo.Dgraph, CancelFunc) {

	direccion := host + ":" + strconv.Itoa(port)
	conn, err := grpc.Dial(direccion, grpc.WithInsecure())
	if err != nil {
		log.Fatal("While trying to dial gRPC")
	}
	dc := api.NewDgraphClient(conn)
	dgraphClient := dgo.NewDgraphClient(dc)
	//------------
	/**/
	//graphql  /mutate /commit.
	//https://gutsy-grape.us-west-2.aws.cloud.dgraph.io/query
	// T1LEBKB4N6+iEdgv6oxYTiW9XQntVsTLzVTwjttyDr4=  ADMIN
	// EpUZOdIYZBGkPqo1BwkVfSJ8j9ALBdcrNBkD2PT21xk=  CLIENT

	/*conn, err := dgo.DialSlashEndpoint("https://gutsy-grape.us-west-2.aws.cloud.dgraph.io/query", "T1LEBKB4N6+iEdgv6oxYTiW9XQntVsTLzVTwjttyDr4=") //"EpUZOdIYZBGkPqo1BwkVfSJ8j9ALBdcrNBkD2PT21xk=")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	dgraphClient := dgo.NewDgraphClient(api.NewDgraphClient(conn))
	*/
	if err != nil {
		log.Fatalf("While trying to login %v", err.Error())
	}

	return dgraphClient, func() {
		if err := conn.Close(); err != nil {
			log.Printf("Error while closing connection:%v", err)
		}
	}
}

func Pruebaget_data() {
	dg, cancel := getDgraphClient()
	defer cancel()

	ctx := context.Background()

	txn := dg.NewTxn()
	defer txn.Discard(ctx)

	q := `query all($a: string) {
		all(func: eq(name, $a)) {
		  name
		}
	  }`
	q = `
	query people(func: has(name)) {
		name
		age
		uid
		followers {
		  uid
		  name
		}
	  }`
	q = `
	  {
		people(func: has(name)) {
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
	//res, err := txn.QueryWithVars(ctx, q, map[string]string{"$a": "Alice"})
	res, err := txn.Query(ctx, q)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", res.Json)
}
func Pruebaset_data() {
	dg, cancel := getDgraphClient()
	defer cancel()

	ctx := context.Background()

	op := &api.Operation{
		Schema: `name: string @index(exact) .`,
	}
	err := dg.Alter(ctx, op)

	type Persona struct {
		Uid   string   `json:"uid,omitempty"`
		Name  string   `json:"name,omitempty"`
		DType []string `json:"dgraph.type,omitempty"`
	}

	p := Persona{
		Uid:   "_:alice",
		Name:  "Alice",
		DType: []string{"Persona"},
	}

	pb, err := json.Marshal(p)
	if err != nil {
		log.Fatal(err)
	}

	mu := &api.Mutation{
		SetJson: pb,
	}

	txn := dg.NewTxn()
	defer txn.Discard(ctx)

	res, err := txn.Mutate(ctx, mu)
	if err != nil {
		log.Fatal(err)
	}

	err2 := txn.Commit(context.Background())
	if err2 != nil {
		log.Fatal(err2)
	}
	/*req := &api.Request{CommitNow:true, Mutations: []*api.Mutation{mu}}
	res, err := txn.Do(ctx, req)
	if err != nil {
	  log.Fatal(err)
	}*/
	fmt.Printf("%s\n", res)
}
func Pruebaset_data2() {
	dg, cancel := getDgraphClient()
	defer cancel()

	ctx := context.Background()

	op := &api.Operation{
		Schema: `name: string @index(exact) .`,
	}
	err := dg.Alter(ctx, op)

	type Persona struct {
		Uid   string   `json:"uid,omitempty"`
		Name  string   `json:"name,omitempty"`
		DType []string `json:"dgraph.type,omitempty"`
	}

	p := Persona{
		Uid:   "0x6", //"_:pedro",
		Name:  "Alice One",
		DType: []string{"Persona"},
	}

	pb, err := json.Marshal(p)
	if err != nil {
		log.Fatal(err)
	}

	mu := &api.Mutation{
		SetJson: pb,
	}

	txn := dg.NewTxn()
	defer txn.Discard(ctx)

	res, err := txn.Mutate(ctx, mu)
	if err != nil {
		log.Fatal(err)
	}

	err2 := txn.Commit(context.Background())
	if err2 != nil {
		log.Fatal(err2)
	}
	/*req := &api.Request{CommitNow:true, Mutations: []*api.Mutation{mu}}
	res, err := txn.Do(ctx, req)
	if err != nil {
	  log.Fatal(err)
	}*/
	fmt.Printf("%s\n", res)
}

func OldInsertNewCliente(clientef Cliente) {
	dg, cancel := getDgraphClient()
	defer cancel()

	ctx := context.Background()

	op := &api.Operation{
		Schema: `name: string @index(exact) .`,
	}
	err := dg.Alter(ctx, op)

	/*
		}p := Cliente{
			UID:    "0x1", //"_:pedro",
			Name:   "Alice One1",
			Age:    30,
			Avatar: "",
			DType:  []string{"Persona"},
		}*/

	clientef.UID = "_:elid"
	pb, err := json.Marshal(clientef)
	if err != nil {
		log.Fatal(err)
	}

	mu := &api.Mutation{
		SetJson: pb,
	}

	txn := dg.NewTxn()
	defer txn.Discard(ctx)

	res, err := txn.Mutate(ctx, mu)
	if err != nil {
		log.Fatal(err)
	}

	err2 := txn.Commit(context.Background())
	if err2 != nil {
		log.Fatal(err2)
	}
	/*req := &api.Request{CommitNow:true, Mutations: []*api.Mutation{mu}}
	res, err := txn.Do(ctx, req)
	if err != nil {
	  log.Fatal(err)
	}*/
	fmt.Printf("%s\n", res)
}

func oganizarDB() {
	dg, cancel := getDgraphClient()
	defer cancel()

	ctx := context.Background()

	op := &api.Operation{
		Schema: `
		name: string @index(exact) .
		avatar: string .
		age: int .  
		price: int .
		buyerid: string .
		ip: string .
		device: string .
		produtids: [string] .
		TRANSID: string @index(exact) .
		PRODID: string @index(exact) .
		UIDOLD: string @index(exact) .
		`,
	}
	err_alter := dg.Alter(ctx, op)
	if err_alter != nil {
		log.Fatal(err_alter)
	}
}

func MutacionDataBase(bbb []byte, fn FunctionBackBytes) {

	dg, cancel := getDgraphClient()
	defer cancel()

	ctx := context.Background()

	txn := dg.NewTxn()
	defer txn.Discard(ctx)

	/*bites, err := json.Marshal(&bbb)
	if err != nil {
		log.Fatal(err)
	}*/

	mu := &api.Mutation{
		SetJson: bbb,
	}

	res, err := txn.Mutate(ctx, mu)
	if err != nil {
		log.Fatal(err)
	}
	//inesperado tenia que hacer esto
	err2 := txn.Commit(context.Background())
	if err2 != nil {
		log.Fatal(err2)
	}

	fn(res.GetJson())
}

func OldMutacionDataBase(bbb Cliente, fn FunctionBackString) {

	dg, cancel := getDgraphClient()
	defer cancel()

	ctx := context.Background()

	op := &api.Operation{
		Schema: `name: string @index(exact) .
		avatar: string .
		age: int .`,
	}
	err_alter := dg.Alter(ctx, op)
	if err_alter != nil {
		log.Fatal(err_alter)
	}

	txn := dg.NewTxn()
	defer txn.Discard(ctx)

	bites, err := json.Marshal(&bbb)
	if err != nil {
		log.Fatal(err)
	}

	mu := &api.Mutation{
		SetJson: bites,
	}

	res, err := txn.Mutate(ctx, mu)
	if err != nil {
		log.Fatal(err)
	}
	//inesperado tenia que hacer esto
	err2 := txn.Commit(context.Background())
	if err2 != nil {
		log.Fatal(err2)
	}
	/*
		req := &api.Request{
			CommitNow: true,
			Mutations: []*api.Mutation{mu},
		}
		res, err := txn.Do(ctx, req)
		if err != nil {
			log.Fatal(err)
		}*/
	fmt.Printf("%s", res)
}

func ConsultaDataBase(query string, fn FunctionBackBytes) {
	dg, cancel := getDgraphClient()
	defer cancel()

	ctx := context.Background()

	txn := dg.NewTxn()
	defer txn.Discard(ctx)

	//clientes := []Cliente{}
	//res, err := txn.QueryWithVars(ctx, q, map[string]string{"$a": "Alice"})
	res, err := txn.Query(ctx, query)

	//s := string(`{"operation": "get", "key": "example"}`)

	if err == nil {
		fmt.Printf("%s\n", res.Json)
		fn(res.Json)
	} else {
		fn(nil)
	}
}

func ConsultaDataBaseJson(query string, fn FunctionBackString) {
	ConsultaDataBase(query, func(data []byte) {
		fn("" + string(data))
	})
}
