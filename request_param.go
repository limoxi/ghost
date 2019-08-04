package ghost

import "mime/multipart"

type RequestParam Map

func (rp RequestParam) GetFile(key string) *multipart.FileHeader{
	if v, ok := rp[key]; ok && v != nil{
		return v.(*multipart.FileHeader)
	}else{
		return nil
	}
}

func (rp RequestParam) GetFiles(key string) []*multipart.FileHeader{
	listKey := key + "[]"
	if v, ok := rp[listKey]; ok && v != nil{
		return v.([]*multipart.FileHeader)
	}else{
		return []*multipart.FileHeader{}
	}
}