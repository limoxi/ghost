package ghost

import "mime/multipart"

type RequestParams struct{
	Map
}

func (rp RequestParams) GetFile(key string) *multipart.FileHeader{
	if v, ok := rp.Map[key]; ok && v != nil{
		return v.(*multipart.FileHeader)
	}else{
		return nil
	}
}

func (rp RequestParams) GetFiles(key string) []*multipart.FileHeader{
	listKey := key + "[]"
	if v, ok := rp.Map[listKey]; ok && v != nil{
		return v.([]*multipart.FileHeader)
	}else{
		return []*multipart.FileHeader{}
	}
}