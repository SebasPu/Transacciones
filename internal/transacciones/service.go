package transacciones

type Service interface {
	LastId() (int, error)
	GetAll() ([]transaccion, error)
	Store(codigoTransaccion, moneda, monto, emisor, receptor, fecha string) (transaccion, error)
	Update(id int, codigoTransaccion, moneda, monto, emisor, receptor, fecha string) (transaccion, error)
	UpdateCod(id int, codigo, monto string) (transaccion, error)
	Delete(id int) error
}

type service struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return &service{
		repository: repository,
	}
}

func (s *service) LastId() (int, error) {
	return s.repository.LastId()
}

func (s *service) GetAll() ([]transaccion, error) {
	return s.repository.GetAll()
}

func (s *service) Store(codigoTransaccion, moneda, monto, emisor, receptor, fecha string) (transaccion, error) {
	id, _ := s.repository.LastId()
	id++
	t := transaccion{id, &codigoTransaccion, &moneda, &monto, &emisor, &receptor, &fecha}
	return s.repository.Store(t)
}

func (s *service) Update(id int, codigoTransaccion, moneda, monto, emisor, receptor, fecha string) (transaccion, error) {
	t := transaccion{id, &codigoTransaccion, &moneda, &monto, &emisor, &receptor, &fecha}
	return s.repository.Update(t)
}

func (s *service) UpdateCod(id int, codigo, monto string) (transaccion, error){
	return s.repository.UpdateCod(id, codigo, monto)
}

func (s *service) Delete(id int) error {
	return s.repository.Delete(id)
}
