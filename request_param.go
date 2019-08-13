package ghost

import "mime/multipart"

type RequestParams Map

func (rp RequestParams) GetFile(key string) *multipart.FileHeader{
	if v, ok := rp[key]; ok && v != nil{
		return v.(*multipart.FileHeader)
	}else{
		return nil
	}
}

func (rp RequestParams) GetFiles(key string) []*multipart.FileHeader{
	listKey := key + "[]"
	if v, ok := rp[listKey]; ok && v != nil{
		return v.([]*multipart.FileHeader)
	}else{
		return []*multipart.FileHeader{}
	}
}