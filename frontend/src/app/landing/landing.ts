import { Component } from '@angular/core';
import { CommonModule } from '@angular/common'; // Added for basic directives
import { Router, RouterModule } from '@angular/router';

// MATERIAL IMPORTS
import { MatCardModule } from '@angular/material/card';
import { MatButtonModule } from '@angular/material/button';
import { MatButtonToggleModule } from '@angular/material/button-toggle';
import { MatIconModule } from '@angular/material/icon';

@Component({
  selector: 'app-landing',
  standalone: true,
  imports: [
    CommonModule,
    RouterModule,
    MatCardModule, 
    MatButtonModule, 
    MatButtonToggleModule, 
    MatIconModule
  ],
  template: `
    <div class="landing-container">
      <mat-card class="landing-card">
        <mat-card-header>
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
export class LandingComponent {
  constructor(private router: Router) { }

  onRoleSelect(roleName: string) {
    this.router.navigate(['/login'], { queryParams: { role: roleName } });
  }
}