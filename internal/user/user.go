package user

import (
	"fmt"
	"math/rand"
	"time"

	mod "github.com/boginskiy/Clicki/internal/model"
)

type User struct {
}

func NewUser() *User {
	return &User{}
}

// generatorVasiliev - Супер функция для генерации Василиев
func (u *User) generatorVasiliev() string {
	return fmt.Sprintf("%s_%d", "Vasiliy", rand.Intn(1000))
}

func (u *User) CreateEmpty() *mod.UserTb {
	return &mod.UserTb{
		Name:      u.generatorVasiliev(),
		CreatedAt: time.Now(),
	}
}
