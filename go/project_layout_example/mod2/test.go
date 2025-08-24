package main

import (
	"fmt"
	"github.com/google/uuid"
	mod1_pkg1 "mod1/pkg1"
	mod2_pkg1 "mod2/pkg1"
)

func main() {
	a := mod1_pkg1.Add(1, 2)
	m := mod2_pkg1.Multi(1, 2)

	// 使用GitHub上的远程模块生成UUID
	uid := uuid.New()
	fmt.Println("Add result:", a)
	fmt.Println("Multi result:", m)
	fmt.Println("Generated UUID:", uid)
}
