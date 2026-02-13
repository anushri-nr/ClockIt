import { Component, OnInit } from '@angular/core';
import { FormsModule } from '@angular/forms'; 
import { Router, ActivatedRoute } from '@angular/router';
import { CommonModule } from '@angular/common';
import { JsonPipe } from '@angular/common';

@Component({
  selector: 'app-register',
  standalone: true,
  imports: [FormsModule, CommonModule, JsonPipe], 
  template: `
    <div style="padding: 20px; font-family: sans-serif; max-width: 400px; margin: 0 auto;">
      <h2>Register as {{ user.role }}</h2>
      
      <form #registerForm="ngForm" (ngSubmit)="registerForm.valid && onSubmit()" style="display: flex; flex-direction: column; gap: 10px;">
        
        <div>
          <label for="fullName">Full Name:</label>
          <input type="text" id="fullName" [(ngModel)]="user.name" name="name" required style="width: 100%;">
          <div *ngIf="registerForm.submitted && !user.name" style="color: red; font-size: 12px;">Name is required</div>
        </div>

        <div>
          <label for="email">Email:</label>
          <input type="email" id="email" [(ngModel)]="user.email" name="email" required email style="width: 100%;">
        </div>

        <div>
          <label for="password">Password:</label>
          <input type="password" id="password" [(ngModel)]="user.password" name="password" required style="width: 100%;">
        </div>

        <div>
          <label for="phone">Phone:</label>
          <input type="tel" id="phone" [(ngModel)]="user.phoneNumber" name="phone" required style="width: 100%;">
        </div>

        <div>
          <label for="address">Address:</label>
          <textarea id="address" [(ngModel)]="user.address" name="address" required style="width: 100%; height: 60px;"></textarea>
        </div>

        <button type="submit" 
                [disabled]="!registerForm.valid"
                style="padding: 10px; cursor: pointer; margin-top: 10px;"
                [style.background-color]="registerForm.valid ? '#007bff' : 'grey'"
                [style.color]="'white'">
          Create Account
        </button>

      </form>
      
      <hr>
      <p style="font-size: 12px; color: gray;">Debug Payload: {{ user | json }}</p>
    </div>
  `
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
    

    alert('Registration Successful for: ' + this.user.name);
    
    if (this.user.role.toLowerCase() === 'worker') {
        this.router.navigate(['/worker-dashboard']);
    } else {
        alert("Supervisor Dashboard not built yet.");
    }
  }
}