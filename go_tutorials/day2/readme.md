# Tóm tắt kiến thức: Goroutine & Channel (Ngày 1)

Đây là bản tóm tắt các khái niệm cốt lõi về lập trình đồng thời trong Go, tập trung vào Goroutine, Channels, và cách chúng tương tác với nhau.

## 1. Goroutine là gì?

* **Định nghĩa:** Một Goroutine là một "luồng" (thread) siêu nhẹ được quản lý bởi Go Runtime (Bộ lập lịch của Go), không phải bởi hệ điều hành.

* **Đặc điểm:**
    * **Siêu nhẹ:** Khởi tạo chỉ với vài KB bộ nhớ (stack), rẻ hơn nhiều so với luồng của HĐH (thường là 1-2MB).
    * **Siêu nhanh (khởi tạo):** Tạo ra một goroutine nhanh hơn nhiều so vơi tạo 1 luồng.
    * **Quản lý (Scheduling):** Go Scheduler tự động phân bổ (multiplex) hàng ngàn Goroutine lên một số lượng nhỏ luồng của HĐH.

* **Cách tạo:** Sử dụng từ khóa `go`.

    ```go
    // Chạy hàm này trên 1 goroutine mới
    go myFunction() 
    
    // Dùng hàm ẩn danh
    go func(msg string) { 
        fmt.Println(msg)
    }("Xin chào")
    ```

* **Lưu ý:** Khi hàm `main` kết thúc, tất cả các goroutine khác cũng bị "giết" ngay lập tức, bất kể chúng đã chạy xong hay chưa.

---

## 2. Channel (Kênh) là gì?

* **Định nghĩa:** Channel là một "ống dẫn" (pipe) có định kiểu (typed), dùng để **giao tiếp và đồng bộ hóa** giữa các Goroutine. Đây là cách an toàn để truyền dữ liệu qua lại.

* **Triết lý của Go:**
    > "Do not communicate by sharing memory; instead, share memory by communicating."
    > (Đừng giao tiếp bằng cách chia sẻ bộ nhớ; hãy chia sẻ bộ nhớ bằng cách giao tiếp.)

* **Cú pháp cơ bản:**

    ```go
    // Khởi tạo (Tạo 1 channel kiểu int)
    ch := make(chan int) 
    
    // Gửi (Send) (Gửi giá trị 10 *vào* channel)
    ch <- 10 
    
    // Nhận (Receive) (Chờ và nhận 1 giá trị *từ* channel)
    value := <-ch 
    
    // Đóng (Close) (Báo hiệu rằng sẽ không có giá trị nào được gửi nữa)
    close(ch) 
    ```

---

## 3. Hai loại Channel (Cực kỳ quan trọng)

### A. Unbuffered Channel (Kênh không đệm)

* **Khởi tạo:** `ch := make(chan string)`
* **Sức chứa (Capacity):** Bằng 0.
* **Phép so sánh:** Cuộc **"Trao tay trực tiếp"** (Direct Handover / Rendezvous).
* **Hành vi:**
    1.  **Gửi (Send) `ch <- "data"`:** Người gửi sẽ **BLOCK (DỪNG LẠI)**. Nó sẽ đứng chờ cho đến khi *chính xác tại thời điểm đó* có một người nhận (`<-ch`) đến để lấy dữ liệu.
    2.  **Nhận (Receive) `<-ch`:** Người nhận sẽ **BLOCK (DỪNG LẠI)**. Nó sẽ đứng chờ cho đến khi *chính xác tại thời điểm đó* có một người gửi (`ch <- "data"`) đến để đưa dữ liệu.
* **Hệ quả:** Đây là một hành động **đồng bộ hóa (synchronization)**. Cả hai bên (gửi và nhận) phải gặp nhau. Dữ liệu không bao giờ "nằm" trong channel.
    
* **Khi nào dùng:** Khi bạn cần đồng bộ hóa. Ví dụ: "Worker A phải hoàn thành và *báo cáo* cho Worker B, và Worker B chỉ được chạy *sau khi* nhận được báo cáo."

### B. Buffered Channel (Kênh có đệm)

