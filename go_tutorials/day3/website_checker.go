package main

import (
	"fmt"
	"net/http"
	"sync"
)

// Worker chỉ làm việc, không trả kết quả về channel
func worker(wg *sync.WaitGroup, id int, url string) {
    // 3. Báo cáo "Done" khi hàm kết thúc
    defer wg.Done()

	fmt.Printf("Worker %d đang xử lý: %s\n", id, url)
	if _, err := http.Get(url); err != nil {
		fmt.Printf("[LỖI] Worker %d, %s: %s\n", id, url, err)
	} else {
		fmt.Printf("[OK] Worker %d: %s\n", id, url)
	}
}

func main() {
	urls := []string{
		"http://google.com",
		"http://facebook.com",
		"http://github.com",
		"http://nonexistent-domain.dev",
	}

    // 1. Khởi tạo WaitGroup
    var wg sync.WaitGroup

	for i, url := range urls {
        // 2. Thêm 1 công việc vào bộ đếm
        wg.Add(1)
		go worker(&wg, i+1, url)
	}

    // 4. Chờ (BLOCK) cho đến khi bộ đếm về 0
    wg.Wait()

    fmt.Println("Tất cả website đã được kiểm tra xong!")
    // KHÔNG CẦN time.Sleep() nữa!
}