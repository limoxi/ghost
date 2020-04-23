package utils

type lister struct {
	list []interface{}
	mset map[interface{}]struct{}
}

func (this *lister) initMSet(){
	if this.list == nil{
		return
	}
	this.mset = make(map[interface{}]struct{})
	for _, a := range this.list{
		this.mset[a] = struct{}{}
	}
}

func (this *lister) Has(el interface{}) bool{
	if len(this.mset) > 0{
		_, ok := this.mset[el]
		return ok
	}
	return false
}

func NewLister(list []interface{}) *lister{
	inst := new(lister)
	inst.list = list
	inst.initMSet()
	return inst
}

func NewListerFromInts(list []int) *lister{
	tl := make([]interface{}, 0, len(list))
	tmset := make(map[interface{}]struct{})
	for _, l := range list{
		tl = append(tl, l)
		tmset[l] = struct{}{}
	}
	inst := NewLister(tl)
	inst.mset = tmset
	return inst
}

func NewListerFromStrings(list []string) *lister{
	tl := make([]interface{}, 0, len(list))
	tmset := make(map[interface{}]struct{})
	for _, l := range list{
		tl = append(tl, l)
		tmset[l] = struct{}{}
	}
	inst := NewLister(tl)
	inst.mset = tmset
	return inst
}