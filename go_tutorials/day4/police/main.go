package main

import (
	"context"
	"fmt"
	"time"
)

// policeWorker là một Goroutine đại diện cho cảnh sát tuần tra
func policeWorker(ctxWithCancel context.Context, name string) {
	fmt.Printf("Cảnh sát %v: Bắt đầu tuần tra...", name)
	for {
		select {
		case <-ctxWithCancel.Done(): // 👈 Mấu chốt: Lắng nghe tín hiệu hủy
			// Khi context bị hủy (cancel), kênh ctx.Done() sẽ được "close"
			// và case này sẽ được kích hoạt.
			fmt.Printf("Cảnh sát %d: Nhận lệnh HỦY! Dừng lại.\n", name)
			return // Thoát khỏi hàm

		default:
			// Giả lập công việc
			fmt.Printf("Cảnh sát %v: Đang làm việc...\n", name)
			time.Sleep(1 * time.Second)
		}
	}
}
func main() {
	ctx := context.Background()
	ctxWithCancel, cancelFunc := context.WithCancel(ctx)
	go policeWorker(ctxWithCancel, "Đặng Bá Chí")
	time.Sleep(5 * time.Second)
	fmt.Println("MAIN: Gửi tín hiệu HỦY!")
	cancelFunc()
	time.Sleep(1 * time.Second)
	fmt.Println("MAIN: Kết thúc.")
}
