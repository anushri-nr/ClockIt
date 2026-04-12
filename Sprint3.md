# Sprint 3 - ClockIt

## Overview
Sprint 3 focused on supervisor workflow improvements, worker shift request flows, pending approvals, released/assigned shift management, overtime risk handling, and frontend dashboard enhancements.

---

## Work Completed in Sprint 3

### Backend
- Added API support for supervisor assigned-shift retrieval.
- Added API support for worker released-shift retrieval.
- Added API support for worker shift request flow.
- Added API support for supervisor shift rejection flow.
- Implemented status-based shift filtering using `shift_assignments`.
- Improved authentication handling in protected controller methods.
- Added and updated controller/service tests for shift workflow endpoints.
- Hardened test DB setup to prevent SQLite state leakage across tests.

### Frontend
- Built supervisor alerts and tasks UI.
- Added supervisor overtime risk views/features.
- Added supervisor schedule view.
- Added released shifts management dashboard for supervisors.
- Added pending approvals UI and request/approval flows.
- Added upcoming assigned shifts tab fixes.
- Added next-shift assigned tab.
- Added worker weekly hours assigned check and implementation.
- Added logout flow for frontend.

---

## Frontend Unit Tests

### Supervisor Service (`supervisor.service.spec.ts`)
- `should fetch assigned shifts for the schedule view` (Verifies GET parameters and response mapping).
- `should send a PATCH request to reject a requested shift` (Verifies proper HTTP method and URL construction).
- `should fetch workers at risk of overtime` (Verifies GET request to overtime risk endpoint).

## Frontend E2E Tests (Cypress)

### Supervisor Dashboard (`supervisor-login.cy.ts`)
- Verified the Schedule View successfully renders assigned shift data.
- Verified the Pending Approvals modal opens and hydrates requested shift data.
- Verified the Alerts & tasks side panel resolves loading states and correctly displays overtime risk warnings.

## Backend Unit Tests

### Controllers
- Supervisor assigned shifts endpoint tests
- Worker released shifts endpoint tests
- Worker request shift endpoint tests
- Supervisor reject shift endpoint tests
- Supervisor get overworked workers endpoint tests


### Services
- `GetReleasedShiftsByEmployeeCompany` tests
- `RequestReleasedShift` tests
- `RejectShiftRequest` tests
- `GetAssignedShifts` tests
- `CreateWorkerAvailability` tests
- `GetWorkersWithOvertimeHours` tests

---

## Backend API Specification

### Supervisor APIs
- `POST /api/supervisors/register` — Register a supervisor.
- `POST /api/supervisors/login` — Authenticate a supervisor and return JWT.
- `GET /api/supervisors/shifts/assigned` — Get assigned shifts for the supervisor’s company.
- `GET /api/supervisors/workers/availability` — Fetch worker availability for scheduling.
- `GET /api/supervisors/shifts/requested` — Get requested shifts for the supervisor’s company.
- `PATCH /api/supervisors/shifts/:shift_id/reject` — Reject a requested shift.
- `GET /api/supervisors/workers/overtime` — Get all workers whose total assigned shift hours for the current week exceed 20 hours.

### Worker APIs
- `POST /api/workers/register` — Register a worker.
- `POST /api/workers/login` — Authenticate a worker and return JWT.
- `GET /api/workers/shifts/released` — Get released shifts from the worker’s company.
- `POST /api/workers/shifts/:shift_id/request` — Request a released shift.
- `GET /api/workers/:employee_id/shifts` — Get assigned shifts for a worker.
- `POST /api/workers/:employee_id/availability` — Create worker availability.

### Shift APIs
- `POST /api/shifts/create` — Create a shift.
- `POST /api/shifts/assign` — Assign a worker to a shift.
- `GET /api/shifts/:shift_id` — Fetch shift details if implemented.
- `PATCH /api/shifts/:shift_id/release` — Release a shift if implemented.

### Company APIs
- `POST /api/companies/create` — Create a company.
- `GET /api/companies/` — List companies.

---