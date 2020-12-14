package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"sync"
)



type TreeNode struct {
	val      string
	children []*TreeNode
}


func (t *TreeNode) SearchAndInsert(value string) (error, string) {
	if t == nil {
		return errors.New("tree is nil"), ""
	}

	folders := strings.Split(value, "/")

	if t.val == folders[1] {
		if folders[2] != "" {
			if t.children != nil {
				for _, tree := range t.children {
					if tree.val == folders[2] {
						return tree.SearchAndInsert(strings.Replace(value, "/"+folders[1], "", 1))
					}
				}
			} else {
				t.children = append(t.children, &TreeNode{val: folders[2]})
				for _, child := range t.children {
					if child.val == folders[2] {
						return child.SearchAndInsert(strings.Replace(value, "/"+folders[1], "", 1))
					}
				}
			}

		} else {
			return nil, "/"+ folders[1]
		}

	} else {
		if t.children != nil {
			for _, child := range t.children {
				if child.val == folders[1] {
					return child.SearchAndInsert(strings.Replace(value, "/"+folders[1], "", 1))
				}
			}
			t.children = append(t.children, &TreeNode{val: folders[1]})
			for _, child := range t.children {
				if child.val == folders[1] {
					return child.SearchAndInsert(value)
				}
			}
		} else {
			t.children = append(t.children, &TreeNode{val: folders[1]})
			for _, child := range t.children {
				if child.val == folders[1] {
					return child.SearchAndInsert(value)
				}
			}
		}

	}

	return nil, ""
}

func (t *TreeNode) countChildren(parentPath string) int {
	if t.children != nil {
		out := make(chan int)
		length := len(t.children)

		var sum int
		go func(out chan int) {
			wg := &sync.WaitGroup{}
			for _, value := range t.children {
				wg.Add(1)
				go func(out chan int, value *TreeNode, wg *sync.WaitGroup) {
					defer wg.Done()
					out <- value.countChildren(parentPath + "/" + t.val)
				}(out, value, wg)
			}
			wg.Wait()
			close(out)
		}(out)
		//for _, value := range t.children {
		//	go func(out chan int, value *TreeNode) {
		//		out <- value.countChildren(parentPath+"/"+t.val)
		//	}(out, value)
		//	//sum = sum + value.countChildren(parentPath+"/"+t.val)
		//	//sum = sum + <-out
		//}

		//for i:=0; i<length; i++ {
		//	sum = sum + <-out
		//}
		//close(out)
		for {
			val, ok := <-out
			if ok == false {
				break // exit break loop
			} else {
				sum = sum + val
			}
		}
		fmt.Printf("%s - %d \n", parentPath+"/"+t.val, sum+length)
		return sum+length
	}

	fmt.Printf("%s - %d \n", parentPath+"/"+t.val, 0)

	return 0
}

func main() {
	var fileLocation string
	path := os.Args
	if len(path) == 1 {
		fileLocation = "locs.txt"
	} else {
		fileLocation = os.Args[1]
	}
	configFile, err := ioutil.ReadFile(fileLocation)
	if err != nil {
		log.Fatal(err)
	}

	configLines := strings.Split(string(configFile), "\n")

	var treeNode TreeNode
	treeNode.val = "root"

	for _, line := range configLines {
		err := lineFilter(line)
		if err != nil {
			log.Println(err)
			continue
		}
		err, _ = treeNode.SearchAndInsert(line)
		if err != nil {
			return
		}

	}

	treeNode.countChildren("")
	fmt.Println("end")

}

func lineFilter(line string) error  {

	folders := strings.Split(line, "/")
	if len(folders) == 1 || len(folders) == 2 {
		return errors.New(fmt.Sprintf("path %s is incorrect (empty)", line))
	}

	if folders[0] != "" || folders[len(folders)-1] != "" {
		return errors.New(fmt.Sprintf("path %s is incorrect", line))
	}

	return nil
}
