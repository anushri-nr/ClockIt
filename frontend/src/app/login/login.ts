import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { ActivatedRoute, Router, RouterModule } from '@angular/router';

// MATERIAL IMPORTS
import { MatCardModule } from '@angular/material/card';
import { MatInputModule } from '@angular/material/input';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatButtonModule } from '@angular/material/button';

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
    MatButtonModule
  ],
  template: `
    <div class="login-container">
      <mat-card class="login-card">
        <mat-card-header>
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
              <button mat-raised-button color="primary" type="submit" 
                      [disabled]="!loginForm.valid" class="full-width">
                Sign In
              </button>
            </div>
          </form>
        </mat-card-content>

        <mat-card-actions class="center-actions">
           <span style="font-size: 12px; color: gray;">New here?</span>
           <button mat-button color="accent" (click)="onRegister()">
             Create {{ currentRole | titlecase }} Account
           </button>
        </mat-card-actions>
      </mat-card>
    </div>
  `,
  styles: [`
    .login-container { display: flex; justify-content: center; align-items: center; min-height: 100vh; background-color: #f5f5f5; padding: 20px; }
    .login-card { width: 100%; max-width: 400px; padding: 20px; }
    .full-width { width: 100%; margin-bottom: 5px; }
    .actions { margin-top: 10px; }
    
    /* NEW CSS RULE FOR CENTERING */
    .center-actions {
      display: flex;
      flex-direction: column;
      align-items: center;
      justify-content: center;
      gap: 5px;
      padding-top: 10px;
    }
  `]
})
export class LoginComponent implements OnInit {
  email = '';
  password = '';
  currentRole = 'Worker';

  constructor(private route: ActivatedRoute, private router: Router) {}

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
    console.log('Login Payload:', { email: this.email, role: this.currentRole });
    
    if (this.currentRole.toLowerCase() === 'worker') {
      this.router.navigate(['/worker-dashboard']);
    } else {
      alert("Supervisor Dashboard not built yet.");
    }
  }
}