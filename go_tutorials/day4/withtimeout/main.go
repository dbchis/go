package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

func checkWebsite(ctx context.Context, wg *sync.WaitGroup, url string) {
	defer wg.Done()

	// 1. Tạo HTTP request, nhưng GẮN context vào
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		fmt.Printf("[LỖI] %s: Không tạo được request: %v\n", url, err)
		return
	}

	fmt.Printf("[Bắt đầu] %s\n", url)

	// 2. Thực thi request
	// Hầu hết các thư viện I/O (DB, HTTP) đều hỗ trợ Context
	resp, err := http.DefaultClient.Do(req)

	// 3. Xử lý kết quả
	// `select` ở đây để kiểm tra LỖI là gì
	select {
	case <-ctx.Done():
		// Lỗi này xảy ra VÌ context bị hủy (timeout)
		fmt.Printf("[TIMEOUT] %s: Đã quá 2 giây!\n", url)
		// ctx.Err() sẽ cho biết lý do (ví dụ: "context deadline exceeded")

	default:
		// Context chưa bị hủy, vậy lỗi là do thứ khác (ví dụ: DNS, 404)
		if err != nil {
			fmt.Printf("[LỖI] %s: %v\n", url, err)
			return
		}

		// Thành công!
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Printf("[THÀNH CÔNG] %s: %d bytes\n", url, len(body))
		resp.Body.Close()
	}
}

func main() {
	var wg sync.WaitGroup
	urls := []string{
		"http://google.com",          // Nhanh, sẽ xong
		"http://httpbin.org/delay/5", // Chậm, sẽ bị timeout
	}

	// 1. Tạo context cha
	ctx := context.Background()

	for _, url := range urls {
		wg.Add(1)

		// 2. TẠO MỘT TIMEOUT MỚI CHO MỖI GOROUTINE
		// Mỗi website có 2 giây để phản hồi
		// ctxWithTimeout, cancel := context.WithTimeout(ctx, 2*time.Second)
		ctxWithTimeout, cancel := context.WithDeadline(ctx, time.Now().Add(5*time.Second))

		// 3. (Rất quan trọng) Dọn dẹp context khi worker xong
		// Dù nó thành công hay timeout, cancel() vẫn cần được gọi
		// để giải phóng tài nguyên.
		// Chúng ta dùng một func bọc lại để `defer` đúng
		go func(url string, ctx context.Context, cancel context.CancelFunc) {
			defer cancel() // 👈 Dọn dẹp
			checkWebsite(ctx, &wg, url)
		}(url, ctxWithTimeout, cancel)
	}

	wg.Wait()
	fmt.Println("MAIN: Hoàn tất.")
}
