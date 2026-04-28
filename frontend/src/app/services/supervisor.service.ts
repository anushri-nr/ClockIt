import { Injectable } from '@angular/core';
import { HttpClient, HttpParams } from '@angular/common/http';
import { Observable } from 'rxjs';

export interface Shift {
  id: number;
  shift_id: number;
  start_time: string;
  end_time: string;
  assigned_to?: string; // Optional because it might be unassigned
  status: string;
}

export interface Worker {
  id: number;
  name: string;
  email: string;
  phone_no: string;
  start_time: string;
  end_time: string;
  company_id: number;
  company_name: string;
}

@Injectable({
  providedIn: 'root'
})
export class SupervisorService {

  private apiUrl = 'http://localhost:8080/api/supervisors';

  constructor(private http: HttpClient) { }

  // 1. GET SHIFTS (Real API - Updated with Status Filter)
  // Backend Route: GET /api/supervisors/:employee_id/shifts?status=XYZ
  getShifts(supervisorId: number, status?: string): Observable<Shift[]> {
    let params = new HttpParams();
    
    // If a status was passed in (e.g., 'Unassigned' or 'Requested'), add it to the query string
    if (status) {
      params = params.set('status', status);
    }

    return this.http.get<Shift[]>(`${this.apiUrl}/${supervisorId}/shifts`, { params });
  }

  // 2. GET AVAILABLE WORKERS (Real API)
  // Backend Route: GET /api/supervisors/workers/availability?date=YYYY-MM-DD&company_id=X
  getAvailableWorkers(date: string, companyId: number): Observable<Worker[]> {
    
    // Setup Query Parameters
    let params = new HttpParams()
      .set('date', date)
      .set('company_id', companyId.toString());

    return this.http.get<Worker[]>(`${this.apiUrl}/workers/availability`, { params });
  }

  // 3. ASSIGN WORKER TO SHIFT
  // Backend Route: POST /api/shifts/assign
  assignWorker(shiftId: number, workerId: number): Observable<any> {
    const payload = {
      shift_id: Number(shiftId),
      employee_id: Number(workerId),
    };
    // Note: The route is actually under 'shifts', not 'supervisors'
    // So we use a different base URL for this specific call
    return this.http.post('http://localhost:8080/api/shifts/assign', payload);
  }

  /// Fetch all assigned shifts for the schedule view
  getAssignedShifts(): Observable<any[]> {
    return this.http.get<any[]>(`${this.apiUrl}/shifts/assigned`);
  }

  // Fetch shifts that workers have requested
  getRequestedShifts(): Observable<any[]> {
    return this.http.get<any[]>(`${this.apiUrl}/shifts/requested`);
  }

  // Reject a requested shift
  rejectShiftRequest(shiftId: number): Observable<any> {
    return this.http.patch(`${this.apiUrl}/shifts/${shiftId}/reject`, {});
  }

  // Fetch workers at risk of overtime
  getOvertimeRiskWorkers(): Observable<any[]> {
    return this.http.get<any[]>(`${this.apiUrl}/workers/overtime`);
  }

  // Delete a shift
  deleteShift(shiftId: number): Observable<any> {
    return this.http.delete(`http://localhost:8080/api/shifts/${shiftId}`);
  }

  // Unassign a worker from a shift
  unassignWorker(shiftId: number): Observable<any> {
    return this.http.patch(`http://localhost:8080/api/shifts/${shiftId}/unassign`, {});
  }
}