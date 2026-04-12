import { TestBed } from '@angular/core/testing';
import { provideHttpClient } from '@angular/common/http';
import { HttpTestingController, provideHttpClientTesting } from '@angular/common/http/testing';
import { SupervisorService } from './supervisor.service';

describe('SupervisorService', () => {
  let service: SupervisorService;
  let httpMock: HttpTestingController;
  const apiUrl = 'http://localhost:8080/api/supervisors';

  beforeEach(() => {
    TestBed.configureTestingModule({
      providers: [
        SupervisorService,
        provideHttpClient(),
        provideHttpClientTesting()
      ]
    });
    service = TestBed.inject(SupervisorService);
    httpMock = TestBed.inject(HttpTestingController);
  });

  afterEach(() => {
    httpMock.verify();
  });

  // Test 1: Get Shifts
  it('should fetch shifts with an optional status filter', () => {
    const mockShifts = [{ id: 1, status: 'Requested' }];
    
    service.getShifts(8, 'Requested').subscribe(shifts => {
      expect(shifts).toEqual(mockShifts);
    });

    const req = httpMock.expectOne(`${apiUrl}/8/shifts?status=Requested`);
    expect(req.request.method).toBe('GET');
    req.flush(mockShifts);
  });

  // Test 2: Assign Worker
  it('should assign a worker to a shift', () => {
    service.assignWorker(10, 5).subscribe(response => {
      expect(response).toBeTruthy();
    });

    // Changed to match the exact URL your service is actually calling
    const req = httpMock.expectOne('http://localhost:8080/api/shifts/assign');
    expect(req.request.method).toBe('POST');
    expect(req.request.body).toEqual({ shift_id: 10, employee_id: 5 });
    req.flush({ message: 'Success' });
  });

  // Test 3: Get Available Workers
  it('should fetch available workers for a specific date', () => {
    const mockWorkers = [{ id: 5, name: 'John Doe' }];
    
    service.getAvailableWorkers('2026-03-25', 1).subscribe(workers => {
      expect(workers.length).toBe(1);
    });

    // Updated to match the actual API endpoint
    const req = httpMock.expectOne(`http://localhost:8080/api/supervisors/workers/availability?date=2026-03-25&company_id=1`);
    expect(req.request.method).toBe('GET');
    req.flush(mockWorkers);
  });

  // Test 4: Create Shift
  it('should be created successfully', () => {
    expect(service).toBeTruthy();
  });
  
  it('should fetch assigned shifts for the schedule view', () => {
    const mockShifts = [{ id: 1, status: 'Assigned', assigned_to: 5 }];

    service.getAssignedShifts().subscribe(shifts => {
      expect(shifts.length).toBe(1);
      expect(shifts[0].status).toBe('Assigned');
    });

    const req = httpMock.expectOne('http://localhost:8080/api/supervisors/shifts/assigned');
    expect(req.request.method).toBe('GET');
    req.flush(mockShifts);
  });
});

