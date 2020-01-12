package main

import (
    "container/list"
    "fmt"
    "strings"
)
// container/list包中的list.List结构内部是双端链表，既可以左右栈使用，又可以当成队列使用
type MyStack struct {
    List *list.List
}

type MyQueue struct {
    List *list.List
}

type BinaryTree struct {
    Value interface{}
    Left *BinaryTree
    Right *BinaryTree
}

type Tree struct {
    Value interface{}
    Children []*Tree
}

// 栈中的元素从顶进从顶出（先进后出）
func (stack *MyStack) pop() interface{} {
    // 使用List的Back方法获取栈顶元素
    if elem := stack.List.Back(); elem != nil {
        // 将列表（实质是双端链表）中的最后一个元素删除
        stack.List.Remove(elem)
        return elem.Value
    }
    return nil
}

func (stack *MyStack) push(elem interface{}) {
    stack.List.PushBack(elem)
}

// 队列从头部出、尾部进
func (queue *MyQueue) pop() interface{} {
    if elem := queue.List.Front(); elem != nil {
        queue.List.Remove(elem)
        return elem.Value
    }
    return nil
}

func (queue *MyQueue) push(elem interface{}) {
    queue.List.PushBack(elem)
}

// 树的先序遍历（递归实现）
func preOrderRecur(node *BinaryTree) {
    if node == nil {
        return
    }

    fmt.Println(node.Value)
    preOrderRecur(node.Left)
    preOrderRecur(node.Right)
}

// 树的中序遍历（递归实现）
func inOrderRecu(node *BinaryTree) {
    if node == nil {
        return
    }

    inOrderRecu(node.Left)
    fmt.Println(node.Value)
    inOrderRecu(node.Right)
}

// 树的后序遍历（递归实现）
func posOrderRecu(node *BinaryTree) {
    if node == nil {
        return
    }

    posOrderRecu(node.Left)
    posOrderRecu(node.Right)
    fmt.Println(node.Value)
}

// 树的先序遍历（非递归实现）
func preOrder(node *BinaryTree) {
    stack := MyStack{List: list.New()}
    stack.push(node)

    elem := stack.pop()
    for elem != nil {
        node, _ := elem.(*BinaryTree)

        fmt.Println(node.Value)
        if right := node.Right; right != nil {
            stack.push(right)
        }
        if left := node.Left; left != nil {
            stack.push(left)
        }

        elem = stack.pop()
    }
}

// 树的中序遍历（非递归实现）
func inOrder(node *BinaryTree) {
    stack := MyStack{List: list.New()}
    current := node

    for stack.List.Len() > 0 || current != nil {
        if current != nil {
            stack.push(current)
            current = current.Left
        } else {
            current = stack.pop().(*BinaryTree)
            fmt.Println(current.Value)
            current = current.Right
        }
    }
}

// 树的后序遍历（非递归实现）
func postOrder(node *BinaryTree) {
    stack1, stack2 := MyStack{List: list.New()}, MyStack{List: list.New()}
    stack1.push(node)

    for stack1.List.Len() > 0 {
        elem := stack1.pop().(*BinaryTree)
        stack2.push(elem)

        if elem.Left != nil {
            stack1.push(elem.Left)
        }

        if elem.Right != nil {
            stack1.push(elem.Right)
        }
    }

    for stack2.List.Len() > 0 {
        elem := stack2.pop().(*BinaryTree)
        fmt.Println(elem.Value)
    }
}

// 二叉树的层次遍历（非递归实现）
func levelOrder(node *BinaryTree) {
    var nlast *BinaryTree
    last := node
    level := 1
    queue := MyQueue{List: list.New()}  // queue存放当前层次的节点

    fmt.Println(fmt.Sprintf("-----this is %d level-----", level))
    queue.push(node)
    for queue.List.Len() > 0 {
        node := queue.pop().(*BinaryTree)  // 从队列中弹出一个节点进行访问

        if node.Left != nil {
            queue.push(node.Left)
            nlast = node.Left
        }

        if node.Right != nil {
            queue.push(node.Right)
            nlast = node.Right
        }

        fmt.Println(node.Value) // 打印当前正在遍历的节点
        if last == node && (node.Left != nil || node.Right != nil) {
            last = nlast
            level++
            fmt.Println()
            fmt.Println(fmt.Sprintf("-----this is %d level-----", level))
        }
    }
}

