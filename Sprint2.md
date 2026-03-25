# Sprint 2 - ClockIt

## Overview
Sprint 2 focused on strengthening backend capabilities, improving authentication and shift APIs, adding worker-facing payroll and availability features, and expanding the frontend experience for supervisors and workers.

---

## User Stories

### 1. **Authentication & Access Control**
- Add login support for supervisors
- Add login support for workers
- Implement backend authentication and authorization flows

### 2. **Shift API Enhancements**
- Update shift APIs to include start and end time
- Return sorted shift lists from the backend
- Modify supervisor shift retrieval behavior

### 3. **Worker Features**
- Add worker wage calculation functionality
- Build worker availability backend API

### 4. **Supervisor & Frontend Features**
- Create supervisor pending approval dashboard
- Add open shifts view on the frontend
- Fix profile display so the logged-in user name is visible
- Add worker wage view on the dashboard
- Add upcoming shift display on the worker dashboard
- Add weekly hours assigned view on the worker dashboard
- Add estimated weekly earnings on the worker dashboard

### 5. **Quality Improvements**
- Expand backend unit test coverage
- Implement Cypress End-to-End (E2E) tests for frontend authentication flows
- Write Angular Unit Tests (Jasmine/Karma) for core frontend services to reach a 1:1 function ratio

---

## Planned Issues & Status

| Issue | Title | Status | Owner |
|-------|-------|--------|-------|
| #35 | Backend: Login Supervisor API | Complete | avantikahollas |
| #36 | Backend: Login Worker API | Complete | farcyson |
| #37 | Backend: Authentication and Authorization | Complete | avantikahollas, farcyson |
| #38 | Backend: Wage Calculation: Worker Functionality | Complete | farcyson |
| #45 | Enchancement: Fix Backend Shifts API to return a sorted list | Complete | anushri-nr |
| #47 | Enhancement: Get Shifts API to have Start and End Time | Complete | avantikahollas |
| #50 | Frontend: Pending Approval: Supervisor Dashboard | Complete | farcyson |
| #51 | Frontend: Fix: Logged in user name must be visible on the profile | Complete | farcyson |
| #52 | Frontend: Open Shifts View | Complete | farcyson |
| #53 | Backend: Enhancement: Modify Get Shifts by Supervisor API | Complete | anushri-nr |
| #59 | Backend: Write Unit Tests | Complete | avantikahollas, anushri-nr |
| #65 | Worker Availability API: Backend | Complete | anushri-nr |
| #43 | Frontend: Wage View - Worker Dashboard | Complete | RahulUmamahesha |
| #44 | Frontend: Upcoming Shift Display - Worker Dashboard | Complete | RahulUmamahesha |
| #46 | Frontend: Weekly Hours Assigned View: Worker Dashboard | Complete | RahulUmamahesha |
| #48 | Frontend: Estimated Earnings for the week: Worker Dashboard | Complete | RahulUmamahesha |

## Successes
**Authentication Added**: Supervisor and worker login APIs were completed along with backend authentication and authorization support.

**Shift APIs Improved**: Shift responses were enhanced with start/end times, supervisor shift retrieval updates, and sorted backend results.

**Frontend Expanded**: The supervisor dashboard, open shifts view, and profile visibility fixes improved the overall usability of the application.

**Worker Support Grew**: Wage calculation and worker availability backend features added important functionality for employee workflows.

**Testing Improved**: Additional backend unit tests strengthened confidence in the sprint deliverables.

**Carryover Identified**: Worker dashboard enhancements for wage view, upcoming shifts, weekly hours, and estimated earnings are actively in progress and will continue into the next sprint if not completed in time.

**Frontend Integration & Testing:** Successfully integrated JWT authentication state across the frontend using a centralized `BehaviorSubject`, and verified the entire Supervisor flow using automated Cypress E2E tests and Angular unit tests.

---

## Blockers / Challenges

No major blockers were recorded in the sprint issue list. The main remaining work involves completing the in-progress worker dashboard features and connecting them cleanly with the backend data already added during this sprint.

---

## Next Steps (Sprint 3)

- Continue expanding automated test coverage across controllers and services
- Integrate the new backend APIs more deeply with frontend workflows
- Complete the in-progress worker dashboard views for wage, hours, earnings, and upcoming shifts
- Refine approval, shift assignment, and availability experiences
- Address remaining open issues and polish end-to-end user flows