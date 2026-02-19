import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';

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
  private apiUrl = 'http://localhost:8080/api/companies/'; 

  constructor(private http: HttpClient) { }

  getCompanies(): Observable<Company[]> {
    return this.http.get<Company[]>(this.apiUrl);
  }
}