import { Component } from '@angular/core';
import { CommonModule } from '@angular/common';
import { Router, RouterModule } from '@angular/router';
import { MatCardModule } from '@angular/material/card';
import { MatButtonModule } from '@angular/material/button';
import { MatButtonToggleModule } from '@angular/material/button-toggle';
import { MatIconModule } from '@angular/material/icon';
import { MatToolbar, MatToolbarModule } from '@angular/material/toolbar';

@Component({
  selector: 'app-landing',
  standalone: true,
  imports: [
    CommonModule,
    RouterModule,
    MatCardModule,
    MatButtonModule,
    MatButtonToggleModule,
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
          <mat-card-title>Welcome to ClockIt</mat-card-title>
          <mat-card-subtitle>Select your portal to continue</mat-card-subtitle>
        </mat-card-header>
        
        <mat-card-content class="content-center">
          <div class="role-selection">
            <mat-button-toggle-group #group="matButtonToggleGroup" vertical="false">
              <mat-button-toggle (click)="onRoleSelect('supervisor')" value="Supervisor">
                <mat-icon>admin_panel_settings</mat-icon> Supervisor
              </mat-button-toggle>
              <mat-button-toggle (click)="onRoleSelect('worker')" value="Worker">
                <mat-icon>engineering</mat-icon> Worker
              </mat-button-toggle>
            </mat-button-toggle-group>
          </div>
        </mat-card-content>
      </mat-card>
    </div>
  `,
  styles: [`

    .auth-card {
      width: 100%;
      max-width: 400px;
      padding: 30px;
      text-align: center;
    }

    .centered-header {
      display: flex;
      flex-direction: column;
      align-items: center;
      margin-bottom: 20px;
    }

    mat-card-title {
      font-size: 2rem;
      font-weight: 700; 
      margin-bottom: 10px;
    }

    mat-icon {
      margin-right: 8px;
    }
    
    .role-selection {
      margin-top: 20px;
      display: flex;
      justify-content: center;
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
export class LandingComponent {
  constructor(private router: Router) { }

  onRoleSelect(roleName: string) {
    this.router.navigate(['/login'], { queryParams: { role: roleName } });
  }
}