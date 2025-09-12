package usecase

import "Test/internal/feature/table/repository"

type TableUseCaseInteractor struct {
	repo repository.TableRepository
}

func NewTableUseCaseInteractor(repo repository.TableRepository) *TableUseCaseInteractor {
	return &TableUseCaseInteractor{repo: repo}
}

//func (uc *TableUseCaseInteractor)
