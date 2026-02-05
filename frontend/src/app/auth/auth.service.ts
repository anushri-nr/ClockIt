import { Injectable } from '@angular/core';
import { Observable, of, throwError } from 'rxjs';
import { delay } from 'rxjs/operators';
import { LoginPayload } from './auth.models';

@Injectable({ providedIn: 'root' })
export class AuthService {
  // TODO: Replace mock logic with real API calls when endpoints are ready.
  login(payload: LoginPayload): Observable<{ token: string }> {
    if (payload.employeeId.toLowerCase().includes('fail')) {
      return throwError(() => new Error('Invalid credentials.'));
    }

    return of({ token: 'mock-token' }).pipe(delay(600));
  }
}
