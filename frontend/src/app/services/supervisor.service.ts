import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable, of, delay } from 'rxjs';

@Injectable({
  providedIn: 'root'
})
export class SupervisorService {

  private apiUrl = 'http://localhost:8080/api/supervisors';

  constructor(private http: HttpClient) { }

  // GET /api/supervisors/workers/availability
  getAvailableWorkers(): Observable<any[]> {
    // UNCOMMENT THIS WHEN BACKEND IS READY
    // return this.http.get<any[]>(`${this.apiUrl}/workers/availability`);

    // MOCK DATA (Matches the structure we expect)
    const dummyWorkers = [
      { id: 101, name: 'Ava N.', role: 'Front Desk' },
      { id: 102, name: 'Malik J.', role: 'Stock' },
      { id: 103, name: 'Nora K.', role: 'Cashier' },
      { id: 104, name: 'Sam R.', role: 'Barista' }
    ];
    
    console.log('Fetching available workers...');
    return of(dummyWorkers).pipe(delay(800)); 
  }
}