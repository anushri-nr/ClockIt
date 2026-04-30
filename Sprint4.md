# Sprint 4 - ClockIt

## Overview
Sprint 4 focused on improving test coverage, documenting the current backend APIs, hardening the shift workflow behavior added in earlier sprints, and deploying the full stack to AWS. We also completed the Supervisor Dashboard frontend, integrating real-time state management and completing the full end-to-end shift lifecycle.
---

## Work Completed in Sprint 4

### Backend
- Expanded controller test coverage from `64.3%` to `83.3%`.
- Expanded services test coverage from `48.8%` to `55.3%`.
- Added coverage for company list DB failures.
- Added shift creation validation tests for invalid timestamps, unauthorized requests, and service-level failures.
- Added assignment coverage for unauthorized assignment requests.
- Added worker shift request coverage for success, unauthorized access, bad shift IDs, and not-found cases.
- Added supervisor requested-shift coverage for unauthorized access and supervisor-not-found cases.
- Added supervisor login/register negative-path coverage.
- Added worker login negative-path coverage.
- Added worker assigned-shift filtering coverage for time windows.
- Added worker availability success coverage.
- Added released-shift lookup not-found coverage.
- Added service coverage for unassigned shift filtering, latest assignment selection, assignment update behavior, overtime threshold exclusion, and company list edge cases.
- Added shift deletion endpoint wrapped in a database transaction to safely cascade and delete assignment rows.
- Added a company-filtered worker directory endpoint (`GetCompanyWorkers`).

### Frontend
- Fully implemented the main dashboard layout using Angular Material components.
- Added a `Today's Coverage` metric that dynamically counts workers scheduled for the current day.
- Added a toggle to switch the main schedule grid between a `Today` and `Week` view, alongside a real-time `Worker Search` filter.
- Built a read-only directory modal to view company employee contact info (Name, Email, Phone).
- Integrated modals for `Create Shift`, `Open Shifts` (with Delete functionality), and `Assign Worker`.
- Engineered the `Assign Worker` dropdown to automatically parse the selected shift's date and fetch only available workers.
- Refactored component architecture so that actions (Approve, Reject, Delete, Unassign) trigger instant data reloads for the Schedule Grid, Pending Approvals, and Overtime Alerts without requiring a page refresh.
- Added consistent pointer cursors, Material icons, and Angular Material `SnackBar` (Toast) notifications for all success and error states.
---

## AWS Deployment

### Backend
- Deployed Go backend to AWS Elastic Beanstalk.
- Backend live at: `http://clockit.us-east-1.elasticbeanstalk.com`

### Frontend
- Configured Angular production build to use `environment.prod.ts` for environment-specific API URLs.
- Deployed Angular app to AWS S3 static website hosting.
- Frontend live at: `http://clockit-app.s3-website-us-east-1.amazonaws.com`

---

## Frontend Unit Tests

### App
- `app.spec.ts`
  - `should create the app`

### Auth
- `auth.spec.ts`
  - `should be created`
- `auth.service.spec.ts`
  - `should be created and have a null initial user`
  - `should handle login and update local storage`
  - `should handle registration and auto-login`
  - `should fetch the list of companies`
- `auth.guard.spec.ts`
  - `should be created`
- `auth.interceptor.spec.ts`
  - `should be created`

### Supervisor Service
- `supervisor.service.spec.ts`
  - `should fetch shifts with an optional status filter`
  - `should assign a worker to a shift`
  - `should fetch available workers for a specific date`
  - `should be created successfully`
  - `should fetch assigned shifts for the schedule view`
  - `should send a PATCH request to reject a requested shift`
  - `should fetch workers at risk of overtime`
  - `should send a DELETE request to completely remove a shift`
  - `should send a PATCH request to unassign a worker from a shift`
  - `should fetch all workers for the supervisor company directory`

### Shift List Component
- `shift-list.component.spec.ts`
  - `should create`

---

## Cypress E2E Tests

### Supervisor Authentication & Dashboard Features
- `supervisor-login.cy.ts`
  - `should successfully log in a supervisor and load the dashboard`
- `supervisor-sprint4.cy.ts`
  - `should open the Worker Directory and display company workers`
  - `should unassign a worker from an active shift`
  - `should completely delete an unassigned shift`

### Worker Authentication
- `worker-login.cy.ts`
  - `should successfully log in a worker and load the dashboard`

### Worker Account Creation
- `worker-create_account.cy.ts`
  - `creates worker account using Panda Express`

