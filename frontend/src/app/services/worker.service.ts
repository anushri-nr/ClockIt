import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';

// Matches 'AssignedShiftResponse' in backend/services/worker_service.go
export interface WorkerShift {
  id: number;          // json:"id" (Assignment ID)
  shift_id: number;    // json:"shift_id"
  start_time: string;  // json:"start_time"
  end_time: string;    // json:"end_time"
  assigned_by: number; // json:"assigned_by"
  status: string;      // json:"status"
}

@Injectable({
  providedIn: 'root'
})
export class WorkerService {

  private apiUrl = 'http://localhost:8080/api/workers';

  constructor(private http: HttpClient) { }

  // GET /api/workers/:employee_id/shifts
  getMyShifts(workerId: number): Observable<WorkerShift[]> {
    return this.http.get<WorkerShift[]>(`${this.apiUrl}/${workerId}/shifts`);
  }
}