* **Khởi tạo:** `ch := make(chan string, 3)` (Kích thước buffer là 3)
* **Sức chứa (Capacity):** Lớn hơn 0 (ví dụ: 3).
* **Phép so sánh:** Một **"Hàng đợi" (Queue)**, **"Kho chứa"** hoặc **"Cái khay"**.
* **Hành vi:**
    1.  **Gửi (Send) `ch <- "data"`:**
        * Nếu kho *chưa đầy* (ví dụ: đang chứa 2/3): Người gửi "quăng" dữ liệu vào kho và **KHÔNG BỊ BLOCK**. Nó đi làm việc khác ngay lập tức.
        * Nếu kho *đã đầy* (ví dụ: đang chứa 3/3): Người gửi sẽ **BLOCK**, chờ cho đến khi có người nhận lấy bớt 1 món ra để giải phóng chỗ.
    2.  **Nhận (Receive) `<-ch`:**
        * Nếu kho *có hàng* (ví dụ: đang chứa 1/3): Người nhận lấy 1 món ra và **KHÔNG BỊ BLOCK**.
        * Nếu kho *trống rỗng* (0/3): Người nhận sẽ **BLOCK**, chờ cho đến khi có người gửi bỏ 1 món mới vào.
* **Hệ quả:** Đây là một hành động **tách rời (decoupling)**. Người gửi và người nhận không cần gặp nhau. Họ làm việc độc lập qua cái kho.
    
* **Khi nào dùng:** (Như dự án Website Checker) Dùng trong Worker Pool. `main` (người gửi) có thể "quăng" 100 công việc vào `jobs` (buffered channel) và đi làm việc khác, trong khi các `worker` (người nhận) từ từ lấy việc ra xử lý.

---

## 4. `select` - Chờ nhiều kênh cùng lúc

* **Định nghĩa:** `select` là một câu lệnh cho phép một Goroutine chờ trên nhiều hoạt động giao tiếp (gửi hoặc nhận) cùng một lúc.
* **Phép so sánh:** Một nhân viên trực tổng đài, chờ nhiều đường dây điện thoại (channel) cùng lúc. Hễ đường dây nào reo (có tín hiệu) trước, anh ta sẽ nhấc máy đó.
* **Cú pháp & Hành vi:**

    ```go
    select {
    case data1 := <-channel1:
        // Xử lý data1 từ channel1
        
    case data2 := <-channel2:
        // Xửlý data2 từ channel2
        
    case channel3 <- "gửi đi":
        // Đã gửi thành công vào channel3
        
    case <-time.After(1 * time.Second):
        // Hết 1 giây mà không có case nào ở trên xảy ra
        // (Đây là cách làm timeout phổ biến)
        
    default:
        // (Nếu có `default`)
        // Nếu không có channel nào sẵn sàng ngay lập tức, chạy case này
        // Làm cho `select` không bị block.
    }
    ```

* **Lưu ý quan trọng:**
    1.  `select` sẽ **BLOCK** cho đến khi *một* trong các `case` của nó có thể thực thi.
    2.  Nếu nhiều `case` sẵn sàng cùng lúc, `select` sẽ chọn **ngẫu nhiên** một trong số chúng.

---

## 5. Bài học từ dự án Website Checker

* **`go worker(...)`:** Bạn đã tạo ra nhiều Goroutine (`numWorkers`) để chạy song song/đồng thời.
* **`jobs := make(chan string, numJobs)`:** Bạn dùng **Buffered Channel** để làm hàng đợi công việc. Điều này rất hiệu quả, vì `main` có thể "quăng" tất cả URL vào đây rất nhanh và `close(jobs)` mà không bị block.
* **`results := make(chan string, numJobs)`:** Bạn dùng **Buffered Channel** để nhận kết quả. Điều này cũng hiệu quả, vì các `worker` có thể "quăng" kết quả vào đây và quay lại lấy việc mới ngay, không cần chờ `main` nhận.
* **`for url := range jobs` (trong `worker`):** Đây là cách tuyệt vời để nhận việc. Vòng `for range` tự động chờ và lấy giá trị từ channel, và tự động kết thúc khi channel `jobs` được `close()`.
* **Vòng lặp `for...; <-results` (trong `main`):** Bạn đã hiểu rằng toán tử `<-results` luôn chỉ lấy ra **một** kết quả tại một thời điểm (theo thứ tự FIFO - Vào trước Ra trước). Đó là lý do bạn cần lặp lại `numJobs` lần để gom đủ kết quả.
* **`select` với `time.After`:** Bạn đã dùng `select` để chờ giữa hai sự kiện: "nhận được kết quả" (`<-results`) HOẶC "hết giờ" (`<-timeout`). Đây là một pattern cực kỳ phổ biến và mạnh mẽ trong thực tế.

go mod init <project-name>