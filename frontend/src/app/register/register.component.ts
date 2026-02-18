import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms'; 
import { Router, ActivatedRoute, RouterModule } from '@angular/router';
import { Observable } from 'rxjs';
import { MatCardModule } from '@angular/material/card';
import { MatInputModule } from '@angular/material/input';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatButtonModule } from '@angular/material/button';
import { AuthService } from '../services/auth.service';
import { MatIconModule } from '@angular/material/icon';
import { MatToolbarModule } from '@angular/material/toolbar';
import { MatSelectModule } from '@angular/material/select';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { CompanyService, Company } from '../services/company.service';

@Component({
  selector: 'app-register',
  standalone: true,
  imports: [
    CommonModule, 
    FormsModule, 
    RouterModule,
    MatCardModule, 
    MatInputModule, 
    MatFormFieldModule, 
    MatButtonModule,
    MatIconModule,
    MatToolbarModule,
    MatSelectModule,
    MatProgressSpinnerModule
  ],
  templateUrl: './register.component.html',
  styleUrls: ['./register.component.scss']
})
export class RegisterComponent implements OnInit {
  user = {
    name: '',
    email: '',
    password: '',    
    phoneNumber: '', 
    address: '',     
    role: 'Worker',
    companyId: ''
  };
  
  companies$!: Observable<Company[]>; 
  
  isSubmitting = false;

  constructor(
    private router: Router,
    private route: ActivatedRoute,
    private authService: AuthService,
    private companyService: CompanyService
  ) {}

  ngOnInit() {
    const roleFromUrl = this.route.snapshot.queryParams['role'];
    if (roleFromUrl) {
      this.user.role = roleFromUrl;
    }

    this.companies$ = this.companyService.getCompanies();
  }

  onSubmit() {
    if(this.isSubmitting) return;
    this.isSubmitting = true;
    
    this.authService.register(this.user).subscribe({
      next: (success) => {
        this.isSubmitting = false;
        if(success) {
          if(this.user.role.toLowerCase() == "worker") {
            this.router.navigate(['worker-dashboard']);
          } else {
            this.router.navigate(['supervisor-dashboard']);
          }
        } else {
          alert("Registration failed. Please try again.");
        }
      },
      error: (err) => {
        console.error(err);
        this.isSubmitting = false;
        alert("An error occurred during registration.");
      }
    });
  }
}