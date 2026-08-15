package todo

import "strings"

type Service struct {
	repository *Repository
}

func NewService(repository *Repository) *Service {
	return &Service{repository: repository}
}

func (s *Service) List() ([]Todo, error) {
	return s.repository.List()
}

func (s *Service) Get(id int) (Todo, error) {
	return s.repository.Get(id)
}

func (s *Service) Create(input Input) (Todo, error) {
	input, err := validate(input)
	if err != nil {
		return Todo{}, err
	}

	return s.repository.Create(input)
}

func (s *Service) Update(id int, input Input) (Todo, error) {
	input, err := validate(input)
	if err != nil {
		return Todo{}, err
	}

	return s.repository.Update(id, input)
}

func (s *Service) Delete(id int) error {
	return s.repository.Delete(id)
}

func validate(input Input) (Input, error) {
	input.Body = strings.TrimSpace(input.Body)
	if input.Body == "" {
		return Input{}, ErrInvalidBody
	}

	return input, nil
}