// 树的层次遍历
func levelTreeOrder(node *Tree) {
    var nlast *Tree
    last := node
    queue := MyQueue{List: list.New()}
    queue.push(node)

    for queue.List.Len() > 0 {
        node := queue.pop().(*Tree)

        for _, elem := range node.Children {
            queue.push(elem)
            nlast = elem
        }

        fmt.Println(node.Value)

        if last == node {
            last = nlast
            fmt.Println()
        }

    }
}

func preOrderToStr(node *BinaryTree) (ret string) {
    if node == nil {
        return "#!"
    }

    ret += fmt.Sprintf("%d!", node.Value)
    ret += preOrderToStr(node.Left)
    ret += preOrderToStr(node.Right)
    return ret
}

func strToBinaryTree(arr []string, index *int) *BinaryTree {
    if *index >= len(arr) {
        return nil
    }

    if arr[*index] == "#" {
        *index++
        return nil
    }

    node := &BinaryTree{}
    node.Value = arr[*index]
    *index++

    node.Left = strToBinaryTree(arr, index)
    node.Right = strToBinaryTree(arr, index)
    return node
}

func main() {
    str := "1!2!4!#!#!5!#!#!3!6!#!#!7!#!#!"
    arr := strings.Split(str, "!")
    index := 0
    node := strToBinaryTree(arr, &index)
    preOrder(node)

    // -----------多叉树层次遍历（打印行号）-----------------
    // node11 := &Tree{Value: 10}
    // node10 := &Tree{Value: 10}
    // node9 := &Tree{Value: 9}
    // node8 := &Tree{Value: 8}
    // node7 := &Tree{Value: 7}
    // node6 := &Tree{Value: 6}
    // node5 := &Tree{Value: 5, Children: []*Tree{node8, node9, node10, node11}}
    // node4 := &Tree{Value: 4}
    // node3 := &Tree{Value: 3}
    // node2 := &Tree{Value: 2, Children: []*Tree{node5, node6, node7}}
    // root := &Tree{Value: 1, Children: []*Tree{node2, node3, node4}}
    // levelTreeOrder(root)

    // -----------二叉树-----------------
    // node7 := &BinaryTree{Value: 7}
    // node6 := &BinaryTree{Value: 6}
    // node5 := &BinaryTree{Value: 5}
    // node4 := &BinaryTree{Value: 4}
    // node3 := &BinaryTree{Value: 3, Left: node6, Right: node7}
    // node2 := &BinaryTree{Value: 2, Left: node4, Right: node5}
    // root := &BinaryTree{Value: 1, Left: node2, Right: node3}

    // fmt.Println(preOrderToStr(root))

    // ---------二叉树前序遍历（递归、非递归）
    // preOrderRecur(root)
    // fmt.Println()
    // preOrder(root)

    // ---------二叉树中序遍历（递归、非递归）
    // inOrderRecu(root)
    // fmt.Println()
    // inOrder(root)

    // ---------二叉树后序遍历（递归、非递归）
    // posOrderRecu(root)
    // fmt.Println()
    // postOrder(root)

    // ---------二叉树层次遍历（打印行号）
    // levelOrder(root)

    // queue := MyQueue{List: list.New()}
    // queue.push(1)
    // queue.push(2)
    // queue.push(3)
    // fmt.Println(queue.pop())
    // fmt.Println(queue.pop())
    // fmt.Println(queue.pop())
    // fmt.Println(queue.pop())

    // fmt.Println(root)
    // stack := MyStack{List: list.New()}
    // stack.push(1)
    // stack.push(2)
    // stack.push(3)

    // fmt.Println(stack.pop())
    // fmt.Println(stack.pop())
    // fmt.Println(stack.pop())
    // fmt.Println(stack.pop())
}