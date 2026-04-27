import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { environment } from '../../environments/environment';

// Interface matching the Go Model exactly
export interface Company {
  id: number;       // matches json:"id"
  name: string;     // matches json:"name"
  address: string;  // matches json:"address"
}

@Injectable({
  providedIn: 'root'
})
export class CompanyService {

  // The endpoint defined in routes.go: companyRoutes.GET("/", controllers.ListCompanies)
  private apiUrl = `${environment.apiUrl}/companies/`; 

  constructor(private http: HttpClient) { }

  getCompanies(): Observable<Company[]> {
    return this.http.get<Company[]>(this.apiUrl);
  }
}