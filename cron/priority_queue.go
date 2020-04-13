package cron

type Comparable interface {
	GetCompareValue() int64
}

type PriorityQueue []Comparable

func (p PriorityQueue) Len() int {
	return len(p)
}

func (p PriorityQueue) Less(i, j int) bool {
	return p[i].GetCompareValue() < p[j].GetCompareValue()
}

func (p PriorityQueue) Swap(i, j int) {
	p[j], p[i] = p[i], p[j]
}

func (p *PriorityQueue) Pop() interface{} {
	len := len(*p)
	ele := (*p)[len-1]
	*p = (*p)[:len-1]
	return ele
}

func (p *PriorityQueue) Push(ele interface{}) {
	*p = append(*p, ele.(Comparable))
}
