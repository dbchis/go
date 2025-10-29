package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

// Đây là hàm của worker
// Nó nhận vào 2 channel: jobs để lấy việc và results để trả kết quả.
func worker(id int, jobs <-chan string, results chan<- string) {
	for url := range jobs {
		fmt.Printf("Worker %d đang xử lý: %s\n", id, url)

		// Thêm "http://" nếu cần
		if !strings.HasPrefix(url, "http") {
			url = "http://" + url
		}

		resp, err := http.Get(url)
		if err != nil {
			results <- fmt.Sprintf("[LỖI] %s: %s", url, err)
			continue
		}

		results <- fmt.Sprintf("%s: %d %s", url, resp.StatusCode, resp.Status)
		resp.Body.Close()
	}
}

func run() {
	// Danh sách công việc
	urls := []string{
		"google.com",
		"facebook.com",
		"golang.org",
		"bing.com",
		"nonexistent-domain.dev", // URL này sẽ lỗi
		"github.com",
	}
	numJobs := len(urls)

	// === BẮT ĐẦU THỰC HÀNH TẠI ĐÂY ===

	// Khởi tạo 2 channel
	// TODO: Thử thay đổi giữa buffered và unbuffered
	jobs := make(chan string, numJobs)
	results := make(chan string, numJobs)
	// jobs := make(chan string)
	// results := make(chan string)

	// 1. Khởi tạo 3 worker
	numWorkers := 3
	for w := 1; w <= numWorkers; w++ {
		go worker(w, jobs, results)
	}

	// 2. Gửi công việc (jobs) vào channel `jobs`
	for _, url := range urls {
		jobs <- url
	}
	close(jobs) // Đóng channel `jobs` để báo cho worker biết đã hết việc

	// 3. Nhận kết quả từ channel `results`
	// for a := 1; a <= numJobs; a++ {
	// 	// TODO: Thay thế vòng lặp này bằng `select`
	// 	fmt.Println(<-results)
	// }
	// (Trong hàm main, thay thế Phần 3)

	// 3. Nhận kết quả với `select` và `timeout`
	timeout := time.After(5 * time.Second) // Đặt timeout tổng là 5 giây

	for i := 0; i < numJobs; i++ {
		select {
		case res := <-results:
			// SỰ KIỆN 1: Nhận được kết quả
			fmt.Println(res)

		case <-timeout:
			// SỰ KIỆN 2: Hết 5 giây
			fmt.Println("QUÁ THỜI GIAN! (Timeout) Đã hết 5 giây.")
			return // Thoát khỏi main
		}
	}
}
