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

  // GET /api/workers/:employee_id/shifts?start_time=...&end_time=...
  getMyShifts(workerId: number, startTime?: string, endTime?: string): Observable<WorkerShift[]> {
    const params: any = {};
    if (startTime && endTime) {
      params.start_time = startTime;
      params.end_time = endTime;
    }
    return this.http.get<WorkerShift[]>(`${this.apiUrl}/${workerId}/shifts`, { params });
  }

  // POST /api/shifts/release
  releaseShift(shiftId: number): Observable<any> {
    return this.http.post('http://localhost:8080/api/shifts/release', { shift_id: shiftId });
  }
}
