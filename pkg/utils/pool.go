package utils

// Resetable — интерфейс, требующий наличия метода Reset()
type Resetable interface {
	Reset()
}

// Pool — собственная обобщенная структура пула, использующая канал.
// T должен реализовывать интерфейс Resetable.
type Pool[T Resetable] struct {
	// Канал используется для хранения и передачи объектов между Get() и Put()
	poolChan chan T
	// Функция-фабрика для создания новых объектов, если пул пуст
	newFunc func() T
}

// NewPool — функция-конструктор для создания пула.
// capacity определяет максимальное количество объектов, которые могут храниться в пуле.
func NewPool[T Resetable](capacity int, newFunc func() T) *Pool[T] {
	if capacity <= 0 {
		capacity = 10 // Устанавливаем разумное значение по умолчанию
	}
	p := &Pool[T]{
		// Буферизованный канал позволяет хранить объекты без блокировки отправителя (Put),
		// пока не будет достигнут лимит capacity.
		poolChan: make(chan T, capacity),
		newFunc:  newFunc,
	}
	return p
}

// Get извлекает объект из пула.
// Если в канале poolChan есть свободный объект, он возвращается.
// Если канал пуст, вызывается функция-фабрика newFunc для создания нового объекта.
func (p *Pool[T]) Get() T {
	select {
	case obj := <-p.poolChan:
		// Объект взят из пула
		return obj
	default:
		// Пул пуст, создаем новый объект с помощью фабрики
		return p.newFunc()
	}
}

// Put помещает объект обратно в пул.
// Перед возвратом объекта вызывается его метод Reset() для сброса состояния.
// Если канал полон (достигнут лимит capacity), объект просто игнорируется
// и сборщик мусора его утилизирует.
func (p *Pool[T]) Put(obj T) error {
	if any(obj) == nil {
		return nil
	}

	// САМЫЙ ВАЖНЫЙ ШАГ: Сброс состояния перед возвратом в пул.
	obj.Reset()

	select {
	case p.poolChan <- obj:
		return nil
	default:
		return NewPoolOvercrowded(cap(p.poolChan))
	}
}
