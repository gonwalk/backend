import main

五种基础排序算法对比
![五种基础排序算法对比](https://upload-images.jianshu.io/upload_images/14576226-ebf818d343eb0343.png?imageMogr2/auto-orient/strip|imageView2/2/w/1200/format/webp)
1:冒泡排序
算法描述
比较相邻的元素。如果第一个比第二个大，就交换它们两个；
对每一对相邻元素作同样的工作，从开始第一对到结尾的最后一对，这样在最后的元素应该会是最大的数；
针对所有的元素重复以上的步骤，除了最后一个；
重复步骤1~3，直到排序完成。
动图演示
冒泡排序动图演示
代码演示
func bubbleSort(arr [6]int) {
    for i := 0; i < len(arr)-1; i++ {
        for j := 0; j < len(arr)-i-1; j++ {
            if arr[j] > arr[j+1] {
                temp := arr[j]
                arr[j] = arr[j+1]
                arr[j+1] = temp
            }
        }
    }
    fmt.Println(arr)
}


// 选择排序
算法描述
n个记录的直接选择排序可经过n-1趟直接选择排序得到有序结果。具体算法描述如下：

初始状态：无序区为R[1..n]，有序区为空；
第i趟排序(i=1,2,3…n-1)开始时，当前有序区和无序区分别为R[1..i-1]和R(i..n）。该趟排序从当前无序区中-选出关键字最小的记录 R[k]，将它与无序区的第1个记录R交换，使R[1..i]和R[i+1..n)分别变为记录个数增加1个的新有序区和记录个数减少1个的新无序区；
n-1趟结束，数组有序化了。
动图演示
选择排序
代码演示
func selectSort(arr [6]int) {
    for i := 0; i < len(arr)-1; i++ {
        min_index := i
        for j := i + 1; j < len(arr); j++ {
            if arr[i] > arr[j] {
                min_index = j
            }
            temp := arr[i]
            arr[i] = arr[min_index]
            arr[min_index] = temp
        }
    }
    fmt.Println(arr)
}
// 插入排序
算法描述
一般来说，插入排序都采用in-place在数组上实现。具体算法描述如下：

从第一个元素开始，该元素可以认为已经被排序；
取出下一个元素，在已经排序的元素序列中从后向前扫描；
如果该元素（已排序）大于新元素，将该元素移到下一位置；
重复步骤3，直到找到已排序的元素小于或者等于新元素的位置；
将新元素插入到该位置后；
重复步骤2~5。
动图演示
插入排序.gif
代码实现
func insertSort(arr [6]int) {
    for i := 0; i < len(arr); i++ {
        for j := i; j > 0; j-- {
            if arr[j] > arr[j-1] {
                temp := arr[j]
                arr[j] = arr[j-1]
                arr[j-1] = temp
            }
        }
    }
    fmt.Println(arr)
}



4:快速排序
快速排序的基本思想：通过一趟排序将待排记录分隔成独立的两部分，其中一部分记录的关键字均比另一部分的关键字小，则可分别对这两部分记录继续进行排序，以达到整个序列有序。

从数列中挑出一个元素，称为 “基准”（pivot）；
重新排序数列，所有元素比基准值小的摆放在基准前面，所有元素比基准值大的摆在基准的后面（相同的数可以到任一边）。在这个分区退出之后，该基准就处于数列的中间位置。这个称为分区（partition）操作；
递归地（recursive）把小于基准值元素的子数列和大于基准值元素的子数列排序。


// 快速排序
func quickSort(arr []int) []int {
    length := len(arr)
    if length <= 1 {
        return arr
    }
    middle := arr[0]
    var left []int
    var right []int
    for i := 1; i < length; i++ {
        if middle < arr[i] {
            right = append(right, []int{arr[i]}...)
        } else {
            left = append(left, []int{arr[i]}...)
        }
    }
    middle_s := []int{middle}

    left = quickSort(left)
    right = quickSort(right)

    result := append(append(left, middle_s...), right...)

    return result
}


type TreeNode struct {
    val int
    left *TreeNode
    right *TreeNode
}

func printBinaryTree(root *TreeNode){
    if root == nil{  // 根节点为空，打印程序退出
        return
    }
    treeNodeList := make([]*TreeNode, 0)  // 节点列表，存放当前层级的节点
    treeNodeList.append(root)             // 访问根节点
    if len(treeNodeList) > 0 {            // 当前层级是否已经遍历完
        curNode := treeNodelList[0]             // （从左到右）从当前层级中取出一个节点
        treeNodeList = treeNodeList[1:]   // 删除已访问的节点，筛选出当前层级中剩余未访问的节点
        fmt.Printf("%d", curNode.val) // 打印当前正在访问的节点
        if curNode.left != nil {
            treeNodeList.append(curNode.left)
        }
        if curNode.right != nil {
            treeNodeList.append(curNode.right)
        }
    }
}

func main(){

}