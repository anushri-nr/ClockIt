import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { ActivatedRoute, Router, RouterModule } from '@angular/router';
import { MatCardModule } from '@angular/material/card';
import { MatInputModule } from '@angular/material/input';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatButtonModule } from '@angular/material/button';
import { AuthService } from '../services/auth.service';
import { MatIconModule } from '@angular/material/icon';
import { MatToolbarModule } from '@angular/material/toolbar';

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
    MatButtonModule,
    MatIconModule,
    MatToolbarModule
  ],
  templateUrl: './login.component.html',
  styleUrls: ['./login.component.scss']
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
    } else {
      this.router.navigate(['/']);
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
        if(this.currentRole.toLowerCase() == "worker") {
          this.router.navigate(['worker-dashboard']);
        } else {
          this.router.navigate(['supervisor-dashboard']);
        }
      } else {
        alert("Login failed!(Check the console)")
      }
    });
  }
}