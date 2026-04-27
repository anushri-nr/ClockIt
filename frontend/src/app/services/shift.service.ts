import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { environment } from '../../environments/environment';

@Injectable({
  providedIn: 'root'
})
export class ShiftService {

  // Update this with your real backend port (usually 8080 for Go)
  private apiUrl = `${environment.apiUrl}/shifts`;

  constructor(private http: HttpClient) { }

  createShift(shiftData: any): Observable<any> {
    // 1. Format the data to match Go's expectations
    // Go usually wants RFC3339 format for time (e.g., "2026-02-15T09:00:00Z")
    // The HTML input gives us "2026-02-15T09:00"
    
    const payload = {
      start_time: new Date(shiftData.startTime).toISOString(), 
      end_time: new Date(shiftData.endTime).toISOString(),
      // We assume the backend handles "created_at" automatically
      // We might need to send the supervisor's ID if the backend requires it
      created_by: Number(shiftData.createdBy) 
    };

    console.log('Creating Shift Payload:', payload);

    return this.http.post(`${this.apiUrl}/create`, payload);
  }
}