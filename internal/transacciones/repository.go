package transacciones

import (
	"fmt"
	"github.com/SebasPu/Transacciones/pkg/store"
)

type transaccion struct {
	Id                int     `json:"id"`
	CodigoTransaccion *string `json:"codigo" binding:"required"`
	Moneda            *string `json:"moneda" binding:"required"`
	Monto             *string `json:"monto" binding:"required"`
	Emisor            *string `json:"emisor" binding:"required"`
	Receptor          *string `json:"receptor" binding:"required"`
	Fecha             *string `json:"fecha" binding:"required"`
}

var listaTransacciones []transaccion
var lastID int

type Repository interface {
	LastId() (int, error)
	GetAll() ([]transaccion, error)
	Store(tran transaccion) (transaccion, error)
	Update(tran transaccion) (transaccion, error)
	UpdateCod(id int, codigo, monto string) (transaccion, error)
	Delete(id int) error
}

type repository struct {
	db store.Store
}

func (t *repository) LastId() (int, error) {
	return lastID, nil
}

func (t *repository) GetAll() ([]transaccion, error) {
	t.db.Read(&listaTransacciones)
	return listaTransacciones, nil
}

func (t *repository) Store(tran transaccion) (transaccion, error) {
	t.db.Read(&listaTransacciones)
	listaTransacciones = append(listaTransacciones, tran)
	lastID = tran.Id
	if err := t.db.Write(listaTransacciones); err != nil {
		return transaccion{}, err
	}
	return tran, nil
}

func (t *repository) Update(tran transaccion) (transaccion, error) {
	t.db.Read(&listaTransacciones)
	update := false
	for i := range listaTransacciones {
		if listaTransacciones[i].Id == tran.Id {
			listaTransacciones[i] = tran
			update = true
		}
	}
	if err := t.db.Write(listaTransacciones); err != nil || !update {
		return transaccion{}, fmt.Errorf("transaccion %d no encontrada", tran.Id)
	}
	return tran, nil
}

func (t *repository) UpdateCod(id int, codigo, monto string) (transaccion, error) {
	var tran transaccion
	update := false
	for i := range listaTransacciones {
		if listaTransacciones[i].Id == id {
			listaTransacciones[i].CodigoTransaccion = &codigo
			listaTransacciones[i].Monto = &monto
			update = true
			tran = listaTransacciones[i]
		}
	}
	if err := t.db.Write(listaTransacciones); err != nil || !update {
		return transaccion{}, fmt.Errorf("transaccion %d no encontrada", tran.Id)
	}
	return tran, nil
}

func (t *repository) Delete(id int) error {
	delete := false
	var index int
	for i := range listaTransacciones {
		if listaTransacciones[i].Id == id {
			index = i
			delete = true
		}
	}
	listaTransacciones = append(listaTransacciones[:index], listaTransacciones[index+1:]...)
	if err := t.db.Write(listaTransacciones); err != nil || !delete {
		return fmt.Errorf("transaccion %d no encontrada", id)
	}
	return nil
}

func NewRepository(db store.Store) Repository {
	return &repository{
		db: db,
	}
}
