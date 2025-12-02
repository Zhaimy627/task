func singleNumber(nums []int) int {
	result := 0
	for _, v := range nums {
		result ^= v
	}
	return result
}
func isPalindrome(x int) bool {
	// 负数或以 0 结尾但不是 0 的数都不是回文
	if x < 0 || (x%10 == 0 && x != 0) {
		return false
	}

	reverted := 0
	for x > reverted {
		reverted = reverted*10 + x%10
		x /= 10
	}

	// 当数字长度是奇数时，reverted 会多一位中间数字
	// 允许 reverted / 10 == x（去掉中间位）
	return x == reverted || x == reverted/10
}
func isValid(s string) bool {
	if len(s) == 0 {
		return true
	}
	if len(s)%2 == 1 {
		return false // 奇数长度不可能匹配
	}

	stack := []byte{}
	pairs := map[byte]byte{
		')': '(',
		'}': '{',
		']': '[',
	}

	for i := 0; i < len(s); i++ {
		ch := s[i]

		// 是右括号
		if left, ok := pairs[ch]; ok {
			if len(stack) == 0 || stack[len(stack)-1] != left {
				return false
			}
			stack = stack[:len(stack)-1] // 出栈
		} else {
			// 是左括号，入栈
			stack = append(stack, ch)
		}
	}

	return len(stack) == 0
}
func longestCommonPrefix(strs []string) string {
	if len(strs) == 0 {
		return ""
	}
	if len(strs) == 1 {
		return strs[0]
	}

	// 以第一个字符串为基准
	for i := 0; i < len(strs[0]); i++ {
		char := strs[0][i]

		// 遍历其他字符串的同一位置
		for j := 1; j < len(strs); j++ {
			// 情况1：当前字符串长度不够
			// 情况2：字符不匹配
			if i >= len(strs[j]) || strs[j][i] != char {
				return strs[0][:i] // 返回 [0, i)
			}
		}
	}

	// 所有字符串都包含 strs[0] 作为前缀
	return strs[0]
}
func longestCommonPrefix(strs []string) string {
	if len(strs) == 0 {
		return ""
	}
	if len(strs) == 1 {
		return strs[0]
	}

	// 以第一个字符串为基准
	for i := 0; i < len(strs[0]); i++ {
		char := strs[0][i]

		// 遍历其他字符串的同一位置
		for j := 1; j < len(strs); j++ {
			// 情况1：当前字符串长度不够
			// 情况2：字符不匹配
			if i >= len(strs[j]) || strs[j][i] != char {
				return strs[0][:i] // 返回 [0, i)
			}
		}
	}

	// 所有字符串都包含 strs[0] 作为前缀
	return strs[0]
}
func plusOne(digits []int) []int {
	n := len(digits)

	// 从最低位开始处理
	for i := n - 1; i >= 0; i-- {
		if digits[i] < 9 {
			digits[i]++ // 不进位，直接返回
			return digits
		}
		digits[i] = 0 // 进位，继续
	}

	// 所有位都是 9，例如 [9,9,9] → [1,0,0,0]
	return append([]int{1}, digits...)
}
func removeDuplicates(nums []int) int {
	if len(nums) == 0 {
		return 0
	}

	slow := 0
	for fast := 1; fast < len(nums); fast++ {
		if nums[fast] != nums[slow] {
			slow++
			nums[slow] = nums[fast] // 覆盖
		}
	}
	return slow + 1 // 长度
}