import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms'; 
import { Router, ActivatedRoute, RouterModule } from '@angular/router';

import { MatCardModule } from '@angular/material/card';
import { MatInputModule } from '@angular/material/input';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatButtonModule } from '@angular/material/button';

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
    MatButtonModule
  ], 
  template: `
    <div class="register-container">
      <mat-card class="register-card">
        <mat-card-header>
          <mat-card-title>Register as {{ user.role | titlecase }}</mat-card-title>
        </mat-card-header>
        
        <mat-card-content>
          <form #registerForm="ngForm" (ngSubmit)="registerForm.valid && onSubmit()">
            
            <mat-form-field appearance="outline" class="full-width">
              <mat-label>Full Name</mat-label>
              <input matInput [(ngModel)]="user.name" name="name" required>
              <mat-error>Name is required</mat-error>
            </mat-form-field>

            <mat-form-field appearance="outline" class="full-width">
              <mat-label>Email</mat-label>
              <input matInput type="email" [(ngModel)]="user.email" name="email" required email>
              <mat-error>Please enter a valid email</mat-error>
            </mat-form-field>

            <mat-form-field appearance="outline" class="full-width">
              <mat-label>Password</mat-label>
              <input matInput type="password" [(ngModel)]="user.password" name="password" required minlength="6">
              <mat-error>Password must be at least 6 characters</mat-error>
            </mat-form-field>

            <mat-form-field appearance="outline" class="full-width">
              <mat-label>Phone Number</mat-label>
              <input matInput type="tel" [(ngModel)]="user.phoneNumber" name="phoneNumber" required pattern="[0-9]*">
              <mat-error>Valid phone number is required</mat-error>
            </mat-form-field>

            <mat-form-field appearance="outline" class="full-width">
              <mat-label>Address</mat-label>
              <textarea matInput [(ngModel)]="user.address" name="address" required rows="3"></textarea>
              <mat-error>Address is required</mat-error>
            </mat-form-field>

            <div class="actions">
              <button mat-raised-button color="primary" type="submit" 
                      [disabled]="!registerForm.valid">
                Create Account
              </button>
              
              <button mat-button color="warn" routerLink="/login" [queryParams]="{ role: user.role }">
                Cancel
              </button>
            </div>

          </form>
        </mat-card-content>
      </mat-card>
    </div>
  `,
  styles: [`
    .register-container { 
      display: flex; 
      justify-content: center; 
      align-items: center; 
      min-height: 100vh; 
      background-color: #f5f5f5; 
      padding: 20px; 
    }
    .register-card { 
      width: 100%; 
      max-width: 450px; 
    }
    .full-width { 
      width: 100%; 
      margin-bottom: 10px; 
    }
    .actions {
      display: flex;
      flex-direction: column;
      gap: 10px;
      margin-top: 10px;
    }
    mat-card-title {
      margin-bottom: 20px;
      display: block;
      text-align: center;
    }
  `]
})
export class RegisterComponent implements OnInit {
  
  user = {
    name: '',
    email: '',
    password: '',    
    phoneNumber: '', 
    address: '',     
    role: 'Worker' 
  };

  constructor(private router: Router, private route: ActivatedRoute) {}

  ngOnInit() {
    const roleFromUrl = this.route.snapshot.queryParams['role'];
    if (roleFromUrl) {
      this.user.role = roleFromUrl;
    }
  }

  onSubmit() {
    console.log('Sending to Backend:', this.user);
    alert('Registration Successful!');
    
    if (this.user.role.toLowerCase() === 'worker') {
        this.router.navigate(['/worker-dashboard']);
    } else {
        alert("Supervisor Dashboard not built yet.");
    }
  }
}