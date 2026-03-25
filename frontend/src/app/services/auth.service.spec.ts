import { TestBed } from '@angular/core/testing';
import { provideHttpClient } from '@angular/common/http';
import { HttpTestingController, provideHttpClientTesting } from '@angular/common/http/testing';
import { AuthService } from './auth.service';

describe('AuthService', () => {
  let service: AuthService;
  let httpMock: HttpTestingController;
  const apiUrl = 'http://localhost:8080/api';

  beforeEach(() => {
    TestBed.configureTestingModule({
      providers: [
        AuthService,
        provideHttpClient(),
        provideHttpClientTesting()
      ]
    });
    service = TestBed.inject(AuthService);
    httpMock = TestBed.inject(HttpTestingController);
    localStorage.clear(); // Start fresh each time
  });

  afterEach(() => {
    httpMock.verify();
  });

  // Test 1: Initialization
  it('should be created and have a null initial user', () => {
    expect(service).toBeTruthy();
    expect(service.currentUserValue).toBeNull();
  });

  // Test 2: Login Function
  it('should handle login and update local storage', () => {
    const mockResponse = {
      token: 'fake-jwt',
      employee: { id: 1, name: 'Test Supervisor', company_id: 1 }
    };

    service.login('test@test.com', 'password', 'supervisor').subscribe(success => {
      expect(success).toBe(true);
      expect(localStorage.getItem('jwt_token')).toBe('fake-jwt');
      expect(service.currentUserValue?.name).toBe('Test Supervisor');
    });

    const req = httpMock.expectOne(`${apiUrl}/supervisors/login`);
    expect(req.request.method).toBe('POST');
    req.flush(mockResponse);
  });

  // Test 3: Register Function
  it('should handle registration and auto-login', () => {
    const newUser = { name: 'New Guy', email: 'new@test.com', password: '123', role: 'worker', companyId: 1 };
    const mockResponse = {
      token: 'new-jwt',
      employee: { id: 2, name: 'New Guy', company_id: 1 }
    };

    service.register(newUser).subscribe(success => {
      expect(success).toBe(true);
      expect(localStorage.getItem('jwt_token')).toBe('new-jwt');
    });

    const req = httpMock.expectOne(`${apiUrl}/workers/register`);
    expect(req.request.method).toBe('POST');
    req.flush(mockResponse);
  });

  // Test 4: Get Companies Function
  it('should fetch the list of companies', () => {
    const mockCompanies = [{ id: 1, name: 'Panda Express' }];

    service.getCompanies().subscribe(companies => {
      expect(companies.length).toBe(1);
      expect(companies[0].name).toBe('Panda Express');
    });

    const req = httpMock.expectOne(`${apiUrl}/companies/`);
    expect(req.request.method).toBe('GET');
    req.flush(mockCompanies);
  });
});