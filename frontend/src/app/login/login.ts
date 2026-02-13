import { Component, OnInit } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { ActivatedRoute, Router } from '@angular/router';
import { CommonModule } from '@angular/common';

@Component({
    selector: 'app-login',
    imports: [FormsModule, CommonModule],
    template: `
  <div style="padding: 20px; font-family: sans-serif; max-width: 400px; margin: 0 auto;">
  <h2>Login</h2>
  <form #registerForm="ngForm" (ngSubmit)="registerForm.valid && onSubmit()">
  <div>
    <label for="email">Email:</label>
          <input type="email" id="email" [(ngModel)]="email" name="email" required email style="width: 100%;">
  </div>

  <div>
    <label for="password">Password:</label>
          <input type="password" id="password" [(ngModel)]="password" name="password" required style="width: 100%;">
    </div>
    <button type="submit" 
                [disabled]="!registerForm.valid"
                style="padding: 10px; cursor: pointer; margin-top: 10px;"
                [style.background-color]="registerForm.valid ? '#007bff' : 'grey'"
                [style.color]="'white'">
          Login
        </button>
  </form>
  <br><br>
  <div (click)=onRegister() class="register-button">
    <button matButton="outlined">Register</button>
  </div>

</div>
  `
})
export class LoginComponent {
    email: string = '';
    password: string = '';
    currentRole: string = '';

    constructor(private route: ActivatedRoute, private router: Router) {}

    ngOnInit() {
      this.currentRole = this.route.snapshot.queryParams['role'];
    }

    onRegister() {
      this.router.navigate(['/register'], { queryParams: { role: this.currentRole } })
    }


    onSubmit() {

        console.log('Username:', this.email);
        console.log('Password:', this.password);
        console.log('Role:', this.currentRole);

        if (this.currentRole.toLowerCase() === 'worker') {
        this.router.navigate(['/worker-dashboard']);
    } else {
        alert("Supervisor Dashboard not built yet.");
    }
    }
}