### Worker Shift Management
- `worker-shift-management.cy.ts`
  - `should display RELEASE SHIFT button on the shift card`
  - `should click RELEASE SHIFT button and release the shift`
  - `should remove the shift from upcoming list after release`
  - `should display forward navigation button for weekly hours`
  - `should navigate to next week when forward button is clicked`
  - `should update weekly hours display when navigating forward`
  - `should allow multiple forward navigations`
  - `should display backward navigation button for weekly hours`
  - `should navigate to previous week when backward button is clicked`
  - `should update weekly hours display when navigating backward`
  - `should allow multiple backward navigations`
  - `should update weekly hours when a shift is released`
  - `should maintain released shift state when navigating weeks`

---

## Backend Unit Tests

### Controller Tests

#### Company Controller
- `TestCreateCompany_Success`
- `TestCreateCompany_ValidationError`
- `TestListCompanies_Success`
- `TestListCompanies_DBFailure`

#### Shift Controller
- `TestCreateShift_Success`
- `TestCreateShift_InvalidStartTimeFormat`
- `TestCreateShift_InvalidEndTimeFormat`
- `TestCreateShift_Unauthorized`
- `TestCreateShift_ServiceError`
- `TestAssignWorkerToShift_Success`
- `TestAssignWorkerToShift_MissingFields`
- `TestAssignWorkerToShift_Unauthorized`
- `TestReleaseShiftForWorker_Success`
- `TestReleaseShiftForWorker_NotFound`
- `TestReleaseShiftForWorker_InvalidStatus_Controller`
- `TestRequestShift_Success`
- `TestRequestShift_Unauthorized`
- `TestRequestShift_BadShiftID`
- `TestRequestShift_NotFound`
- `TestRejectShiftRequest_Unauthorized`
- `TestRejectShiftRequest_BadShiftID`
- `TestRejectShiftRequest_NotFound`
- `TestRejectShiftRequest_Success`

#### Supervisor Controller
- `TestSupervisorLogin_Success`
- `TestSupervisorLogin_InvalidPassword`
- `TestSupervisorLogin_UserNotFound`
- `TestSupervisorLogin_BadRequest_InvalidPayload`
- `TestSupervisorRegister_Success`
- `TestSupervisorRegister_InvalidCompany`
- `TestSupervisorRegister_ValidationError`
- `TestGetWorkerAvailability_Success`
- `TestGetWorkerAvailability_InvalidDate`
- `TestGetWorkerAvailability_MissingQuery`
- `TestGetShiftsByCompany_Controller_ValidationAndSuccess`
- `TestGetShiftsByCompany_Controller_SupervisorNotFound`
- `TestGetRequestedShifts_Controller_Success`
- `TestGetRequestedShifts_Unauthorized`
- `TestGetRequestedShifts_SupervisorNotFound`
- `TestGetAssignedShifts_Unauthorized`
- `TestGetAssignedShifts_SupervisorNotFound`
- `TestGetAssignedShifts_Success_Empty`
- `TestGetWorkersWithOvertimeHours_Controller_Success`
- `TestGetWorkersWithOvertimeHours_Controller_DBFailure`

#### Worker Controller
- `TestWorkerRegister_BadRequest_InvalidPayload`
- `TestWorkerLogin_BadRequest_InvalidPayload`
- `TestGetShiftsForWorker_BadRequest_InvalidEmployeeID`
- `TestCreateAvailability_BadRequest_InvalidEmployeeID`
- `TestCreateAvailability_BadRequest_InvalidBody`
- `TestCreateAvailability_WorkerNotFound`
- `TestCreateAvailability_MalformedJSON`
- `TestGetShiftsForWorker_NotFound`
- `TestGetShiftsForWorker_Success_Empty`
- `TestWorkerRegister_Success`
- `TestWorkerRegister_RegisterEmployeeError_InvalidCompany`
- `TestWorkerLogin_Success`
- `TestWorkerLogin_InvalidPassword`
- `TestWorkerLogin_UserNotFound`
- `TestGetShiftsForWorker_Success_WithWindowFilter`
- `TestGetShiftsForWorker_BadRequest_IncompleteWindow`
- `TestGetShiftsForWorker_BadRequest_InvalidWindowOrder`
- `TestCreateAvailability_Success`
- `TestGetReleasedShiftsForWorkerCompany_WorkerNotFound`
- `TestGetReleasedShiftsForWorkerCompany_Unauthorized`
- `TestGetReleasedShiftsForWorkerCompany_Success`

### Service Tests

#### Company Service
- `TestCreateCompanyService_Success`
- `TestCreateCompanyService_PersistsAddress`
- `TestCreateCompanyService_InvalidName`
- `TestListCompaniesService_Success`
- `TestListCompaniesService_Empty`

