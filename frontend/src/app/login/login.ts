import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { ActivatedRoute, Router, RouterModule } from '@angular/router';
import { MatCardModule } from '@angular/material/card';
import { MatInputModule } from '@angular/material/input';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatButtonModule } from '@angular/material/button';
import { AuthService } from '../services/auth';

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
  .landing-container { 
    display: flex; 
    justify-content: center; 
    align-items: center; 
    height: 100vh; 
    /* No background color here, we let the global gradient shine through */
  }

  .landing-card { 
    width: 400px; 
    padding: 40px; /* More breathing room */
    text-align: center; 
    background: rgba(255, 255, 255, 0.95); /* Slightly transparent white */
    backdrop-filter: blur(10px); /* The "Frosted Glass" effect (Modern!) */
  }

  mat-card-title { 
    font-size: 2.2rem; 
    font-weight: 300; /* Thinner, more elegant font */
    color: #333;
    margin-bottom: 5px; 
  }

  mat-card-subtitle {
    font-size: 1rem;
    color: #666;
    margin-bottom: 30px; /* Push content down */
  }

  /* Make the icons in the buttons pop */
  mat-icon { 
    vertical-align: middle; 
    margin-right: 5px;
  }
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
        alert("Login Succesful");
        if(this.currentRole.toLowerCase() == "worker") {
          this.router.navigate(['worker-dashboard']);
        }
        else {
          alert("Supervisor Dashboard not built yet");
        }
      }
      else {
        alert("Login failed!(Check the console)")
      }
    });
  }
}