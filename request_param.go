package ghost

import "mime/multipart"

type RequestParams struct{
	GMap
}

func (rp RequestParams) GetFile(key string) *multipart.FileHeader{
	if v, ok := rp.GMap[key]; ok && v != nil{
		return v.(*multipart.FileHeader)
	}else{
		return nil
	}
}

func (rp RequestParams) GetFiles(key string) []*multipart.FileHeader{
	listKey := key + "[]"
	if v, ok := rp.GMap[listKey]; ok && v != nil{
		return v.([]*multipart.FileHeader)
	}else{
		return []*multipart.FileHeader{}
	}
}