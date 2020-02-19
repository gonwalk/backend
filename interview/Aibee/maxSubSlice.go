package main

import (
	"fmt"
)


// golang动态规划求解最大连续子数组和

/*

*/

func main(){
	s := []int{1, 3, -2, 2, -1, 5, 2, -6, 7 }
	m := maxSubArray(s)
	fmt.Println(m)
}

// 求最大连续子数组和 https://blog.csdn.net/QQ245671051/article/details/72162320
func maxSubArray(arr []int) int {
    currSum := 0
    maxSum := arr[0]

    for _, v := range arr {
        if currSum > 0 {		// currSum表示当前遍历子串的和，如果为正数就继续累加
            currSum += v
        } else {				// 当前子序列的和为负数，则子序列的位置从当前遍历的元素位置开始重新计算和
            currSum = v
        }
        if maxSum < currSum {
            maxSum = currSum
        }
    }
    return maxSum
}


// leetcode(17):连续子串和/乘积最大：https://blog.csdn.net/qq_35082030/article/details/79975912
// 1.连续子串和
// （1）问题描述
/* Given an integer array nums, find the contiguous subarray (containing at least one number) which has the largest sum and return its sum.
给定一个整数数组 nums ，找到一个具有最大和的连续子数组（子数组最少包含一个数），返回其最大和。
Example:
Input: [-2,1,-3,4,-1,2,1,-5,4],
Output: 6
Explanation: [4,-1,2,1] has the largest sum = 6.
原文链接：https://blog.csdn.net/qq_35082030/article/details/79975912
*/

// （2）思路
/*
思路
这同样是一个动态规划的题目。我们只需要一个局部最大值，然后求解全局最大值即可。
关键是局部最优值的来源是什么：其一当前值，另一个是之前的局部最大值+当前值。
看局部最大值是来源于哪，然后再和全局最优进行比较。
*/

// （3）代码实现
public int maxSubArray(int[] nums) {
	if(nums==null||nums.length==0) return 0;
	#全局最大
	int max=Integer.MIN_VALUE;
	#局部最大
	int tempmax=0;
	for(int i=0;i<nums.length;i++){
		tempmax=tempmax+nums[i];
		//状态转移
		tempmax=tempmax>nums[i]?tempmax:nums[i];
		//和全局比较
		max=max>tempmax?max:tempmax;
	}
	return max;
}

// 2.Maximum Product Subarray
// （1）题目描述
/*
Find the contiguous subarray within an array (containing at least one number) which has the largest product.
找出一个序列中乘积最大的连续子序列（该序列至少包含一个数）。
For example, given the array [2,3,-2,4],
the contiguous subarray [2,3] has the largest product = 6.
*/

// （2）思路
/*
和第一题思路一样，只不过这里要考虑正负值，因为负负得正，这样最小值有可能直接变成最大值了。
*/

// （3）代码实现
// 原文链接：https://blog.csdn.net/qq_35082030/article/details/79975912
public int maxProduct(int[] nums) {
	if(nums==null){
		return 0;
	}
	int max=nums[0];
	int min=nums[0];
	int globalmax=nums[0];
	//int globalmin=nums[0];
	for(int i=1;i<nums.length;i++){
		int tempmax=max*nums[i];
		int tempmin=min*nums[i];
		//如果最小值和最大值互换了
		if(tempmin>tempmax){
			int temp=tempmax;
			tempmax=tempmin;
			tempmin=temp;
		}
		//局部最大，局部最小
		max=tempmax>nums[i]?tempmax:nums[i];
		min=tempmin<nums[i]?tempmin:nums[i];
		//全局最大
		globalmax=globalmax>max?globalmax:max;
		//全局最小
		//globalmin=globalmin<min?globalmin:min;
	}
	return globalmax;
}


// 3.子矩阵最大累加和
// （1）问题描述
/*
给定一个矩阵Matrix, 其中的值由正, 有负, 有0, 返回子矩阵的最大累加和.
*/

// （2）思路
/*
这个思路和第一题一样，只不过现在从1维变为2维了，我们由第一题可以知道，最大连续子串和的时间复杂度维O（n），那么二维遍历一次，也应当是O（n^2），因此最少也要O（n^3）才能完成。
我们首先考虑第一层的最大累加和，如同第一题一样就可以解决。随后在判断第一层和第二层一起的最大累加和，即把第一层和第二层按位加，就变成了1维的，然后同样采取第一题的做法。以此类推，直到加为最后一层即可。然后找到最大值，如果还需要记录该矩阵的大小，就在每次遍历时，加入当前遍历次数和当前最优的连续数目即可，最终再进行一个乘积运算。
*/

// （3）代码实现
public static int MaximumSum(int[][] Matrix){
	int max=Integer.MIN_VALUE;
	for(int i=0;i<Matrix.length;i++){
		int[] sum=new int[Matrix[0].length];
		for(int j=i;j<Matrix.length;j++){
			//进行累加
			for(int col=0;col<Matrix[0].length;col++){
				sum[col]+=Matrix[j][col];
			}
			int tempmax=maxSubArray(sum);
			max=max>tempmax?max:tempmax;
		}
	}
	return max;
}
/*
其实变成向量的乘积，也是这个道理，只不过把累加变成累乘即可。
原文链接：https://blog.csdn.net/qq_35082030/article/details/79975912
*/