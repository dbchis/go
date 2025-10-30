// package main

// import (
// 	"fmt"
// 	"sync"
// )

// var counter int
// var mu sync.Mutex // Khai báo khóa
// var wg sync.WaitGroup

// func increment() {
// 	// 1. Khóa lại trước khi vào "vùng nguy hiểm"
// 	mu.Lock()

// 	// 2. Dùng defer để đảm bảo LUÔN MỞ KHÓA
// 	// (Kể cả khi hàm panic, nó vẫn được gọi)
// 	defer mu.Unlock()

// 	counter++
// }

// func main() {
// 	for i := 0; i < 1000; i++ {
// 		wg.Add(1)
// 		go func() {
// 			defer wg.Done()
// 			increment()
// 		}()
// 	}

// 	wg.Wait()
// 	fmt.Println("Kết quả cuối cùng:", counter) // LUÔN LUÔN LÀ 1000
// }