#### Shift Service
- `TestReleaseShiftForWorker_Success`
- `TestReleaseShiftForWorker_NotFound`
- `TestReleaseShiftForWorker_InvalidStatus`
- `TestCreateShift_Success`
- `TestCreateShift_SupervisorNotFound`
- `TestCreateShift_InvalidTimeOrder`
- `TestCreateShift_SameStartAndEndTime`
- `TestAssignWorkerToShift_Success`
- `TestAssignWorkerToShift_UpdatesExistingRequest`
- `TestAssignWorkerToShift_ShiftNotFound`
- `TestAssignWorkerToShift_WorkerNotFound`
- `TestAssignWorkerToShift_SupervisorNotAuthorized`
- `TestAssignWorkerToShift_WorkerHasWrongRole`
- `TestRejectShiftRequest_Success`
- `TestRejectShiftRequest_NotFound_WrongCompany`
- `TestRejectShiftRequest_NotFound_WhenSupervisorMissing`
- `TestRequestReleasedShift_WorkerWrongCompany`
- `TestRequestReleasedShift_WorkerNotFound`
- `TestShiftService_DeleteShift`
- `TestShiftService_UnassignWorker`

#### Supervisor Service
- `TestSupervisorService_GetShiftsByCompany_NotFound`
- `TestSupervisorService_GetShiftsByCompany_Success`
- `TestSupervisorService_GetShiftsByCompany_UnassignedFilter`
- `TestSupervisorService_GetShiftsByCompany_NoFilterUsesLatestAssignment`
- `TestWorkerService_GetWorkersAvailable_Success`
- `TestWorkerService_GetWorkersAvailable_NoResults`
- `TestSupervisorService_GetWorkersWithOvertimeHours_Success`
- `TestSupervisorService_GetWorkersWithOvertimeHours_NotFound`
- `TestSupervisorService_GetWorkersWithOvertimeHours_ExcludesWorkersAtOrBelowThreshold`
- `TestSupervisorService_GetCompanyWorkers`

#### Worker Service
- `TestCreateWorkerAvailability`
- `TestGetAssignedShifts_Success`
- `TestGetAssignedShifts`
- `TestGetReleasedShiftsByEmployeeCompany_Success`
- `TestGetReleasedShiftsByEmployeeCompany_EmployeeNotWorker`
- `TestGetReleasedShiftsByEmployeeCompany_Empty`
- `TestRequestReleasedShift_Success`
- `TestRequestReleasedShift_NotReleased`

---

## Backend API Documentation

### Health
- `GET /health`
  - Returns `200 OK` when the backend is running.

### Company APIs
- `POST /api/companies/create`
  - Auth: public.
  - Body:
    ```json
    {
      "name": "Acme Inc",
      "address": "123 Main St"
    }
    ```
  - Success: `201 Created`.
  - Error: `400 Bad Request` for validation or service errors.

- `GET /api/companies/`
  - Auth: public.
  - Success: `200 OK` with a list of companies.
  - Error: `500 Internal Server Error` if companies cannot be fetched.

### Supervisor APIs
- `POST /api/supervisors/register`
  - Auth: public.
  - Body:
    ```json
    {
      "name": "Supervisor One",
      "email": "supervisor@example.com",
      "password": "password123",
      "address": "123 Main St",
      "phone_no": "5551234567",
      "company_id": 1,
      "wage": 40
    }
    ```
  - Success: `201 Created` with JWT token and employee object.
  - Error: `400 Bad Request` for invalid input or registration errors.

- `POST /api/supervisors/login`
  - Auth: public.
  - Body:
    ```json
    {
      "email": "supervisor@example.com",
      "password": "password123"
    }
    ```
  - Success: `200 OK` with JWT token and employee object.
  - Error: `400 Bad Request`, `401 Unauthorized`, or `500 Internal Server Error`.

- `GET /api/supervisors/workers/availability?date=YYYY-MM-DD&company_id=1`
  - Auth: JWT + supervisor role.
  - Success: `200 OK` with available workers for the requested date.
  - Error: `400 Bad Request` for missing/invalid query params, `500 Internal Server Error` for fetch failures.

- `GET /api/supervisors/:employee_id/shifts?status=Assigned`
  - Auth: JWT + supervisor role.
  - Optional status values: `Assigned`, `Released`, `Rejected`, `Requested`, `Unassigned`.
  - Success: `200 OK` with shifts for the supervisor's company.
  - Error: `400 Bad Request`, `404 Not Found`, or `500 Internal Server Error`.

- `GET /api/supervisors/shifts/requested`
  - Auth: JWT + supervisor role.
  - Success: `200 OK` with requested shifts.
  - Error: `401 Unauthorized`, `404 Not Found`, or `500 Internal Server Error`.

