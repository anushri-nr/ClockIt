import { Injectable } from '@angular/core';
import { of, Observable } from 'rxjs'; 
import { delay } from 'rxjs/operators'; 

@Injectable({
  providedIn: 'root'
})
export class AuthService {

  constructor() { }

  
  login(email: string, password: string, role: string): Observable<boolean> {
    
    console.log(`Attempting login for ${email}, ${password} as ${role}`);
    
    return of(true).pipe(delay(1000)); 
  }

  register(userData: any): Observable<boolean> {
    console.log('Sending to DB:', userData);
    return of(true).pipe(delay(1000));
  }
}