// package main

// import (
//     "log"
//     "net/http"
//     "os"

//     "urlshortener/internal/handlers"
//     "urlshortener/internal/service"
//     "urlshortener/internal/storage"

//     "github.com/joho/godotenv"
// )

// func main() {
//     // Load environment variables from .env file
//     godotenv.Load(".env")

//     // Kết nối Redis
//     store := storage.NewRedisStore()

//     // Tạo services
//     shortenerService := service.NewShortenerService(store)
//     userService := service.NewUserService(store)  // Fix: Tạo đúng instance

//     // Tạo handlers
//     handler := handlers.NewHandler(shortenerService)
//     authHandler := handlers.NewAuthHandler(userService, os.Getenv("JWT_SECRET_KEY"))  // Fix: Dùng userService đúng

//     // Routes
//     http.HandleFunc("/signup", authHandler.SignupHandler)
//     http.HandleFunc("/login", authHandler.LoginHandler)
//     http.HandleFunc("/shorten", authHandler.AuthMiddleware(handler.ShortenHandler))  // Fix: Thêm middleware để require auth
//     http.HandleFunc("/", handler.RedirectHandler)  // Không cần auth

//     log.Println("🚀 URL Shortener đang chạy tại http://localhost:8080")
//     log.Println("Ví dụ: http://localhost:8080/shorten?url=https://google.com")

//     log.Fatal(http.ListenAndServe(":8080", nil))
// }

package main

import (
    "fmt"
    "log"
    "net/http"
    "os"

    "urlshortener/internal/handlers"
    "urlshortener/internal/service"
    "urlshortener/internal/storage"

    "github.com/joho/godotenv"
)

// Struct để lưu route info
type Route struct {
    Method      string
    Path        string
    Description string
}

func main() {
    // Load environment variables from .env file
    godotenv.Load(".env")

    // Kết nối Redis
    store := storage.NewRedisStore()

    // Tạo services
    shortenerService := service.NewShortenerService(store)
    userService := service.NewUserService(store)

    // Tạo handlers
    handler := handlers.NewHandler(shortenerService)
    authHandler := handlers.NewAuthHandler(userService, os.Getenv("JWT_SECRET_KEY"))

    // Slice lưu routes
    var routes []Route

    // Define routes và lưu vào slice
    routes = append(routes, Route{"POST", "/signup", "User signup"})
    http.HandleFunc("/signup", authHandler.SignupHandler)

    routes = append(routes, Route{"POST", "/login", "User login"})
    http.HandleFunc("/login", authHandler.LoginHandler)

    routes = append(routes, Route{"POST", "/shorten", "Create short URL (require auth)"})
    http.HandleFunc("/shorten", authHandler.AuthMiddleware(handler.ShortenHandler))

    routes = append(routes, Route{"GET", "/", "Redirect short URL"})
    http.HandleFunc("/", handler.RedirectHandler)

    // In routes ra terminal (giống Laravel route:list)
    fmt.Println("📋 Registered Routes:")
    fmt.Printf("%-8s %-20s %s\n", "Method", "Path", "Description")
    fmt.Println("---------------------------------------------")
    for _, route := range routes {
        fmt.Printf("%-8s %-20s %s\n", route.Method, route.Path, route.Description)
    }
    fmt.Println()

    log.Println("🚀 URL Shortener đang chạy tại http://localhost:8080")
    log.Println("Ví dụ: http://localhost:8080/shorten?url=https://google.com")

    log.Fatal(http.ListenAndServe(":8080", nil))
}