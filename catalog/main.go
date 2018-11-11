package main

import (
	"log"
	"os"
    "fmt"
    "errors"
    "strconv"
	"time"
    "path/filepath"
    "net/http"
    "encoding/json"
	"context"

	"google.golang.org/grpc"
    "google.golang.org/grpc/credentials"

    pb "microservices-grpc-go-python/catalog/ecommerce"
)

func getDiscountConnection (host string) (*grpc.ClientConn, error){
    wd, _ := os.Getwd()
    parentDir := filepath.Dir(wd)
    certFile := filepath.Join(parentDir, "keys", "cert.pem")
    creds, _ := credentials.NewClientTLSFromFile(certFile, "")
	return grpc.Dial(host, grpc.WithTransportCredentials(creds))
}

func findCustomerByID(id int) (pb.Customer, error) {
    c1 := pb.Customer{Id: 1, FirstName: "John", LastName: "Snow"}
    c2 := pb.Customer{Id: 2, FirstName: "Daenerys", LastName: "Targaryen"}
    customers := map[int]pb.Customer{
        1: c1,
        2: c2,
    }
    found, ok := customers[id]
    if ok {
        return found, nil
    }
    return found, errors.New("Customer not found.")
}

func getFakeProducts() []*pb.Product {
    p1 := pb.Product{Id: 1, Slug: "iphone-x", Description: "64GB, black and iOS 12", PriceInCents: 99999}
    p2 := pb.Product{Id: 2, Slug: "notebook-avell-g1511", Description: "Notebook Gamer Intel Core i7", PriceInCents: 150000}
    p3 := pb.Product{Id: 3, Slug: "playstation-4-slim", Description: "1TB Console", PriceInCents: 32999}
    return []*pb.Product{&p1, &p2, &p3}
}

func getProductsWithDiscountApplied(customer pb.Customer, products []*pb.Product) []*pb.Product {
    host := os.Getenv("DISCOUNT_SERVICE_HOST")
    if len(host) == 0 {
        host = "localhost:11443"
    }
    conn, err := getDiscountConnection(host)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
    defer conn.Close()

    c := pb.NewDiscountClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second * 10)
	defer cancel()

    productsWithDiscountApplied := make([]*pb.Product, 0)
    for _, product := range products {
        r, err := c.ApplyDiscount(ctx, &pb.DiscountRequest{Customer: &customer, Product: product})
        if err == nil {
            productsWithDiscountApplied = append(productsWithDiscountApplied, r.GetProduct())
        } else {
            log.Println("Failed to apply discount.", err)
        }
    }

    if len(productsWithDiscountApplied) > 0 {
        return productsWithDiscountApplied
    }
    return products
}

func handleGetProducts(w http.ResponseWriter, req *http.Request) {
    products := getFakeProducts()
    w.Header().Set("Content-Type", "application/json")

    customerID := req.Header.Get("X-USER-ID")
    if customerID == "" {
        json.NewEncoder(w).Encode(products)
        return
    }
    id, err := strconv.Atoi(customerID)
    if err != nil {
        http.Error(w, "Customer ID is not a number.", http.StatusBadRequest)
        return
    }

    customer, err := findCustomerByID(id)
    if err != nil {
        json.NewEncoder(w).Encode(products)
        return
    }

    productsWithDiscountApplied := getProductsWithDiscountApplied(customer, products)
    json.NewEncoder(w).Encode(productsWithDiscountApplied)
}

func main() {
	port := "11080"
	if len(os.Args) > 1 {
		port = os.Args[1]
	}

    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "It is working.")
    })
    http.HandleFunc("/products", handleGetProducts)

    fmt.Println("Server running on", port)
    http.ListenAndServe(":"+port, nil)
}
