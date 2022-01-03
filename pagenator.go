package ghost

import "gorm.io/gorm"

const DEFAULT_COUNT_PER_PAGE = 20

type Paginator struct {
	curPage        int
	pageSize       int
	totalItemCount int
	totalPage      int
}

func (this *Paginator) getMaxPage() int {
	if this.totalPage != 0 {
		return this.totalPage
	}
	totalPage := this.totalItemCount / this.pageSize
	if this.totalItemCount%this.pageSize == 0 {
		if totalPage == 0 {
			totalPage = 1
		}
	} else {
		totalPage += 1
	}
	this.totalPage = totalPage
	return totalPage
}

func (this *Paginator) getPageRange() (start, end int) {
	start = (this.curPage - 1) * this.pageSize
	end = start + this.pageSize
	return
}

func (this *Paginator) Paginate(db *gorm.DB) *gorm.DB {
	var c int64
	db.Count(&c)
	this.totalItemCount = int(c)
	this.totalPage = this.getMaxPage()
	// 如果浏览页数超过最大页数，则显示最后一页数据
	if this.curPage > this.totalPage {
		this.curPage = this.totalPage
	}
	start, _ := this.getPageRange()
	db = db.Limit(this.pageSize).Offset(start)
	return db
}

func (this *Paginator) MockPaginate(records []interface{}) []interface{} {
	this.totalItemCount = len(records)
	this.totalPage = this.getMaxPage()
	// 如果浏览页数超过最大页数，则显示最后一页数据
	if this.curPage > this.totalPage {
		this.curPage = this.totalPage
	}
	start, end := this.getPageRange()
	return records[start:end]
}

func (this *Paginator) ToMap() map[string]int {
	return map[string]int{
		"cur_page":    this.curPage,
		"page_size":   this.pageSize,
		"total_count": this.totalItemCount,
		"max_page":    this.totalPage,
	}
}

func NewPaginator(curPage int, args ...int) *Paginator {
	countPerPage := DEFAULT_COUNT_PER_PAGE
	switch len(args) {
	case 1:
		if args[0] != 0 {
			countPerPage = args[0]
		}
	}
	inst := new(Paginator)
	inst.curPage = curPage
	inst.pageSize = countPerPage
	return inst
}
