# Tóm tắt kiến thức: `sync` Package (Ngày 2)

Ngày 2 tập trung vào cách các Goroutine **đồng bộ hóa (synchronize)** và **bảo vệ (protect)** dữ liệu chia sẻ. Chúng ta sử dụng các công cụ từ `sync` package để thay thế cho việc giao tiếp qua channel.

---

## 1. `sync.WaitGroup` (Người giám sát)

`WaitGroup` là một cơ chế đồng bộ hóa đơn giản, hoạt động như một "bộ đếm công việc". Nó dùng để **chờ một nhóm Goroutine hoàn thành** trước khi tiếp tục.

### Vấn đề nó giải quyết

Nếu không có cơ chế chờ, hàm `main` (Goroutine chính) sẽ kết thúc ngay lập tức, và "giết" chết tất cả các Goroutine con đang chạy dở.

> **🛑 Lỗi phổ biến:** Dùng `time.Sleep()` ở cuối `main` để "cầu may". Cách này không đáng tin cậy và lãng phí tài nguyên.

### Phép so sánh (Analogy)

`WaitGroup` giống như một người quản lý với một bộ đếm:
1.  **`wg.Add(n)`:** Quản lý nói: "Tôi có `n` việc cần làm."
2.  **`wg.Done()`:** Mỗi nhân viên (Goroutine) khi làm xong báo cáo: "Tôi xong 1 việc." (Bộ đếm giảm đi 1).
3.  **`wg.Wait()`:** Quản lý **đứng chờ (BLOCK)** tại chỗ cho đến khi bộ đếm về `0`.

### Các phương thức chính

* **`var wg sync.WaitGroup`**: Khai báo.
* **`wg.Add(n int)`**: Tăng bộ đếm lên `n` (thường gọi ở `main` *trước khi* chạy Goroutine).
* **`wg.Done()`**: Giảm bộ đếm đi 1 (luôn gọi ở Goroutine con, thường dùng với `defer`).
* **`wg.Wait()`**: Block Goroutine hiện tại cho đến khi bộ đếm về 0.

### Ví dụ: Website Checker (Phiên bản `WaitGroup`)


    ```
    package main

    import (
        "fmt"
        "net/http"
        "sync"
    )

    func worker(wg *sync.WaitGroup, id int, url string) {
        // 3. Đảm bảo báo "Done" khi hàm kết thúc
        defer wg.Done()

        fmt.Printf("Worker %d đang xử lý: %s\n", id, url)
        if _, err := http.Get(url); err != nil {
            fmt.Printf("[LỖI] Worker %d, %s\n", id, url)
        } else {
            fmt.Printf("[OK] Worker %d: %s\n", id, url)
        }
    }

    func main() {
        urls := []string{
            "[http://google.com](http://google.com)",
            "[http://facebook.com](http://facebook.com)",
            "[http://github.com](http://github.com)",
        }

        // 1. Khai báo WaitGroup
        var wg sync.WaitGroup

        for i, url := range urls {
            // 2. Thêm 1 công việc vào bộ đếm
            wg.Add(1)
            go worker(&wg, i+1, url)
        }

        // 4. Chờ (BLOCK) cho đến khi bộ đếm về 0
        wg.Wait()

        fmt.Println("Tất cả website đã được kiểm tra xong!")
    }```

---


## 2. `sync.Mutex` (Ổ khóa)

`Mutex` (Mutual Exclusion - Loại trừ lẫn nhau) là một **ổ khóa**. Nó dùng để bảo vệ "vùng nguy hiểm" (critical section) - nơi có biến đang được chia sẻ.

### Vấn đề nó giải quyết

**Race Condition (Điều kiện Đua)**. Đây là lỗi xảy ra khi:
* Nhiều Goroutine...
* Truy cập (đọc/ghi) cùng một biến...
* Tại cùng một thời điểm...
* Và *ít nhất một trong số chúng là hành động GHI*.

Hậu quả là dữ liệu cuối cùng bị sai lệch, không thể đoán trước.


> **Phát hiện lỗi:** Go cung cấp công cụ để phát hiện Race Condition:
> `go run -race ten_file_cua_ban.go`

### Phép so sánh (Analogy)

