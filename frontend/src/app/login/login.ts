import { Component, OnInit } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { ActivatedRoute, Router } from '@angular/router';

@Component({
    selector: 'app-login',
    imports: [FormsModule],
    template: `
  <div>
  <h2>Login</h2>
  <form (ngSubmit)="onSubmit()">
    <label for="username">Username:</label>
    <input type="text" id="username" [(ngModel)]="username" name="username" required>

    <label for="password">Password:</label>
    <input type="password" id="password" [(ngModel)]="password" name="password" required>

    <button type="submit">Login</button>
  </form>

  <div (click)=onRegister() class="register-button">
    <button matButton="outlined">Register</button>
  </div>

</div>
  `
})
export class LoginComponent {
    username: string = '';
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
        // Implement your login logic here
        console.log('Username:', this.username);
        console.log('Password:', this.password);
        console.log('Role:', this.currentRole);
        // Add authentication logic and navigate to the next page upon successful login
    }
}
