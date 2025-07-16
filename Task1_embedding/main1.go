package main

import "fmt"

// родительская структура
type Human struct {
	Name   string
	Age    int
	Weight float64
	Height float64
}

//методы родительской структуры
func (s Human) Hello() string {
	return fmt.Sprintf("Привет, моё имя %s\n", s.Name)
}

func (p Human) Info() string {
	return fmt.Sprintf(
		"Возраст: %d\nВес: %.2f кг\nРост: %.2f м\n",
		p.Age, p.Weight, p.Height,
	)
}

// 'наследование' реализовано через встраивание (композицию); Human без имени поля, значит, методы и поля доступны напрямую из Action
type Action struct {
	DaySteps int
	Human
}

func main() {

	steps := Action{DaySteps: 10000, Human: Human{Name: "Киса", Age: 33, Weight: 66.6, Height: 1.77}}

	//Action имеет доступ ко всем методам Human
	fmt.Println(steps.Hello())
	fmt.Println("Параметры перед началoм тренировки:")
	fmt.Println(steps.Info())
	fmt.Println("Сегодня я прошла:", steps.DaySteps, "шагов")
}

/*
Встраивание структур

Дана структура Human (с произвольным набором полей и методов).
Реализовать встраивание методов в структуре Action от родительской структуры Human (аналог наследования).
Подсказка: используйте композицию (embedded struct), чтобы Action имел все методы Human.
*/
