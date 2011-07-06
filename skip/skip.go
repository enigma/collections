package skip

import (
	. "github.com/badgerodon/collections"
	"fmt"
	"rand"
	"time"
)

type (
	node struct {
		neighbors []*node
		key Any
		value Any
	}
	SkipList struct {
		root *node
		size int
		less func(Any,Any)bool
		gen *rand.Rand
		probability float64
	}
)
// Create a new skip list
func New(less func(Any,Any)bool) *SkipList {
	gen := rand.New(rand.NewSource(time.Nanoseconds()))
	n := &node{make([]*node, 0),nil,nil}
	return &SkipList{n, 0, less, gen, 0.75}
}
func (this *SkipList) Do(f func(Any)bool) {
	if this.size == 0 {
		return
	}
	cur := this.root.neighbors[0]
	for cur != nil {
		if !f(cur.value) {
			break
		}
		cur = cur.neighbors[0]
	}
}
// Get an item from the skip list
func (this *SkipList) Get(key Any) Any {
	if this.size == 0 {
		return nil
	}
	
	cur := this.root
	for i := len(cur.neighbors)-1; i >= 0; i-- {
		for this.less(cur.neighbors[i].key, key) {
			cur = cur.neighbors[i]
		}
	}
	
	if this.equals(cur.key, key) {
		return cur.value
	}
	
	return nil
}
// Insert a new item into the skip list
func (this *SkipList) Insert(key Any, value Any) {
	h := len(this.root.neighbors)
	
	updates := this.getUpdateTable(key)
	
	nh := this.pickHeight()
	n := &node{make([]*node, nh),key,value}
	
	if nh > h {
		this.root.neighbors = append(this.root.neighbors, n)
	}
	
	for i := 0; i < nh; i++ {
		if i < h {
			n.neighbors[i] = updates[i].neighbors[i]
			updates[i].neighbors[i] = n
		}
	}
	
	this.size++
}
// Remove an item from the skip list
func (this *SkipList) Remove(key Any) Any {
	updates := this.getUpdateTable(key)
	cur := updates[0].neighbors[0]
	
	// If we found it
	if cur != nil && this.equals(key, cur.key) {
		// Change all the linked lists
		for i := 0; i < len(updates); i++ {
			if this.equals(updates[i].neighbors[i].key, cur.key) {
				updates[i].neighbors[i] = cur.neighbors[i]
			} else {
				break
			}
		}
		
		// Kill off the upper links if they're nil
		for i := len(this.root.neighbors)-1; i>=0; i-- {
			if this.root.neighbors[i] == nil {
				this.root.neighbors = this.root.neighbors[:i]
			} else {
				break
			}
		}
		
		this.size--
		
		return cur.value
	}
	
	return nil
}
// String representation of the list
func (this *SkipList) String() string {	
	str := "{"
	if len(this.root.neighbors) > 0 {
		cur := this.root.neighbors[0]
		for cur != nil {
			str += fmt.Sprint(cur.key)
			str += ":"
			str += fmt.Sprint(cur.value)
			str += " "
			cur = cur.neighbors[0]
		}
	}
	str += "}"
	
	return str
}
// Get a vertical list of nodes of all the things that occur
//  immediately before "key"
func (this *SkipList) getUpdateTable(key Any) []*node {
	cur := this.root
	h := len(cur.neighbors)
	updates := make([]*node, h)
	for i := h-1; i >= 0; i-- {
		for cur.neighbors[i] != nil && this.less(cur.neighbors[i].key, key) {
			cur = cur.neighbors[i]
		}
		updates[i] = cur
	}
	return updates
}
// Defines an equals method in terms of "less"
func (this *SkipList) equals(a, b Any) bool {
	return !this.less(a,b) && !this.less(b,a)
}
// Pick a random height
func (this *SkipList) pickHeight() int {
	h := 1
	for this.gen.Float64() > this.probability {
		h++
	}
	if h > len(this.root.neighbors) {
		return h + 1
	}
	return h
}
