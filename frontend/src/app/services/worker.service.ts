import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { environment } from '../../environments/environment';

// Matches 'AssignedShiftResponse' in backend/services/worker_service.go
export interface WorkerShift {
  id: number;          // json:"id" (Assignment ID)
  shift_id: number;    // json:"shift_id"
  start_time: string;  // json:"start_time"
  end_time: string;    // json:"end_time"
  assigned_by: number; // json:"assigned_by"
  status: string;      // json:"status"
}

export interface AvailabilityPayload {
  day_of_week: number; // 0 = Sunday, 1 = Monday, etc.
  start_time: string;  // e.g., "09:00"
  end_time: string;    // e.g., "17:00"
}

@Injectable({
  providedIn: 'root'
})
export class WorkerService {

  private apiUrl = `${environment.apiUrl}/workers`;

  constructor(private http: HttpClient) { }

  // GET /api/workers/:employee_id/shifts
  getMyShifts(workerId: number): Observable<WorkerShift[]> {
    return this.http.get<WorkerShift[]>(`${this.apiUrl}/${workerId}/shifts`);
  }

  createAvailability(workerId: number, payload: AvailabilityPayload): Observable<any> {
    return this.http.post(`${this.apiUrl}/${workerId}/availability`, payload);
  }

  // Tell the Go backend to release the shift
  releaseShift(shiftId: number): Observable<any> {
    return this.http.post(`${environment.apiUrl}/shifts/release`, { shift_id: shiftId });
  }
}