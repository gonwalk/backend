
把正整数数组的元素拼成一个最小的数：https://blog.csdn.net/brucebupt/article/details/78024675

转自几个面试经典算法题（http://www.cnblogs.com/sunniest/p/4596182.html）题目四

/*
这个算法的思想和快速排序的思想相似。先把第一个元素p当作中间数，数组前后各有一个索引begin、end。先从后索引对应元素e起，如果ep>pe或e=p，那么后索引往前移动，直到后索引等于前索引，或ep<pe，交换e、p元素。再从前索引对应元素b起，如果pb>bp或b=p，那么前索引后移，直到前索引等于后索引或pb<bq，交换p、b。这时，“小于”p的都在p左边（这里小于的意思是，两个数a、b拼成一个数时，a小于b指ab<ba），“大于”p的都在p右边（这里大于的意思是，两个数a、b拼成一个数时，a大于b指ab>ba）。再对p的两边子数组各自进行递归处理。
*/

package main

import (
	"fmt"
	"strconv"
)

func main() {
	var a = []int {2, 3, 1, 6, 2, 5}
	quickSortMin(a, 0, len(a) - 1)
	for i := range a {
		fmt.Print(a[i])
	}
	fmt.Println()
}

func MBiggerThanN(m, n int) bool {
	if m == n {
		return false
	}

	sm := strconv.Itoa(m)
	sn := strconv.Itoa(n)
	s1 := sm + sn
	s2 := sn + sm

	i1, err1 := strconv.Atoi(s1)
	if err1 != nil {
		fmt.Println("first parameter transfer fails")
		return false
	}
	i2, err2 := strconv.Atoi(s2)
	if err2 != nil {
		fmt.Println("seconde parameter transfer fails")
		return false
	}

	if i1 > i2 {
		return true
	} else {
		return false
	}
}

func getPivotPosition(array []int, left int, right int) int {
	pivot := array[left]
	for left < right {
		for (left < right && (pivot == array[right] || MBiggerThanN(array[right], pivot))) {
			right--
		}
		array[left], array[right] = array[right], array[left]
		for (left < right && (pivot == array[left] || MBiggerThanN(pivot, array[left]))) {
			left++
		}
		array[left], array[right] = array[right], array[left]
	}
	array[left] = pivot
	return left
}

func quickSortMin(array []int, left int, right int ) {
	if left < right {
		position := getPivotPosition(array, left, right)
		quickSort(array, left, position - 1)
		quickSort(array, position + 1, right)
	}
}

/*
问题描述：输入一个正整数数组，将它们连接起来排成一个数，输出能排出的所有数字中最小的一个。例如输入数组{32,  321}，则输出这两个能排成的最小数字32132。请给出解决问题的算法，并证明该算法。

      思路：先将整数数组转为字符串数组，然后字符串数组进行排序，最后依次输出字符串数组即可。这里注意的是字符串的比较函数需要重新定义，不是比较a和b，而是比较ab与 ba。如果ab < ba，则a < b；如果ab > ba，则a > b；如果ab = ba，则a = b。比较函数的定义是本解决方案的关键。

      证明：为什么这样排个序就可以了呢？简单证明一下。根据算法，如果a < b，那么a排在b前面，否则b排在a前面。可利用反证法，假设排成的最小数字为xxxxxx，并且至少存在一对字符串满足这个关系：a > b，但是在组成的数字中a排在b前面。根据a和b出现的位置，分三种情况考虑：

      （1）xxxxab，用ba代替ab可以得到xxxxba，这个数字是小于xxxxab，与假设矛盾。因此排成的最小数字中，不存在上述假设的关系。

      （2）abxxxx，用ba代替ab可以得到baxxxx，这个数字是小于abxxxx，与假设矛盾。因此排成的最小数字中，不存在上述假设的关系。

      （3）axxxxb，这一步证明麻烦了一点。可以将中间部分看成一个整体ayb，则有ay < ya，yb < by成立。将ay和by表示成10进制数字形式，则有下述关系式，这里a，y，b的位数分别为n，m，k。

        关系1： ay < ya => a * 10^m + y < y * 10^n + a => a * 10^m - a < y * 10^n - y => a( 10^m - 1)/( 10^n - 1) < y

        关系2： yb < by => y * 10^k + b < b * 10^m + y => y * 10^k - y < b * 10^m - b => y < b( 10^m -1)/( 10^k -1) 

        关系3： a( 10^m - 1)/( 10^n - 1) < y < b( 10^m -1)/( 10^k -1)  => a/( 10^n - 1)< b/( 10^k -1) => a*10^k - a < b * 10^n - b =>a*10^k + b < b * 10^n + a => a < b

       这与假设a > b矛盾。因此排成的最小数字中，不存在上述假设的关系。

       综上所述，得出假设不成立，从而得出结论：对于排成的最小数字，不存在满足下述关系的一对字符串：a > b，但是在组成的数字中a出现在b的前面。从而得出算法是正确的。

解题笔记（25）——把数组排成最小的数：https://blog.csdn.net/wuzhekai1985/article/details/6704902
*/