package utils


func StringInList(str string, list []interface{}) bool{
	for _, s := range list{
		if s.(string) == str{
			return true
		}
	}
	return false
}