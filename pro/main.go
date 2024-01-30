package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// User struct represents a user in the system
type User struct {
    Username string
    Password string
}

var client *mongo.Client

func main() {
    // Connect to MongoDB
    clientOptions := options.Client().ApplyURI("mongodb+srv://thashmigaacs20:123@cluster0.3gg8g3b.mongodb.net")
    var err error
    client, err = mongo.Connect(context.Background(), clientOptions)
    if err != nil {
        fmt.Println(err)
        return
    }
    defer client.Disconnect(context.Background())

    // Test the connection
    err = client.Ping(context.Background(), nil)
    if err != nil {
        fmt.Println(err)
        return
    }

    http.HandleFunc("/", LoginPage)
    http.HandleFunc("/login", LoginPage)
    http.HandleFunc("/welcome", WelcomePage)
    http.HandleFunc("/signup", SignupPage)

    // Serve static files from the "static" directory.
    http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

    // Start the server on port 8080.
    fmt.Println("Server started on http://localhost:8080")
    http.ListenAndServe(":8080", nil)
}

// SignupPage is the handler for the signup page.
func SignupPage(w http.ResponseWriter, r *http.Request) {
	log.Println("sign")
    if r.Method == http.MethodPost {
        // Retrieve signup form data.
        username := r.FormValue("username")
        password := r.FormValue("password")

        // Perform signup logic here (e.g., store user data in MongoDB).
        if insertUser(username, password) {
            fmt.Printf("New user signup: Username - %s, Password - %s\n", username, password)
            http.Redirect(w, r, "/welcome", http.StatusSeeOther)
            return
        } else {
            http.Error(w, "Error creating user", http.StatusInternalServerError)
            return
        }
    }

    // If not a POST request, serve the signup page template.
    tmpl, err := template.ParseFiles("templates/signup.html")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    tmpl.Execute(w, nil)
}

// insertUser inserts a new user into MongoDB
func insertUser(username, password string) bool {
	log.Println("insert")
    collection := client.Database("Login").Collection("users")

    // Check if the user already exists
    existingUser := User{}
    filter := map[string]interface{}{
        "username": username,
    }

    err := collection.FindOne(context.Background(), filter).Decode(&existingUser)
    if err == nil {
        // User already exists
        fmt.Println("User already exists")
        return false
    }

    // If the user does not exist, insert the new user
    newUser := User{
        Username: username,
        Password: password,
    }

    _, err = collection.InsertOne(context.Background(), newUser)
    if err != nil {
        fmt.Println(err)
        return false
    }

    return true
}

// LoginPage is the handler for the login page.
func LoginPage(w http.ResponseWriter, r *http.Request) {
    if r.Method == http.MethodPost {
        username := r.FormValue("username")
        password := r.FormValue("password")

        // Perform authentication logic here (e.g., check against MongoDB).
        if authenticateUser(username, password) {
            // Successful login, redirect to a welcome page.
            http.Redirect(w, r, "/welcome", http.StatusSeeOther)
            return
        }

        // Invalid credentials, show the login page with an error message.
        fmt.Fprintf(w, "Invalid credentials. Please try again.")
        return
    }

    // If not a POST request, serve the login page template.
    tmpl, err := template.ParseFiles("templates/login.html")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    tmpl.Execute(w, nil)
}

// authenticateUser checks if the provided username and password match a user in MongoDB
func authenticateUser(username, password string) bool {
    collection := client.Database("Login").Collection("users")

    var user User
    filter := map[string]interface{}{
        "username": username,
        "password": password,
    }

    err := collection.FindOne(context.Background(), filter).Decode(&user)
    if err != nil {
        fmt.Println(err)
        return false
    }

    return true
}

// WelcomePage is the handler for the welcome page.
func WelcomePage(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Welcome, you have successfully logged in!")
}

// import (
// 	"context"
// 	"fmt"
// 	"log"

// 	"go.mongodb.org/mongo-driver/mongo"
// 	"go.mongodb.org/mongo-driver/mongo/options"
// )

// func main() {
// 	// Set the MongoDB Atlas connection string
// 	connectionString := "mongodb+srv://thashmigaacs20:123@cluster0.3gg8g3b.mongodb.net"
// 	// "mongodb+srv://userone:L2ORM855sl5XzunU@cluster0.9eqsnyg.mongodb.net/"
// 	// Set up client options
// 	clientOptions := options.Client().ApplyURI(connectionString)

// 	// Connect to MongoDB
// 	client, err := mongo.Connect(context.TODO(), clientOptions)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	// Ping the MongoDB server to verify the connection
// 	err = client.Ping(context.TODO(), nil)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	fmt.Println("Connected to MongoDB!")

// 	// Do your MongoDB operations here...

// 	// Disconnect from MongoDB
// 	err = client.Disconnect(context.TODO())
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	fmt.Println("Connection to MongoDB closed.")
// }
