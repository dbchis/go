package main

import (
	"fmt"
	"math/rand"
	"time"
)

// Ý tưởng bài toán: Khu Ẩm thực "Buzzer"
// Bối cảnh: Có 3 Khách hàng (Customer) đến một khu ẩm thực. Họ xếp hàng (tưởng tượng) để đặt món tại một Quầy Bếp (Kitchen) duy nhất.

// Luồng hoạt động:

// Mỗi Khách hàng (một Goroutine) sẽ tạo một "Yêu cầu Món ăn" (Order).

// Khi đặt món, Khách hàng đưa Yêu cầu này vào Hàng đợi của Bếp (một Buffered Channel).

// Đồng thời, Khách hàng nhận được một "Thiết bị rung" (Buzzer) (một Unbuffered Channel cá nhân).

// Khách hàng đi về bàn và BLOCK (dừng lại) để chờ Buzzer của mình rung.

// Bếp (một Goroutine) xử lý các Yêu cầu trong Hàng đợi. Khi làm xong món, Bếp sẽ "gửi tín hiệu" (gửi Món ăn) qua Buzzer của đúng Khách hàng đó.

// Khách hàng đang chờ <-Buzzer sẽ được "unblock", nhận Món ăn và đi về.

// Thử thách (dùng select): Nếu Khách hàng chờ quá 3 giây mà Buzzer không rung, họ sẽ bực mình và bỏ về (timeout).
// --- Structs (Như bên trên) ---
type Dish struct {
	Content string
}

type Order struct {
	DishName string
	Buzzer   chan<- Dish // Kênh chỉ-gửi
}

// chan<- Dish: Mũi tên <- chỉ vào chan.
// Nghĩa là: Dữ liệu (Dish) chỉ có thể đi vào channel.
// Đây là kênh CHỈ-GỬI (Send-Only).
// <-chan Dish: Mũi tên <- đi ra từ chan.
// Nghĩa là: Dữ liệu (Dish) chỉ có thể đi ra khỏi channel.
// Đây là kênh CHỈ-NHẬN (Receive-Only).

// --- Các nhân vật ---
// customer (Khách hàng) - Một Goroutine
func customer(id int, dishName string, orderQueue chan<- Order) {
	fmt.Printf("[Khách %d] 🙋: Em muốn một '%s'\n", id, dishName)

	// 1. Tạo Buzzer CÁ NHÂN (Unbuffered)
	// Đây là kênh "trao tay" món ăn
	myBuzzer := make(chan Dish)

	// 2. Tạo Order, đính kèm Buzzer của mình
	order := Order{
		DishName: dishName,
		Buzzer:   myBuzzer, // Gửi kênh "chỉ-gửi" cho Bếp
	}

	// 3. Gửi Order vào Hàng đợi của Bếp (Buffered)
	orderQueue <- order
	fmt.Printf("[Khách %d] 👍: Đã đặt xong, ngồi chờ buzzer...\n", id)

	// 4. CHỜ (BLOCK) tín hiệu từ Buzzer
	//    SỬ DỤNG SELECT để chờ Buzzer HOẶC bị Timeout
	select {
	case dish := <-myBuzzer:
		// SỰ KIỆN 1: Buzzer rung, nhận được món
		fmt.Printf("[Khách %d] 😋: Đã nhận được '%s'. Ngon!\n", id, dish.Content)

	case <-time.After(5 * time.Second):
		// SỰ KIỆN 2: Chờ 3 giây mà không thấy gì
		fmt.Printf("[Khách %d] 😡: Đợi 3 giây rồi, lâu quá, BỎ VỀ!\n", id)
	}
}

// kitchen (Bếp) - Một Goroutine
func kitchen(orderQueue <-chan Order, nameKit string) {
	fmt.Println("[" + nameKit + "] 🔥: Bếp đã mở cửa! Sẵn sàng nhận đơn...")

	// Sử dụng 'for range' để liên tục xử lý các order trong hàng đợi
	// Vòng lặp này sẽ tự dừng khi 'orderQueue' bị đóng (close)
	for order := range orderQueue {

		fmt.Printf("["+nameKit+"] 🍳: Nhận được đơn '%s'. Bắt đầu nấu...\n", order.DishName)

		// Giả lập thời gian nấu (từ 1 đến 4 giây)
		cookTime := time.Duration(rand.Intn(3)+1) * time.Second
		time.Sleep(cookTime)

		// Món ăn hoàn thành
		dish := Dish{
			Content: "Món " + order.DishName + " nóng hổi (Nấu trong " + cookTime.String() + ")",
		}

		fmt.Printf("["+nameKit+"]🛎️: Món '%s' đã xong! Bấm buzzer...\n", order.DishName)

		// 5. Gửi (trao tay) món ăn qua Buzzer (Unbuffered)
		// Bếp sẽ BLOCK ở đây nếu Khách hàng... bỏ về (do timeout) và không nhận
		// Chúng ta sẽ xử lý việc này sau (với select), tạm thời cứ gửi
		order.Buzzer <- dish
	}

	fmt.Println("[" + nameKit + "]😴: Hết đơn, Bếp đóng cửa.")
}

// --- Hàm Main (Quản lý) ---
func main() {
	// Số lượng khách
	numCustomers := 10
	// Hàng đợi của Bếp (Buffered) - Tối đa 5 order cùng lúc
	orderQueue := make(chan Order, 5)

	// Danh sách món ăn
	menu := []string{"Phở", "Bún", "Cơm rang", "Mì xào"}

	// 1. Khởi động Bếp (1 Goroutine)
	// go kitchen(orderQueue)
	//1. Khởi động 2 Bếp (2 Goroutines)
	numKitchens := 4
	for k := 1; k <= numKitchens; k++ {
		go kitchen(orderQueue, fmt.Sprintf("Bếp-%d", k))
	}

	// 2. Các Khách hàng lần lượt đến (3 Goroutines)
	for i := 1; i <= numCustomers; i++ {
		// Chọn món ngẫu nhiên
		dishName := menu[rand.Intn(len(menu))]
		go customer(i, dishName, orderQueue)
		time.Sleep(500 * time.Millisecond) // Giả lập khách đến cách nhau
	}

	// Chờ một lúc cho các khách hàng đặt món
	time.Sleep(3 * time.Second)

	// 3. Đã hết giờ nhận đơn
	fmt.Println("\n[Quản lý] ⛔: HẾT GIỜ NHẬN ĐƠN. Đóng hàng đợi Bếp.\n")
	close(orderQueue) // Báo cho Bếp biết không còn order mới

	// 4. Chờ 10 giây để Bếp và Khách hàng xử lý xong
	// (Đây là cách đơn giản để main không bị thoát)
	time.Sleep(10 * time.Second)
	fmt.Println("[Quản lý] 🏁: Nhà hàng đóng cửa.")
}
