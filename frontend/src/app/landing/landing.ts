import { Component } from '@angular/core';
import { MatCardModule } from '@angular/material/card';
import { MatButtonModule } from '@angular/material/button';
import { Router, RouterModule } from '@angular/router';
import { MatIconModule } from '@angular/material/icon';
import { MatRipple } from '@angular/material/core';

@Component({
  selector: 'app-landing',
  standalone: true,
  imports: [MatCardModule, MatButtonModule, RouterModule, MatIconModule, MatRipple],
  template: `
    <div class="landing-container">
      <div class="role-selection">
        
        <mat-card class="role-card supervisor" matRipple (click)="onRoleSelect('supervisor')"> 
          <mat-card-header>
            <mat-card-title>Supervisor</mat-card-title>
          </mat-card-header>
          </mat-card>

        <mat-card class="role-card worker" matRipple (click)="onRoleSelect('worker')">
          <mat-card-header>
            <mat-card-title>Worker</mat-card-title>
          </mat-card-header>
          </mat-card>

      </div>
    </div>
  `,
  styles: [/* ... styles ... */]
})
export class LandingComponent {
  constructor(private router: Router) {}

  onRoleSelect(roleName: string) {
    this.router.navigate(['/login'], {queryParams: { role:roleName }})
  }
}