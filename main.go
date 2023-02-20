package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"github.com/gorilla/mux"
)

type valCalculadora struct {
	Val1      float32 `json:"Val1"`
	Operador  string  `json:"Operador"`
	Val2      float32 `json:"Val2"`
	Resultado float32 `json:"Resultado"`
}

type CalculadoraBD struct {
	Val1      float32
	Operador  string
	Val2      float32
	Resultado float32
}

var calcu = []CalculadoraBD{}

var valores valCalculadora

var scripts = ""

func main() {
	request()
}



func createOperacion(w http.ResponseWriter, req *http.Request) {
	reqBody, err := ioutil.ReadAll(req.Body)
	if err != nil {
		fmt.Fprintf(w, "Insertar operacion valida")
	}

	json.Unmarshal(reqBody, &valores)
	var result float32 = 0

	switch valores.Operador {
	case "+":
		result = valores.Val1 + valores.Val2
	case "-":
		result = valores.Val1 - valores.Val2
	case "*":
		result = valores.Val1 * valores.Val2
	case "/":
		if valores.Val2 == 0 {
			result = 0
		} else {
			result = valores.Val1 / valores.Val2
		}
	}

	respuesta := valCalculadora{
		Val1:      valores.Val1,
		Operador:  valores.Operador,
		Val2:      valores.Val2,
		Resultado: result,
	}

	datosJson, err := json.Marshal(respuesta)
	if err != nil {
		fmt.Println(err)
	}

	json.Unmarshal(datosJson, &valores)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(valores)
}

func enableCORS(router *mux.Router) {
	router.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
	}).Methods(http.MethodOptions)
	router.Use(middlewareCors)
}

func middlewareCors(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, req *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
			w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
			next.ServeHTTP(w, req)
		})
}

func request() {
	router := mux.NewRouter().StrictSlash(false)
	enableCORS(router)
	
	router.HandleFunc("/operacion", createOperacion).Methods("POST")

	log.Println("Escuchando en http://localhost:8000")
	log.Fatal(http.ListenAndServe(":8000", router))
}



