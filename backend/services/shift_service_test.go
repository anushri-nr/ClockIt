package services

import (
	"errors"
	"strings"
	"testing"
	"time"

	"clockit/backend/database"
	"clockit/backend/models"

	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestReleaseShiftForWorker_Success(t *testing.T) {
	setupServiceTestDB(t)

	comp := models.Company{Name: "RelCo"}
	require.NoError(t, database.DB.Create(&comp).Error)

	sup := models.Employee{Name: "Sup", Email: "sup@rel.test", Password: "x", Role: models.RoleSupervisor, CompanyID: comp.ID}
	require.NoError(t, database.DB.Create(&sup).Error)

	worker := models.Employee{Name: "W", Email: "w@rel.test", Password: "x", Role: models.RoleWorker, CompanyID: comp.ID}
	require.NoError(t, database.DB.Create(&worker).Error)

	start := time.Now().Add(24 * time.Hour)
	end := start.Add(8 * time.Hour)
	sft := models.Shift{StartTime: start, EndTime: end, CreatedBy: sup.ID}
	require.NoError(t, database.DB.Create(&sft).Error)

	sa := models.ShiftAssignment{ShiftID: sft.ID, EmployeeID: worker.ID, AssigneeID: sup.ID, AssignedAt: time.Now(), Status: models.StatusAssigned}
	require.NoError(t, database.DB.Create(&sa).Error)

	out, err := ReleaseShiftForWorker(sft.ID, worker.ID)
	require.NoError(t, err)
	require.NotNil(t, out)
	require.Equal(t, models.StatusReleased, out.Status)

	var persisted models.ShiftAssignment
	require.NoError(t, database.DB.First(&persisted, sa.ID).Error)
	require.Equal(t, models.StatusReleased, persisted.Status)
}

func TestReleaseShiftForWorker_NotFound(t *testing.T) {
	setupServiceTestDB(t)

	comp := models.Company{Name: "RelCo2"}
	require.NoError(t, database.DB.Create(&comp).Error)

	worker := models.Employee{Name: "W2", Email: "w2@rel.test", Password: "x", Role: models.RoleWorker, CompanyID: comp.ID}
	require.NoError(t, database.DB.Create(&worker).Error)

	_, err := ReleaseShiftForWorker(9999, worker.ID)
	require.Error(t, err)
	require.ErrorIs(t, err, gorm.ErrRecordNotFound)
}

func TestReleaseShiftForWorker_InvalidStatus(t *testing.T) {
	setupServiceTestDB(t)

	comp := models.Company{Name: "RelCo3"}
	require.NoError(t, database.DB.Create(&comp).Error)

	sup := models.Employee{Name: "Sup2", Email: "sup2@rel.test", Password: "x", Role: models.RoleSupervisor, CompanyID: comp.ID}
	require.NoError(t, database.DB.Create(&sup).Error)

	worker := models.Employee{Name: "W2", Email: "w2@rel.test", Password: "x", Role: models.RoleWorker, CompanyID: comp.ID}
	require.NoError(t, database.DB.Create(&worker).Error)

	start := time.Now().Add(24 * time.Hour)
	end := start.Add(8 * time.Hour)
	sft := models.Shift{StartTime: start, EndTime: end, CreatedBy: sup.ID}
	require.NoError(t, database.DB.Create(&sft).Error)

	// create an assignment but mark as Requested
	sa := models.ShiftAssignment{ShiftID: sft.ID, EmployeeID: worker.ID, AssigneeID: sup.ID, AssignedAt: time.Now(), Status: models.StatusRequested}
	require.NoError(t, database.DB.Create(&sa).Error)

	out, err := ReleaseShiftForWorker(sft.ID, worker.ID)
	require.Error(t, err)
	require.Nil(t, out)
	require.ErrorIs(t, err, ErrAssignmentNotAssigned)
}

func createTestCompany(t *testing.T, name string) models.Company {
	t.Helper()
	c := models.Company{Name: name}
	require.NoError(t, database.DB.Create(&c).Error)
	return c
}

func createTestEmployee(
	t *testing.T,
	name, email string,
	role models.EmployeeRole,
	companyID uint,
) models.Employee {
	t.Helper()
	e := models.Employee{
		Name:      name,
		Email:     email,
		Password:  "x",
		Role:      role,
		CompanyID: companyID,
		Wage:      20,
	}
	require.NoError(t, database.DB.Create(&e).Error)
	return e
}

func createTestShift(t *testing.T, createdBy uint, start, end time.Time) models.Shift {
	t.Helper()
	s := models.Shift{
		StartTime: start,
		EndTime:   end,
		CreatedBy: createdBy,
	}
	require.NoError(t, database.DB.Create(&s).Error)
	return s
}

func TestCreateShift_Success(t *testing.T) {
	setupServiceTestDB(t)

	co := createTestCompany(t, "CreateShiftCo")
	supervisor := createTestEmployee(t, "Sup", "sup-create@test.local", models.RoleSupervisor, co.ID)

	start := time.Now().Add(24 * time.Hour).Truncate(time.Second)
	end := start.Add(8 * time.Hour)

	shift, err := CreateShift(start, end, supervisor.ID)
	require.NoError(t, err)
	require.NotNil(t, shift)
	require.NotZero(t, shift.ID)
	require.Equal(t, supervisor.ID, shift.CreatedBy)
	require.Equal(t, start.UTC().Unix(), shift.StartTime.UTC().Unix())
	require.Equal(t, end.UTC().Unix(), shift.EndTime.UTC().Unix())
}

func TestCreateShift_SupervisorNotFound(t *testing.T) {
	setupServiceTestDB(t)

	start := time.Now().Add(24 * time.Hour)
	end := start.Add(8 * time.Hour)

	shift, err := CreateShift(start, end, 999999)
	require.Error(t, err)
	require.Nil(t, shift)
	require.Contains(t, strings.ToLower(err.Error()), "supervisor")
	require.Contains(t, strings.ToLower(err.Error()), "not found")
}

func TestCreateShift_InvalidTimeOrder(t *testing.T) {
	setupServiceTestDB(t)

	co := createTestCompany(t, "InvalidOrderCo")
	supervisor := createTestEmployee(t, "Sup2", "sup-order@test.local", models.RoleSupervisor, co.ID)

	start := time.Now().Add(24 * time.Hour)
	end := start.Add(-1 * time.Hour) // end before start

	shift, err := CreateShift(start, end, supervisor.ID)
	require.Error(t, err)
	require.Nil(t, shift)
}

func TestCreateShift_SameStartAndEndTime(t *testing.T) {
	setupServiceTestDB(t)

	co := createTestCompany(t, "SameTimeCo")
	supervisor := createTestEmployee(t, "Sup3", "sup-same@test.local", models.RoleSupervisor, co.ID)

	start := time.Now().Add(24 * time.Hour).Truncate(time.Second)
	end := start

	shift, err := CreateShift(start, end, supervisor.ID)
	require.Error(t, err)
	require.Nil(t, shift)
}

func TestAssignWorkerToShift_Success(t *testing.T) {
	setupServiceTestDB(t)

	co := createTestCompany(t, "AssignSuccessCo")
	supervisor := createTestEmployee(t, "Sup4", "sup-assign@test.local", models.RoleSupervisor, co.ID)
	worker := createTestEmployee(t, "Worker1", "worker-assign@test.local", models.RoleWorker, co.ID)

	start := time.Now().Add(48 * time.Hour).Truncate(time.Second)
	end := start.Add(8 * time.Hour)
	shift := createTestShift(t, supervisor.ID, start, end)

	asg, err := AssignWorkerToShift(shift.ID, worker.ID, supervisor.ID)
	require.NoError(t, err)
	require.NotNil(t, asg)
	require.NotZero(t, asg.ID)
	require.Equal(t, shift.ID, asg.ShiftID)
	require.Equal(t, worker.ID, asg.EmployeeID)
	require.Equal(t, supervisor.ID, asg.AssigneeID)
	require.Equal(t, models.StatusAssigned, asg.Status)
}

func TestAssignWorkerToShift_ShiftNotFound(t *testing.T) {
	setupServiceTestDB(t)

	co := createTestCompany(t, "AssignNoShiftCo")
	supervisor := createTestEmployee(t, "Sup5", "sup-noshift@test.local", models.RoleSupervisor, co.ID)
	worker := createTestEmployee(t, "Worker2", "worker-noshift@test.local", models.RoleWorker, co.ID)

	asg, err := AssignWorkerToShift(999999, worker.ID, supervisor.ID)
	require.Error(t, err)
	require.Nil(t, asg)
}

func TestAssignWorkerToShift_WorkerNotFound(t *testing.T) {
	setupServiceTestDB(t)

	co := createTestCompany(t, "AssignNoWorkerCo")
	supervisor := createTestEmployee(t, "Sup6", "sup-noworker@test.local", models.RoleSupervisor, co.ID)

	start := time.Now().Add(48 * time.Hour).Truncate(time.Second)
	end := start.Add(8 * time.Hour)
	shift := createTestShift(t, supervisor.ID, start, end)

	asg, err := AssignWorkerToShift(shift.ID, 999999, supervisor.ID)
	require.Error(t, err)
	require.Nil(t, asg)
}

func TestAssignWorkerToShift_SupervisorNotAuthorized(t *testing.T) {
	setupServiceTestDB(t)

	co := createTestCompany(t, "AssignAuthCo")
	supervisor := createTestEmployee(t, "Sup7", "sup-auth@test.local", models.RoleSupervisor, co.ID)
	worker := createTestEmployee(t, "Worker3", "worker-auth@test.local", models.RoleWorker, co.ID)

	// non-supervisor tries to assign
	nonSupervisor := createTestEmployee(t, "WorkerAsAssigner", "assigner-worker@test.local", models.RoleWorker, co.ID)

	start := time.Now().Add(48 * time.Hour).Truncate(time.Second)
	end := start.Add(8 * time.Hour)
	shift := createTestShift(t, supervisor.ID, start, end)

	asg, err := AssignWorkerToShift(shift.ID, worker.ID, nonSupervisor.ID)
	require.Error(t, err)
	require.Nil(t, asg)
}

func TestAssignWorkerToShift_WorkerHasWrongRole(t *testing.T) {
	setupServiceTestDB(t)

	co := createTestCompany(t, "AssignWrongRoleCo")
	supervisor := createTestEmployee(t, "Sup8", "sup-wrongrole@test.local", models.RoleSupervisor, co.ID)

	// employee exists but is not worker
	employee := createTestEmployee(t, "AnotherSup", "another-sup@test.local", models.RoleSupervisor, co.ID)

	start := time.Now().Add(48 * time.Hour).Truncate(time.Second)
	end := start.Add(8 * time.Hour)
	shift := createTestShift(t, supervisor.ID, start, end)

	asg, err := AssignWorkerToShift(shift.ID, employee.ID, supervisor.ID)
	require.Error(t, err)
	require.Nil(t, asg)
}

func TestRejectShiftRequest_Success(t *testing.T) {
	setupServiceTestDB(t) // use your existing service DB helper

	company := models.Company{Name: "Reject Co"}
	require.NoError(t, database.DB.Create(&company).Error)

	supervisor := models.Employee{
		Name:      "Sup",
		Email:     "sup.reject@test.local",
		Role:      models.RoleSupervisor,
		CompanyID: company.ID,
		Wage:      40,
	}
	worker := models.Employee{
		Name:      "Worker",
		Email:     "worker.reject@test.local",
		Role:      models.RoleWorker,
		CompanyID: company.ID,
		Wage:      20,
	}
	require.NoError(t, database.DB.Create(&supervisor).Error)
	require.NoError(t, database.DB.Create(&worker).Error)

	now := time.Now().UTC()
	shift := models.Shift{
		StartTime: now.Add(2 * time.Hour),
		EndTime:   now.Add(10 * time.Hour),
		CreatedBy: supervisor.ID,
	}
	require.NoError(t, database.DB.Create(&shift).Error)

	assignment := models.ShiftAssignment{
		ShiftID:    shift.ID,
		EmployeeID: worker.ID,
		AssigneeID: supervisor.ID,
		Status:     models.StatusRequested,
		AssignedAt: now,
	}
	require.NoError(t, database.DB.Create(&assignment).Error)

	svc := SupervisorService{}
	out, err := svc.RejectShiftRequest(supervisor.ID, shift.ID)
	require.NoError(t, err)
	require.NotNil(t, out)
	require.Equal(t, models.StatusReleased, out.Status)
	require.Equal(t, supervisor.ID, out.AssigneeID)

	var dbRow models.ShiftAssignment
	require.NoError(t, database.DB.Where("shift_id = ?", shift.ID).First(&dbRow).Error)
	require.Equal(t, models.StatusReleased, dbRow.Status)
	require.Equal(t, supervisor.ID, dbRow.AssigneeID)
}

func TestRejectShiftRequest_NotFound_WrongCompany(t *testing.T) {
	setupServiceTestDB(t) // use your existing service DB helper

	companyA := models.Company{Name: "A"}
	companyB := models.Company{Name: "B"}
	require.NoError(t, database.DB.Create(&companyA).Error)
	require.NoError(t, database.DB.Create(&companyB).Error)

	supA := models.Employee{
		Name:      "SupA",
		Email:     "supa.reject@test.local",
		Role:      models.RoleSupervisor,
		CompanyID: companyA.ID,
		Wage:      40,
	}
	supB := models.Employee{
		Name:      "SupB",
		Email:     "supb.reject@test.local",
		Role:      models.RoleSupervisor,
		CompanyID: companyB.ID,
		Wage:      40,
	}
	workerB := models.Employee{
		Name:      "WorkerB",
		Email:     "workerb.reject@test.local",
		Role:      models.RoleWorker,
		CompanyID: companyB.ID,
		Wage:      20,
	}
	require.NoError(t, database.DB.Create(&supA).Error)
	require.NoError(t, database.DB.Create(&supB).Error)
	require.NoError(t, database.DB.Create(&workerB).Error)

	now := time.Now().UTC()
	shiftB := models.Shift{
		StartTime: now.Add(3 * time.Hour),
		EndTime:   now.Add(11 * time.Hour),
		CreatedBy: supB.ID,
	}
	require.NoError(t, database.DB.Create(&shiftB).Error)

	require.NoError(t, database.DB.Create(&models.ShiftAssignment{
		ShiftID:    shiftB.ID,
		EmployeeID: workerB.ID,
		AssigneeID: supB.ID,
		Status:     models.StatusRequested,
		AssignedAt: now,
	}).Error)

	svc := SupervisorService{}
	out, err := svc.RejectShiftRequest(supA.ID, shiftB.ID)
	require.Error(t, err)
	require.Nil(t, out)
	require.True(t, errors.Is(err, gorm.ErrRecordNotFound))
}
