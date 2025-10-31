// package main

// import (
// 	"fmt"
// 	"sync"
// )

// var config = make(map[string]string)
// var mu sync.RWMutex
// var wg sync.WaitGroup

// // Hàng trăm Goroutine có thể gọi hàm này cùng lúc
// func readConfig(key string) string {
// 	mu.RLock() // Lấy khóa ĐỌC
// 	defer mu.RUnlock()
// 	return config[key]
// }

// // Chỉ một Goroutine gọi hàm này
// func writeConfig(key, value string) {
// 	mu.Lock() // Lấy khóa GHI (độc quyền)
// 	defer mu.Unlock()
// 	config[key] = value
// }

// func main() {
// 	// 1 Goroutine ghi
// 	wg.Add(1)
// 	go func() {
// 		defer wg.Done()
// 		writeConfig("server_url", "http://prod.com")
// 	}()

// 	// 5 Goroutines đọc
// 	for i := 0; i < 5; i++ {
// 		wg.Add(1)
// 		go func(id int) {
// 			defer wg.Done()
// 			fmt.Printf("Reader %d đọc được: %s\n", id, readConfig("server_url"))
// 		}(i)
// 	}

// 	wg.Wait()
// }
