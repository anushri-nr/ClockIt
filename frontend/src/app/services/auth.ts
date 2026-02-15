import { Injectable } from '@angular/core';
import { of, Observable, delay } from 'rxjs';  
import { HttpClient } from '@angular/common/http';

@Injectable({
  providedIn: 'root'
})
export class AuthService {

  private apiUrl = 'http://localhost:8080/api';

  constructor(private http: HttpClient) { }

  
  login(email: string, password: string, role: string): Observable<boolean> {
    
    console.log(`Attempting login for ${email}, ${password} as ${role}`);
    
    return of(true).pipe(delay(1000)); 
  }

  register(user: any): Observable<any> {

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

    console.log('Sending to Backend:', backendPayload);

    const endpoint = user.role.toLowerCase() === 'supervisor'
      ? `${this.apiUrl}/supervisors/register`
      : `${this.apiUrl}/workers/register`;

    // return this.http.post(endpoint, backendPayload);

    return of(true).pipe(delay(1000));
  }

  getCompanies(): Observable<any[]> {
    const dummyCompanies = [
      { id: 1, name: 'Google' },
      { id: 2, name: 'Microsoft' },
      { id: 3, name: 'UF' },
      { id: 4, name: 'Chick-fil-A' },
    ];
    console.log('SERVICE: getCompanies called. Returning data in 1s...');
    return of(dummyCompanies).pipe(delay(1000));
  }
}