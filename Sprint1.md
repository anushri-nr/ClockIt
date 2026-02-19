# Sprint 1 - ClockIt

## Overview
Sprint 1 focused on establishing the project foundation, setting up the development environment, and implementing core registration, login, and API infrastructure for the ClockIt shift management application.

---

## User Stories

### 1. **Project Setup & Infrastructure**
- Initialize Go module and basic server structure
- Set up Angular frontend with CI/CD pipeline
- Configure SQLite database with GORM ORM

### 2. **Registration & Login**
- Supervisor login and registration (Frontend & Backend)
- Worker login and registration (Frontend & Backend)

### 3. **Dashboard & UI**
- Supervisor dashboard for shift management
- Worker dashboard for viewing assigned shifts
- Material Component Library integration

### 4. **Core APIs**
- API to fetch all companies
- Supervisor API to list all shifts created by supervisor
- Worker API to list all shifts assigned to worker
- API to fetch worker availability for a specific day
- API to assign shifts to workers

### 5. **Database Design**
- Design and implement database schema
- Add worker availability tracking

---

## Planned Issues & Status

| Issue | Title | Status | Assignee |
|-------|-------|--------|----------|
| #2 | Project Skeleton: Initialize Go Module & Basic Server | Complete | Anushri |
| #3 | Project Skeleton: Initialize Frontend (Angular) & CI Setup | Complete | Rahul, Varshith |
| #4 | Database: Setup GORM with SQLite connection | Complete | Avantika |
| #5 | UI: Install Material Component Library | Complete | Varshith, Rahul |
| #22 | Database Design | Complete | Avantika, Anushri |
| #13 | Landing Page: ClockIt Frontend | Complete | Varshith |
| #15 | Worker Registration: Backend | Complete | Anushri |
| #18 | Supervisor: Post API to assign shift to Worker | Complete | Anushri |
| #27 | Worker: API to list all shifts assigned to a worker | Complete | Avantika |
| #28 | Supervisor: API to list all shifts created by a supervisor | Complete | Avantika |
| #29 | API to fetch all companies | Complete | Avantika |
| #14 | Supervisor: Login and Registration: Frontend | Complete | Rahul |
| #16 | Supervisor Registration: Backend | Complete | Avantika |
| #17 | Supervisor: API to fetch Worker Availability for the day | Complete | Avantika |
| #19 | Supervisor: Dashboard View: Frontend | Complete | Rahul |
| #20 | Worker: Dashboard View: Frontend | Complete | Varshith |
| #21 | Worker: Login and Registration: Frontend | Complete | Rahul |
| #25 | Implement supervisor registration and fetch worker availability for a specific day | Complete | Avantika |
| #31 | Supervisor: Create Shift | Complete | Anushri |
| #32 | Company: Create API | Complete | Anushri |

## Successes
**No Spillover**: All planned tasks were completed.

**Project Foundation**: Established complete Go backend and Angular frontend structure.

**Database**: SQLite database configured with proper ORM integration. 

**Registration and Login**: Core login/registration flows implemented for both user types.

**Core APIs**: All planned shift management and assignment APIs completed.

**UI Framework**: Material Design system integrated for consistent frontend.  

---

## Blockers / Challenges

None identified at this stage. All tasks are progressing through normal review cycles. Some tasks are planned for Sprint 2, so we have mocked data for now.

---

## Next Steps (Sprint 2)

- Authentication and Authorization
- Expand worker availability tracking features
- Implement additional shift management features
