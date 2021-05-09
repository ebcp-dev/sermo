package app

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/ebcp-dev/gorest-api/model"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type App struct {
	Router *mux.Router
	DB     *sql.DB
}

// Receives database credentials and connects to database.
func (a *App) Initialize(user, password, dbname string) {
	connectionString :=
		fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", user, password, dbname)

	var err error
	a.DB, err = sql.Open("postgres", connectionString)
	// Log errors.
	if err != nil {
		log.Fatal(err)
	}

	a.Router = mux.NewRouter()
	a.initializeRoutes()
}

// Starts the application.
func (a *App) Run(addr string) {
	log.Printf("Server listening on port: %s", addr)
	log.Fatal(http.ListenAndServe(":8010", a.Router))
}

// Defines routes.
func (a *App) initializeRoutes() {
	a.Router.HandleFunc("/products", a.getProducts).Methods("GET")
	a.Router.HandleFunc("/product", a.createProduct).Methods("POST")
	a.Router.HandleFunc("/product/{id:[0-9]+}", a.getProduct).Methods("GET")
	a.Router.HandleFunc("/product/{id:[0-9]+}", a.updateProduct).Methods("PUT")
	a.Router.HandleFunc("/product/{id:[0-9]+}", a.deleteProduct).Methods("DELETE")
}

// Route handlers

// Retrieves product from db using id from URL.
func (a *App) getProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	// Convert id string variable to int.
	id, err := strconv.Atoi(vars["id"])
	// Respond with error if id is wrong format/type.
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid product ID")
		return
	}

	p := model.Product{ID: id}
	if err := p.GetProduct(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			// Respond with 404 if product not found in db.
			respondWithError(w, http.StatusNotFound, "Product not found")
		default:
			// Respond if internal server error.
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	// If product found respond product object.
	respondWithJSON(w, http.StatusOK, p)
}

// Gets list of product with count and start variables from URL.
func (a *App) getProducts(w http.ResponseWriter, r *http.Request) {
	// Convert count and start string variables to int.
	count, _ := strconv.Atoi(r.FormValue("count"))
	start, _ := strconv.Atoi(r.FormValue("start"))

	// Default and limit of count is 10.
	if count > 10 || count < 1 {
		count = 10
	}
	// Min start is 0;
	if start < 0 {
		start = 0
	}

	products, err := model.GetProducts(a.DB, start, count)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, products)
}

// Inserts new product into db.
func (a *App) createProduct(w http.ResponseWriter, r *http.Request) {
	var p model.Product
	// Gets JSON object from request body.
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&p); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	defer r.Body.Close()

	if err := p.CreateProduct(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	// Respond with newly created product.
	respondWithJSON(w, http.StatusCreated, p)
}

// Updates product in db using id from URL.
func (a *App) updateProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	// Convert id string variable to int.
	id, err := strconv.Atoi(vars["id"])
	// Respond with error if id is wrong format/type.
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid product ID")
		return
	}

	var p model.Product
	// Gets JSON object from request body.
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&p); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid resquest payload")
		return
	}

	defer r.Body.Close()
	p.ID = id

	if err := p.UpdateProduct(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	// Respond with updated product.
	respondWithJSON(w, http.StatusOK, p)
}

// Deletes product in db using id from URL.
func (a *App) deleteProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	// Convert id string variable to int.
	id, err := strconv.Atoi(vars["id"])
	// Respond with error if id is wrong format/type.
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Product ID")
		return
	}

	p := model.Product{ID: id}
	if err := p.DeleteProduct(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	// Respond with success message if operation is completed.
	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

// Helper functions

// Error message response.
func respondWithError(w http.ResponseWriter, code int, message string) {
	log.Printf("Error: %s", message)
	respondWithJSON(w, code, map[string]string{"error": message})
}

// JSON http response.
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
