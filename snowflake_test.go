package snowflake

import (
	"fmt"
	"testing"
)

func TestNewNode(t *testing.T) {

	_, err := NewNode(0)
	if err != nil {
		t.Fatalf("error creating NewNode, %s", err)
	}

	_, err = NewNode(5000)
	if err == nil {
		t.Fatalf("no error creating NewNode, %s", err)
	}

}

func TestGenerateDuplicateID(t *testing.T) {

	node, _ := NewNode(1)

	var x, y ID
	for i := 0; i < 1000000; i++ {
		y = node.Generate()
		if x == y {
			t.Errorf("x(%d) & y(%d) are the same", x, y)
		}
		x = y
	}
}

func TestSnowflakeGenerate(t *testing.T) {
	node, _ := NewNode(1)

	ch := make(chan ID)
	count := 1000000
	open := make(chan struct{}, 1)

	for i := 0; i < count; i++ {
		go func(i int) {
			if i == count-1 {
				close(open)
			}
			<-open
			id := node.Generate()
			ch <- id
		}(i)
	}

	defer close(ch)

	m := make(map[ID]int)

	for i := 0; i < count; i++ {
		id := <-ch
		// 如果 map 中存在为 id 的 key, 说明生成的 snowflake ID 有重复
		_, ok := m[id]
		if ok {
			t.Error("ID is not unique!\n")
			return
		}
		// 将 id 作为 key 存入 map
		m[id] = i
	}
	fmt.Println("snowflake ID is success")
}

func TestGenerateId(t *testing.T) {
	id := GenerateId()
	if id == 0 {
		t.Errorf("id is zero")
	}
}
