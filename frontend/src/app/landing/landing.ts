import { Component } from '@angular/core';
import { MatCardModule } from '@angular/material/card';
import { MatButtonModule } from '@angular/material/button';
import { Router, RouterModule } from '@angular/router';
import { MatIconModule } from '@angular/material/icon';
import { MatRipple } from '@angular/material/core';
import {MatButtonToggleModule} from '@angular/material/button-toggle';

@Component({
  selector: 'app-landing',
  standalone: true,
  imports: [MatCardModule, MatButtonModule, RouterModule, MatIconModule, MatRipple, MatButtonToggleModule],
  template: `
  <div class="landing-container">
    
    <div class="header-section">
      <h1> 
      Welcome to ClockIt
         </h1>
      <p>
        What are you...
         </p>
    </div>

    <div class="role-selection">
      <mat-button-toggle-group name="fontStyle" aria-label="Font Style">
        <mat-button-toggle (click)="onRoleSelect('supervisor')" value="Supervisor">Supervisor</mat-button-toggle>
        <mat-button-toggle (click)="onRoleSelect('worker')" value="Worker">Worker</mat-button-toggle>
      </mat-button-toggle-group>
    </div>

  </div>
`,
// ... keep your imports and logic
  styles: [/* ... styles ... */]
})
export class LandingComponent {
  constructor(private router: Router) { }

  onRoleSelect(roleName: string) {
    this.router.navigate(['/login'], { queryParams: { role: roleName } })
  }
}