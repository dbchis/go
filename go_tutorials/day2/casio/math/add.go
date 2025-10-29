// Tên package phải trùng với tên thư mục
package math

// TÊN HÀM VIẾT HOA ĐỂ EXPORT
// Nếu bạn viết "add" (thường), package main sẽ không thấy nó.
func Add(a, b int) int {
	return a + b
}

// Hàm này (viết thường) sẽ là private, chỉ dùng nội bộ trong package mathutil
func internalHelper() {
	// ...
}
