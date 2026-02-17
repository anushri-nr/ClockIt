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
  template: `
    <div class="page-container flex-center">
  
      <mat-toolbar class="auth-toolbar">
        <span class="brand">ClockIt</span>
        <span class="spacer"></span>
      </mat-toolbar>

      <mat-card class="auth-card">
        <div class="brand-logo">
          <mat-icon class="logo-icon">schedule</mat-icon>
        </div>

        <mat-card-header class="centered-header">
          <mat-card-title>Create {{ user.role | titlecase }} Account</mat-card-title>
        </mat-card-header>
        
        <mat-card-content>
          
          <ng-container *ngIf="companies$ | async as companiesList; else loading">
            
            <form #registerForm="ngForm" (ngSubmit)="registerForm.valid && onSubmit()">
              
              <mat-form-field appearance="outline" class="full-width">
                <mat-label>Select Company</mat-label>
                <mat-select [(ngModel)]="user.companyId" name="companyId" required>
                  <mat-option *ngFor="let company of companiesList" [value]="company.id">
                    {{ company.name }}
                  </mat-option>
                </mat-select>
                <mat-error>Company is required</mat-error>
              </mat-form-field>

              <mat-form-field appearance="outline" class="full-width">
                <mat-label>Full Name</mat-label>
                <input matInput [(ngModel)]="user.name" name="name" required>
              </mat-form-field>

              <mat-form-field appearance="outline" class="full-width">
                <mat-label>Email</mat-label>
                <input matInput type="email" [(ngModel)]="user.email" name="email" required email>
              </mat-form-field>

              <mat-form-field appearance="outline" class="full-width">
                <mat-label>Password</mat-label>
                <input matInput type="password" [(ngModel)]="user.password" name="password" required minlength="6">
              </mat-form-field>

              <mat-form-field appearance="outline" class="full-width">
                <mat-label>Phone Number</mat-label>
                <input matInput type="tel" [(ngModel)]="user.phoneNumber" name="phoneNumber" required>
              </mat-form-field>

              <mat-form-field appearance="outline" class="full-width">
                <mat-label>Address</mat-label>
                <textarea matInput [(ngModel)]="user.address" name="address" required rows="2"></textarea>
              </mat-form-field>

              <div class="actions">
                <button mat-flat-button color="primary" type="submit" 
                        [disabled]="!registerForm.valid" class="full-width large-btn">
                  Create Account
                </button>
                
                <button mat-button color="warn" routerLink="/login" [queryParams]="{ role: user.role }">
                  Cancel
                </button>
              </div>

            </form>
          </ng-container>

          <ng-template #loading>
            <div class="spinner-container">
              <mat-spinner diameter="40"></mat-spinner>
              <p>Loading companies...</p>
            </div>
          </ng-template>

        </mat-card-content>
      </mat-card>
    </div>
  `,
  styles: [`
    .auth-card { width: 100%; max-width: 450px; padding: 30px; margin-top: 100px; }
    .brand-logo { display: flex; justify-content: center; margin-bottom: 20px; }
    .logo-icon { font-size: 48px; height: 48px; width: 48px; color: #0f5f5c; background: #d7efe9; padding: 15px; border-radius: 50%; }
    .centered-header { justify-content: center; margin-bottom: 20px; }
    .full-width { width: 100%; margin-bottom: 8px; }
    .large-btn { padding: 25px 0; font-size: 1.1rem; }
    .actions { display: flex; flex-direction: column; gap: 10px; margin-top: 15px; }
    .auth-toolbar { position: absolute; top: 0; left: 0; width: 100%; background-color: #0f5f5c; color: white; box-shadow: 0 4px 12px rgba(15, 95, 92, 0.2); padding: 0 24px; box-sizing: border-box; }
    .brand { font-weight: 700; letter-spacing: 1px; font-size: 1.2rem; }
    .spacer { flex: 1 1 auto; }
    .spinner-container { display: flex; flex-direction: column; align-items: center; justify-content: center; padding: 40px; color: #666; }
  `]
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
    
    this.authService.register(this.user).subscribe((success) => {
      this.isSubmitting = false;
      if(success) {
        if(this.user.role.toLowerCase() == "worker") {
          this.router.navigate(['worker-dashboard']);
        } else {
          this.router.navigate(['supervisor-dashboard']);
        }
      } else {
        alert("Registration failed");
      }
    });
  }
}