package services

import (
	"github.com/kitabisa/loan-engine/internal/repositories"
	"github.com/kitabisa/loan-engine/pkg/external"
)

type ServiceFactory struct {
	RepoFactory    *repositories.RepositoryFactory
	EmailService   external.EmailService
	StorageService external.StorageService
	JwtSecret      string
}

func NewServiceFactory(
	repoFactory *repositories.RepositoryFactory,
	emailService external.EmailService,
	storageService external.StorageService,
	jwtSecret string,
) *ServiceFactory {
	return &ServiceFactory{
		RepoFactory:    repoFactory,
		EmailService:   emailService,
		StorageService: storageService,
		JwtSecret:      jwtSecret,
	}
}

func (f *ServiceFactory) BorrowerService() BorrowerService {
	return NewBorrowerService(f.RepoFactory.BorrowerRepository())
}

func (f *ServiceFactory) LoanService() LoanService {
	return NewLoanService(
		f.RepoFactory.LoanRepository(),
		f.RepoFactory.LoanApprovalRepository(),
		f.RepoFactory.LoanDisbursementRepository(),
		f.RepoFactory.LoanInvestmentRepository(),
		f.RepoFactory.LoanStateHistoryRepository(),
		f.RepoFactory.InvestorRepository(),
		f.EmailService,
		f.StorageService,
	)
}

func (f *ServiceFactory) InvestorService() InvestorService {
	return NewInvestorService(f.RepoFactory.InvestorRepository())
}

func (f *ServiceFactory) AuthService() AuthService {
	return NewAuthService(f.RepoFactory.UserRepository(), f.JwtSecret)
}