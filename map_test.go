package ghost

import (
	"fmt"
	"testing"
)

func TestMap(t *testing.T){

	data := map[string]interface{}{
		"int": 12,
		"float": 12.56,
		"string": "abc",
		"list": []int{1,2,3},
		"dict": map[string]string{
			"a": "aaa",
			"b": "bbb",
			"c": "ccc",
		},
	}
	m := NewMapFromData(data)

	t.Log(m.Get("list").([]int))
}

// 选择排序
// 每一轮，找到数组中最小的数，与当前索引位置数调换位置
func selectionSort(nums []int){
	for i, num := range nums{
		min := num
		minIndex := i
		for j:=i+1; j<len(nums);j++{
			if min > nums[j]{
				min = nums[j]
				minIndex = j
			}
		}
		nums[i], nums[minIndex] = nums[minIndex], nums[i]
	}
}

// 插入排序
// 将目标元素插入到恰当的位置前，后续所有的元素都要向右移动一位
func insertionSort(nums []int){
	for i:=1;i<len(nums);i++{
		for j:=i;j>0&&nums[j-1]>nums[j];j--{
			nums[j-1], nums[j] = nums[j], nums[j-1]
		}
	}
}

// 希尔排序
// 设置一个递增序列，使数组中任意间隔为h的元素都是有序的，即构造多个局部有序，
// 并在最终通过插入排序完成整个数组的排序
func shellSort(nums []int){
	l := len(nums)
	// 构造递增序列
	h := 1
	for ; h<l/3;{
		h = 3* h + 1
	}
	fmt.Println(h, "======")
	for ;h>=1;{
		for i:=h;i<l;i++{
			for j:=i;j>=h && nums[j]<nums[j-h];j-=h{
				nums[j], nums[j-h] = nums[j-h], nums[j]
			}
		}
		h = h/3
		fmt.Println(h, "--------", nums)
	}
}

// 原地归并排序
// 1、复制数组到另一个新数组中
// 将新数组分为两部分，并按序归并到原数组中
func mergeSort1(nums []int, start, mid, end int){
	tmp := make([]int, 0, len(nums))
	for _, num := range nums{
		tmp = append(tmp, num)
	}
	s := start
	e := mid + 1
	for k:=start;k<=end;k++{
		if s > mid{
			nums[k] = tmp[e]
			e+=1
		}else if e>end{
			nums[k] = tmp[s]
			s+=1
		}else if tmp[s]<tmp[e]{
			nums[k] = tmp[s]
			s+=1
		}else{
			nums[k] = tmp[e]
			e+=1
		}
	}
}

// 快速排序
// 随机选取一个基准数a切分数组为两个子数组，左边数组小于a，右边数组大于a
func quickSort(){

}

func TestSort(t *testing.T){
	ip := []int{4,2,6,12,35,5,11,34,16}
	mergeSort1(ip, 0, 4, 8)
	t.Log(ip)

}