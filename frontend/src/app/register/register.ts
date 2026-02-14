import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms'; 
import { Router, ActivatedRoute, RouterModule } from '@angular/router';
import { MatCardModule } from '@angular/material/card';
import { MatInputModule } from '@angular/material/input';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatButtonModule } from '@angular/material/button';
import { AuthService } from '../services/auth';

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
    /* The global gradient from styles.scss provides the background */
    padding: 20px; 
  }

  .register-card { 
    width: 100%; 
    max-width: 500px; /* Slightly wider than login for better spacing */
    padding: 30px;
    
    /* THE GLASS EFFECT */
    background: rgba(255, 255, 255, 0.95);
    backdrop-filter: blur(10px);
    border: 1px solid rgba(255, 255, 255, 0.2);
    box-shadow: 0 8px 32px 0 rgba(31, 38, 135, 0.37);
  }

  .full-width { 
    width: 100%; 
    margin-bottom: 5px; 
  }

  .actions {
    display: flex;
    flex-direction: column;
    gap: 15px; /* More breathing room between buttons */
    margin-top: 20px;
  }

  mat-card-title {
    font-size: 1.8rem;
    text-align: center;
    color: #333;
    margin-bottom: 20px;
    font-weight: 300;
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

  isLoading = false;

  constructor(
    private router: Router,
    private route: ActivatedRoute,
    private authService: AuthService,
  ) {}

  ngOnInit() {
    const roleFromUrl = this.route.snapshot.queryParams['role'];
    if (roleFromUrl) {
      this.user.role = roleFromUrl;
    }
  }

  onSubmit() {
    if(this.isLoading) return;
    this.isLoading = true;
    
    this.authService.register(this.user)
    .subscribe((success) => {
      if(success) {
        alert("Registration Succesful");
        if(this.user.role.toLowerCase() == "worker") {
          this.router.navigate(['worker-dashboard']);
        }
        else {
          alert("Supervisor Dashboard not built yet");
        }
      }
      else {
        alert("Registration failed (Check the console)");
      }
    });
  }
}