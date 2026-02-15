import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { ActivatedRoute, Router, RouterModule } from '@angular/router';
import { MatCardModule } from '@angular/material/card';
import { MatInputModule } from '@angular/material/input';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatButtonModule } from '@angular/material/button';
import { AuthService } from '../services/auth';
import { MatIconModule } from '@angular/material/icon';
import { MatToolbarModule } from '@angular/material/toolbar';

@Component({
  selector: 'app-login',
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
    MatToolbarModule
  ],
  template: `
    <div class="page-container flex-center">
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
          <mat-card-title>{{ currentRole | titlecase }} Login</mat-card-title>
        </mat-card-header>
        
        <mat-card-content>
          <form #loginForm="ngForm" (ngSubmit)="loginForm.valid && onSubmit()">
            
            <mat-form-field appearance="outline" class="full-width">
              <mat-label>Email</mat-label>
              <input matInput type="email" [(ngModel)]="email" name="email" required email>
              <mat-error>Valid email is required</mat-error>
            </mat-form-field>

            <mat-form-field appearance="outline" class="full-width">
              <mat-label>Password</mat-label>
              <input matInput type="password" [(ngModel)]="password" name="password" required>
              <mat-error>Password is required</mat-error>
            </mat-form-field>

            <div class="actions">
              <button mat-flat-button color="primary" type="submit" 
                      [disabled]="!loginForm.valid" class="full-width large-btn">
                Sign In
              </button>
            </div>
          </form>
        </mat-card-content>

        <mat-card-actions class="center-actions">
           <span style="font-size: 14px; color: #666;">New here?</span>
           <button mat-button color="accent" (click)="onRegister()">
             Create Account
           </button>
        </mat-card-actions>
      </mat-card>
    </div>
  `,
  styles: [`
    .auth-card {
      width: 100%;
      max-width: 400px;
      padding: 30px;
    }

    .centered-header {
      justify-content: center;
      margin-bottom: 20px;
    }

    .full-width {
      width: 100%;
      margin-bottom: 5px;
    }

    .large-btn {
      padding: 25px 0; 
      font-size: 1.1rem;
    }

    .center-actions {
      display: flex;
      flex-direction: column;
      align-items: center;
      gap: 5px;
      margin-top: 10px;
    }
      .auth-toolbar {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  background-color: #0f5f5c; 
  color: white;
  box-shadow: 0 4px 12px rgba(15, 95, 92, 0.2);
  padding: 0 24px;
  box-sizing: border-box; 
}

.brand {
  font-weight: 700;
  letter-spacing: 1px;
  font-size: 1.2rem;
}

.spacer { flex: 1 1 auto; }
  `]
})
export class LoginComponent implements OnInit {
  email = '';
  password = '';
  currentRole = 'Worker';
  isLoading = false;

  constructor(
    private route: ActivatedRoute,
    private router: Router,
    private authService: AuthService,
    ) {}

  ngOnInit() {
    const roleFromUrl = this.route.snapshot.queryParams['role'];
    if (roleFromUrl) {
      this.currentRole = roleFromUrl;
    }
  }

  onRegister() {
    this.router.navigate(['/register'], { queryParams: { role: this.currentRole } });
  }

  onSubmit() {
    if(this.isLoading) return;
    this.isLoading = true;
    
    this.authService.login(this.email, this.password, this.currentRole)
    .subscribe((success) => {
      this.isLoading = false;
      if(success) {
        if(this.currentRole.toLowerCase() == "worker") {
          this.router.navigate(['worker-dashboard']);
        } else {
          this.router.navigate(['supervisor-dashboard']);
        }
      } else {
        alert("Login failed!(Check the console)")
      }
    });
  }
}