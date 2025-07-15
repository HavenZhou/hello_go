package main

import (
	"fmt"
	"sort"
	"strconv"
)

func main1() {
	// 题目01：数组中唯一数
	arrInt := [...]int{1, 1, 5, 2, 2, 4, 4, 2}
	fmt.Println("the single num =", singleNum(arrInt[:]))

	// 题目02：判断一个整数是否是回文数
	fmt.Println(isPalindrome(12344321))

	// 题目03：给定一个只包括 '('，')'，'{'，'}'，'['，']' 的字符串，判断字符串是否有效
	// 分析 JSON字符串格式校验
	fmt.Println(isValid("{[()]}")) // true

	// 题目04：查找字符串数组中的最长公共前缀
	strArr := []string{"我们的abcd4211", "我们abcdasda", "我们abcd12313"}
	fmt.Println("查找字符串数组中的最长公共前缀为：", longestCommonPrefix(strArr)) // true

	// 题目05：给定一个由整数组成的非空数组所表示的非负整数，在该数的基础上加一
	// 例如：输入：[1, 2, 3] → 输出：[1, 2, 4]（123 + 1 = 124）
	intArr := []int{1, 2, 3, 9, 9}
	fmt.Println("数组加1后：", plusOne(intArr))

	// 题目06：删除重复元素
	intArr06 := []int{1, 2, 2, 3, 3, 4, 9, 9}
	fmt.Println("删除重复元素后，新数组长度为：", removeDuplicates(intArr06), intArr06)

	// 题目07：. 合并区间
	intervals := [][]int{{1, 3}, {2, 6}, {8, 10}, {9, 14}, {15, 18}}
	fmt.Println("合并区间后，intervals = ", merge(intervals))

	// 题目08：给定一个整数数组 nums 和一个目标值 target，请你在该数组中找出和为目标值的那两个整数
	intArr08 := [...]int{1, 2, 3, 4, 5, 6, 7}
	fmt.Println("两个整数为", twoSum(intArr08[:], 8))
}

// 题目1：
// 给定一个非空整数数组，除了某个元素只出现一次以外，其余每个元素均出现两次。找出那个只出现了一次的元素。
// 可以使用 for 循环遍历数组，结合 if 条件判断和 map 数据结构来解决，例如通过 map 记录每个元素出现的次数，然后再遍历 map 找到出现次数为1的元素
func singleNum(arrInt []int) int {
	// 定义map
	freq := make(map[int]int) // 频率计算

	// 计算频率
	for _, num := range arrInt {
		freq[num]++
	}

	// 遍历map，找到只出现一次的数值
	for num, count := range freq {
		if count == 1 {
			return num
		}
	}

	return -1
}

// 题目02 判断回文数
func isPalindrome(num int) bool {
	if num < 0 {
		return false
	}

	str := strconv.Itoa(num)
	left, right := 0, len(str)-1

	for left < right {
		if str[left] != str[right] {
			return false
		}

		left++
		right--
	}
	return true
}

// 题目03：给定一个只包括 '('，')'，'{'，'}'，'['，']' 的字符串，判断字符串是否有效
// "{[()]}"
func isValid(s string) bool {
	stack := []rune{}
	pairs := map[rune]rune{
		')': '(',
		'}': '{',
		']': '[',
	}

	for _, char := range s {
		switch char {
		case '(', '{', '[':
			stack = append(stack, char)
		case ')', '}', ']':
			if len(stack) == 0 || stack[len(stack)-1] != pairs[char] {
				return false
			}
			stack = stack[:len(stack)-1]
		}
	}
	return true
}

// 题目04：查找字符串数组中的最长公共前缀
func longestCommonPrefix(strs []string) string {
	if len(strs) == 0 {
		return ""
	}

	base := strs[0]
	for i := 1; i <= len(base); i++ {
		str := base[:i]
		for j := 1; j <= len(strs)-1; j++ {
			if !(str == strs[j][:i]) {
				return str[:i-1]
			}
		}

		if i == len(base) {
			return base
		}
	}
	return ""
}

// 题目05：给定一个由整数组成的非空数组所表示的非负整数，在该数的基础上加一
func plusOne(digits []int) []int {
	n := len(digits)
	for i := n - 1; i >= 0; i-- {
		if digits[i] < 9 {
			digits[i]++
			return digits
		}
		digits[i] = 0
	}

	return append([]int{1}, digits...)
}

// 题目06：删除有序数组中的重复项：
// 给你一个有序数组 nums ，请你原地删除重复出现的元素，使每个元素只出现一次，返回删除后数组的新长度。
// 不要使用额外的数组空间，你必须在原地修改输入数组并在使用 O(1) 额外空间的条件下完成。
// 可以使用双指针法，一个慢指针 i 用于记录不重复元素的位置，一个快指针 j 用于遍历数组，
// 当 nums[i] 与 nums[j] 不相等时，将 nums[j] 赋值给 nums[i + 1]，并将 i 后移一位。
func removeDuplicates(nums []int) int {
	fmt.Println("-----------nums=", nums)
	if len(nums) == 0 {
		return 0
	}
	i := 0 // 慢指针
	for j := 1; j < len(nums); j++ {
		if nums[i] != nums[j] {
			i++
			nums[i] = nums[j]
			fmt.Println("i=", i, "j=", j, "nums=", nums)
		}
	}
	return i + 1
}

// 题目07：合并区间：以数组 intervals 表示若干个区间的集合，其中单个区间为 intervals[i] = [starti, endi] 。
// 请你合并所有重叠的区间，并返回一个不重叠的区间数组，该数组需恰好覆盖输入中的所有区间。
// 可以先对区间数组按照区间的起始位置进行排序，然后使用一个切片来存储合并后的区间，
// 遍历排序后的区间数组，将当前区间与切片中最后一个区间进行比较，
// 如果有重叠，则合并区间；如果没有重叠，则将当前区间添加到切片中。
// 输入：intervals = [[1,3],[2,6],[8,10],[15,18]]
// 输出：[[1,6],[8,10],[15,18]]
func merge(intervals [][]int) [][]int {
	if len(intervals) == 0 {
		return nil
	}

	// 按区间排序
	sort.Slice(intervals, func(i, j int) bool {
		return intervals[i][0] < intervals[j][0]
	})

	// 排序后合并空间
	merged := [][]int{intervals[0]}
	for i := 1; i < len(intervals); i++ {
		last := merged[len(merged)-1]
		if last[len(last)-1] >= intervals[i][0] {
			last[len(last)-1] = intervals[i][len((intervals[i]))-1]
		} else {
			merged = append(merged, intervals[i])
		}
	}

	return merged
}

// 题目08：给定一个整数数组 nums 和一个目标值 target，请你在该数组中找出和为目标值的那两个整数
func twoSum(nums []int, target int) map[int]int {
	numMap := make(map[int]int)
	for i, num := range nums {
		complete := target - num
		if _, ok := numMap[complete]; ok {
			return map[int]int{num: complete}
		}
		numMap[num] = i
	}
	return nil
}
