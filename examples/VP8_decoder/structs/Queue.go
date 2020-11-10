package structs

import (
	"container/list"
	"sync"
)

type Queue struct {
	sync.Mutex
	Data *list.List
	Alive bool
	maxLen int
	Name string
}
func (c *Queue) IsEmpty() bool{
	c.Lock()
	defer c.Unlock()
	if c.Data!=nil && c.Data.Len()>0 {
		return false
	}
	return true
}
func (c *Queue) Create(maxLen int){
	c.Lock()
	defer c.Unlock()
	c.Data=list.New()
	c.Alive=true
	c.maxLen=maxLen
}
func (c *Queue) Push(val interface{}) bool{
	c.Lock()
	defer c.Unlock()
	if c.Data!=nil && c.maxLen>c.Data.Len(){
		if  val!=nil{
			c.Data.PushFront(val)
			return true
		}
	}
	return false
}
func (c *Queue) PushRotate(val interface{}) bool{
	c.Lock()
	defer c.Unlock()
	if c.Data!=nil{
		if  c.maxLen>c.Data.Len(){
			el:=c.Data.Back()
			if el!=nil {
				c.Data.Remove(el)
			}
		}
		if  val!=nil{
			c.Data.PushFront(val)
			return true
		}
	}
	return false
}
func (c *Queue) Pull() interface{}{
	c.Lock()
	defer c.Unlock()
	if c.Data!=nil &&c.Data.Len() > 0 {
		el:=c.Data.Back()
		fs := el.Value
		c.Data.Remove(el)
		return fs
	}
	return nil
}
func (c *Queue) GetLength() int {
	c.Lock()
	defer c.Unlock()
	if c.Data!=nil{
		return c.Data.Len()
	}
	return 0
}
func (c *Queue) Kill() {
	c.Lock()
	defer c.Unlock()
	if c.Alive==true{
		c.Data=nil
	}
	c.Alive=false
}
func (c *Queue) IsAlive() bool{
	c.Lock()
	defer c.Unlock()
	return c.Alive
}
