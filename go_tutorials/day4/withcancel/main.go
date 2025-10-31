package main

import (
	"context"
	"fmt"
	"time"
)

// worker là một Goroutine chạy mãi mãi cho đến khi bị bảo dừng
func worker(ctx context.Context, id int) {
	fmt.Printf("Worker %d: Bắt đầu...\n", id)
	for {
		select {
		case <-ctx.Done(): // 👈 Mấu chốt: Lắng nghe tín hiệu hủy
			// Khi context bị hủy (cancel), kênh ctx.Done() sẽ được "close"
			// và case này sẽ được kích hoạt.
			fmt.Printf("Worker %d: Nhận lệnh HỦY! Dừng lại.\n", id)
			return // Thoát khỏi hàm

		default:
			// Giả lập công việc
			fmt.Printf("Worker %d: Đang làm việc...\n", id)
			time.Sleep(1 * time.Second)
		}
	}
}

func main() {
	// 1. Tạo context gốc
	ctx := context.Background()

	// 2. Tạo context con có thể hủy
	ctxWithCancel, cancelFunc := context.WithCancel(ctx)

	// 3. Khởi động worker với context con
	go worker(ctxWithCancel, 1)

	// 4. Chờ 3 giây
	time.Sleep(3 * time.Second)

	// 5. Gửi tín hiệu HỦY
	fmt.Println("MAIN: Gửi tín hiệu HỦY!")
	cancelFunc() // 👈 Gọi hàm cancel!

	// 6. Chờ worker thoát
	time.Sleep(1 * time.Second)
	fmt.Println("MAIN: Kết thúc.")
}
