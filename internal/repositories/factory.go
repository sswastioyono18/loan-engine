package repositories

import ()

type RepositoryFactory struct {
	driver Driver
}

func NewRepositoryFactory(driver Driver) *RepositoryFactory {
	return &RepositoryFactory{
		driver: driver,
	}
}

func (f *RepositoryFactory) BorrowerRepository() BorrowerRepository {
	return NewBorrowerRepository(f.driver)
}

func (f *RepositoryFactory) LoanRepository() LoanRepository {
	return NewLoanRepository(f.driver)
}

func (f *RepositoryFactory) LoanApprovalRepository() LoanApprovalRepository {
	return NewLoanApprovalRepository(f.driver)
}

func (f *RepositoryFactory) LoanDisbursementRepository() LoanDisbursementRepository {
	return NewLoanDisbursementRepository(f.driver)
}

func (f *RepositoryFactory) InvestorRepository() InvestorRepository {
	return NewInvestorRepository(f.driver)
}

func (f *RepositoryFactory) LoanInvestmentRepository() LoanInvestmentRepository {
	return NewLoanInvestmentRepository(f.driver)
}

func (f *RepositoryFactory) LoanStateHistoryRepository() LoanStateHistoryRepository {
	return NewLoanStateHistoryRepository(f.driver)
}

func (f *RepositoryFactory) UserRepository() UserRepository {
	return NewUserRepository(f.driver)
}