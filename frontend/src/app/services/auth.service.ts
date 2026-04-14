import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable, of, BehaviorSubject } from 'rxjs';
import { tap, map, catchError } from 'rxjs/operators';

// 1. Define exactly what a user looks like
export interface User {
  id: number;
  name: string;
  email: string;
  role: string;
  company_id: number;
}

@Injectable({
  providedIn: 'root'
})
export class AuthService {

  private apiUrl = 'http://localhost:8080/api';

  // 2. The Single Source of Truth for the logged-in user
  private currentUserSubject = new BehaviorSubject<User | null>(null);
  public currentUser$ = this.currentUserSubject.asObservable();

  constructor(private http: HttpClient) { 
    this.loadUserFromStorage(); // Automatically reload user if they refresh the page
  }

  // Helper to synchronously grab the current user whenever a component needs it
  public get currentUserValue(): User | null {
    return this.currentUserSubject.value;
  }

  // Checks local storage on startup
  private loadUserFromStorage() {
    const token = localStorage.getItem('jwt_token');
    const userData = localStorage.getItem('user_data');
    if (token && userData) {
      try {
        this.currentUserSubject.next(JSON.parse(userData));
      } catch (e) {
        console.error('Failed to parse user data from storage');
      }
    }
  }

  // Centralized function to save data on login/register
  private handleAuthResponse(response: any) {
    if (response && response.token && response.employee) {
      localStorage.setItem('jwt_token', response.token);
      localStorage.setItem('user_data', JSON.stringify(response.employee)); // Store as one clean JSON object
      this.currentUserSubject.next(response.employee); // Broadcast to the app
      console.log('Auth centralized state updated for:', response.employee.name);
    }
  }

  login(email: string, password: string, role: string): Observable<boolean> {
    const endpoint = role.toLowerCase() === 'supervisor'
      ? `${this.apiUrl}/supervisors/login`
      : `${this.apiUrl}/workers/login`;

    return this.http.post<any>(endpoint, { email, password }).pipe(
      tap(response => this.handleAuthResponse(response)),
      map(() => true),
      catchError(err => {
        console.error('Login rejected by backend:', err);
        return of(false);
      })
    );
  }

  register(user: any): Observable<boolean> {
    const backendPayload = {
      name: user.name,
      email: user.email,
      password: user.password,
      address: user.address,
      phone_no: user.phoneNumber,
      company_id: Number(user.companyId),
      role: user.role.toLowerCase(),
      wage: 15.00
    };

    const endpoint = user.role.toLowerCase() === 'supervisor'
      ? `${this.apiUrl}/supervisors/register`
      : `${this.apiUrl}/workers/register`;

    return this.http.post<any>(endpoint, backendPayload).pipe(
      tap(response => this.handleAuthResponse(response)),
      map(() => true),
      catchError(err => {
        console.error('Registration failed:', err);
        return of(false);
      })
    );
  }

  getCompanies(): Observable<any[]> {
    return this.http.get<any[]>(`${this.apiUrl}/companies/`);
  }

  logout() {
    localStorage.removeItem('jwt_token');
    localStorage.removeItem('user_data');
    this.currentUserSubject.next(null);
  }
}
