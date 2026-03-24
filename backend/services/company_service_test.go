package services

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateCompanyService_Success(t *testing.T) {
	setupServiceTestDB(t)

	c, err := CreateCompany("TestCo", "Addr")
	require.NoError(t, err)
	require.NotNil(t, c)
	require.NotZero(t, c.ID)
	require.Equal(t, "TestCo", c.Name)
}

func TestCreateCompanyService_InvalidName(t *testing.T) {
	setupServiceTestDB(t)

	c, err := CreateCompany("", "Addr")
	require.Error(t, err)
	require.Nil(t, c)
}

func TestListCompaniesService_Success(t *testing.T) {
	setupServiceTestDB(t)

	_, err := CreateCompany("C1", "")
	require.NoError(t, err)
	_, err = CreateCompany("C2", "")
	require.NoError(t, err)

	list, err := ListCompanies()
	require.NoError(t, err)
	require.Len(t, list, 2)
	names := []string{list[0].Name, list[1].Name}
	require.Contains(t, names, "C1")
	require.Contains(t, names, "C2")
}
