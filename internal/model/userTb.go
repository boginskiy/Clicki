package model

import "time"

type UserTb struct {
	ID          int       `json:"id"`              // Уникальный идентификатор пользователя
	Name        string    `json:"name"`            // Имя пользователя (может содержать ФИО)
	Email       string    `json:"email"`           // Электронная почта (обычно служит уникальным идентификатором)
	Password    string    `json:"-"`               // Хешированный пароль (не отображается при сериализации JSON)
	CreatedAt   time.Time `json:"created_at"`      // Дата регистрации пользователя
	UpdatedAt   time.Time `json:"updated_at"`      // Последнее обновление профиля
	LastLoginAt time.Time `json:"last_login_at"`   // Время последнего входа
	IsActive    bool      `json:"is_active"`       // Флаг активности пользователя (заблокирован или активен)
	Roles       []string  `json:"roles,omitempty"` // Роли пользователя (администратор, модератор, обычный пользователь)
}