`Mutex` là **"khóa cửa nhà vệ sinh"**:
1.  **`mu.Lock()`:** Ai muốn vào (truy cập biến) phải khóa cửa.
2.  Nếu cửa đã khóa, người đến sau phải **xếp hàng chờ (BLOCK)**.
3.  Làm xong, người bên trong phải **`mu.Unlock()`** để người tiếp theo vào.



### Các phương thức chính

* **`var mu sync.Mutex`**: Khai báo.
* **`mu.Lock()`**: Khóa. Nếu đang bị khóa, Goroutine sẽ chờ.
* **`mu.Unlock()`**: Mở khóa.

### Ví dụ: Sửa lỗi "Bộ đếm"

    ```go
    var counter int
    var mu sync.Mutex // Khai báo khóa
    var wg sync.WaitGroup

    func increment() {
        // 1. Khóa lại trước khi vào "vùng nguy hiểm"
        mu.Lock()
        
        // 2. Dùng defer để đảm bảo LUÔN MỞ KHÓA
        // (Kể cả khi hàm panic, nó vẫn được gọi)
        defer mu.Unlock() 
        
        counter++ // Chỉ một Goroutine được làm điều này tại một thời điểm
    }

    func main() {
        for i := 0; i < 1000; i++ {
            wg.Add(1)
            go func() {
                defer wg.Done()
                increment()
            }()
        }
        
        wg.Wait()
        fmt.Println("Kết quả cuối cùng:", counter) // Luôn là 1000
    }```

###Lưu ý về defer: defer mu.Unlock() là một pattern rất an toàn. Nó đảm bảo khóa luôn được mở, ngay cả khi hàm increment bị panic (lỗi nghiêm trọng) giữa chừng. Nếu không defer, khóa có thể bị "kẹt" vĩnh viễn, gây ra deadlock.

---

## 3. `sync.RWMutex` (Khóa Đọc/Viết)

`RWMutex` là một phiên bản `Mutex` được tối ưu hóa cho kịch bản **Đọc nhiều, Ghi ít**.

### Vấn đề nó giải quyết

`Mutex` quá nghiêm ngặt. Nó bắt cả những Goroutine chỉ muốn **ĐỌC** (một hành động vô hại) cũng phải xếp hàng 1-1.

`RWMutex` cho phép **nhiều độc giả (Readers)** vào cùng lúc, nhưng **chỉ một người viết (Writer)** được vào (và khi đó không ai được đọc).

### Phép so sánh (Analogy)

`RWMutex` là **"Thư viện"**:

* **Đọc (Read Lock):** Nhiều người (`mu.RLock()`) có thể vào đọc sách cùng lúc.
* **Ghi (Write Lock):** Khi nhân viên (`mu.Lock()`) muốn *cập nhật* sách, anh ta sẽ:
    1.  Chờ tất cả độc giả hiện tại ra ngoài (`mu.RLock()` sẽ chặn các độc giả mới).
    2.  Khóa cửa thư viện (độc quyền).
    3.  Cập nhật sách.
    4.  Mở cửa (`mu.Unlock()`).



### Các phương thức chính

* `var mu sync.RWMutex`: Khai báo.
* **Khi ĐỌC:**
    * `mu.RLock()`: Khóa đọc.
    * `mu.RUnlock()`: Mở khóa đọc.
* **Khi GHI:**
    * `mu.Lock()`: Khóa ghi (giống `Mutex`).
    * `mu.Unlock()`: Mở khóa ghi (giống `Mutex`).

### Ví dụ: Quản lý Config Map

    ```go
    var config = make(map[string]string)
    var mu sync.RWMutex // Dùng RWMutex
    var wg sync.WaitGroup

    // Hàng trăm Goroutine có thể gọi hàm này cùng lúc
    func readConfig(key string) string {
        mu.RLock() // Lấy khóa ĐỌC (nhiều người vào được)
        defer mu.RUnlock()
        return config[key]
    }

    // Chỉ một Goroutine gọi hàm này
    func writeConfig(key, value string) {
        mu.Lock() // Lấy khóa GHI (độc quyền, chờ độc giả ra hết)
        defer mu.Unlock()
        config[key] = value
    }```