- `GET /api/supervisors/shifts/assigned`
  - Auth: JWT + supervisor role.
  - Success: `200 OK` with assigned shifts.
  - Error: `401 Unauthorized`, `404 Not Found`, or `500 Internal Server Error`.

- `PATCH /api/supervisors/shifts/:shift_id/reject`
  - Auth: JWT + supervisor role.
  - Success: `200 OK` with the updated shift assignment.
  - Error: `400 Bad Request`, `401 Unauthorized`, `404 Not Found`, or `500 Internal Server Error`.

- `GET /api/supervisors/workers/overtime`
  - Auth: JWT + supervisor role.
  - Success: `200 OK` with workers over 20 assigned hours for the current week.
  - Error: `401 Unauthorized`, `404 Not Found`, or `500 Internal Server Error`.

- `GET /api/supervisors/company/workers`
  - Auth: JWT + supervisor role.
  - Success: `200 OK` with a list of workers (`id`, `name`, `email`, `phone_no`) belonging to the supervisor's company.
  - Error: `404 Not Found` or `500 Internal Server Error`.

### Worker APIs
- `POST /api/workers/register`
  - Auth: public.
  - Body:
    ```json
    {
      "name": "Worker One",
      "email": "worker@example.com",
      "password": "password123",
      "address": "123 Main St",
      "phone_no": "5551234567",
      "company_id": 1,
      "wage": 22.5
    }
    ```
  - Success: `201 Created` with JWT token and employee object.
  - Error: `400 Bad Request` for invalid input or registration errors.

- `POST /api/workers/login`
  - Auth: public.
  - Body:
    ```json
    {
      "email": "worker@example.com",
      "password": "password123"
    }
    ```
  - Success: `200 OK` with JWT token and employee object.
  - Error: `400 Bad Request`, `401 Unauthorized`, or `500 Internal Server Error`.

- `GET /api/workers/:employee_id/shifts`
  - Auth: JWT + worker role.
  - Optional query params: `start_time` and `end_time` in RFC3339 format. Both must be provided together.
  - Success: `200 OK` with assigned shifts, sorted by start time. Responses include duration and earnings.
  - Error: `400 Bad Request`, `404 Not Found`, or `500 Internal Server Error`.

- `POST /api/workers/:employee_id/availability`
  - Auth: JWT + worker role.
  - Body:
    ```json
    {
      "day_of_week": 1,
      "start_time": "09:00:00",
      "end_time": "17:00:00"
    }
    ```
  - Success: `201 Created`.
  - Error: `400 Bad Request` or `404 Not Found`.

- `GET /api/workers/shifts/released`
  - Auth: JWT + worker role.
  - Success: `200 OK` with released shifts from the worker's company.
  - Error: `401 Unauthorized`, `404 Not Found`, or `500 Internal Server Error`.

- `POST /api/workers/shifts/:shift_id/request`
  - Auth: JWT + worker role.
  - Success: `200 OK` with the updated requested assignment.
  - Error: `400 Bad Request`, `401 Unauthorized`, or `404 Not Found`.

### Shift APIs
- `POST /api/shifts/create`
  - Auth: JWT + supervisor role.
  - Body:
    ```json
    {
      "start_time": "2026-02-15T09:00:00Z",
      "end_time": "2026-02-15T17:00:00Z"
    }
    ```
  - Success: `201 Created`.
  - Error: `400 Bad Request` or `401 Unauthorized`.

- `POST /api/shifts/assign`
  - Auth: JWT + supervisor role.
  - Body:
    ```json
    {
      "shift_id": 1,
      "employee_id": 2
    }
    ```
  - Success: `201 Created`.
  - Error: `400 Bad Request` or `401 Unauthorized`.

- `POST /api/shifts/release`
  - Auth: JWT + worker role.
  - Body:
    ```json
    {
      "shift_id": 1
    }
    ```
  - Success: `200 OK`.
  - Error: `400 Bad Request`, `401 Unauthorized`, `404 Not Found`, or `500 Internal Server Error`.

- `DELETE /api/shifts/:shift_id`
  - Auth: JWT + supervisor role.
  - Success: `200 OK` shift and underlying assignments permanently deleted.
  - Error: `403 Forbidden` (if shift has active assignments or wrong company), `404 Not Found`, or `500 Internal Server Error`.

- `PATCH /api/shifts/:shift_id/unassign`
  - Auth: JWT + supervisor role.
  - Success: `200 OK` worker unassigned, shift status reset.
  - Error: `403 Forbidden` (wrong company), `404 Not Found`, or `500 Internal Server Error`.

---

## Verification
- `go test ./controllers -cover` passed with `83.3%` controller coverage.
- `go test ./services -cover` passed with `55.3%` service coverage.
- `go test ./...` passed for the backend.
