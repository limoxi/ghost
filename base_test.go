package ghost

import (
	"fmt"
	"testing"
)

type s struct{
	a string
}

func (this *s) cc () *s{
	fmt.Print(this, *this)
	ns := *this
	ns.a = "heheda"
	return &ns
}

func (this *s) c1(){
	fmt.Printf("inner======%p \n", this)
}

func (this *s) c2() {
	n := *this
	//n.a = "444"
	fmt.Printf("c2====%p \n", &n)
}


func TestBase (t *testing.T){

	s1 := s{a: "123"}
	s2 := s1
	s3 := s1

	s2.a = "adasd"

	fmt.Printf("%p, %p, %p \n", &s1, &s2, &s3)
	fmt.Println(s1, s2, s3)
	fmt.Println(&s1, &s2, &s3)
}