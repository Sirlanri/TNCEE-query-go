package main

//获取最大值
func getMax(nums []int) int {
	maxnum := nums[0]
	for i := 1; i < len(nums); i++ {
		if nums[i] > maxnum {
			maxnum = nums[i]
		}
	}
	return maxnum
}
func getMin(nums []int) int {
	minnum := nums[0]
	for i := 1; i < len(nums); i++ {
		if nums[i] < minnum {
			minnum = nums[i]
		}
	}
	return minnum
}
