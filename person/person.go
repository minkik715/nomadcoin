package Person

import "fmt"

type Person struct {
	name string
	age  int
}

func (p *Person) Constructor(name string, age int) Person {
	(*p).name = name
	(*p).age = age
	return *p
}

func (p Person) SayHello() {
	fmt.Println("hello", p.name, p.age)
}
