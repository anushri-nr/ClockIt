import { Component, inject } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormBuilder, ReactiveFormsModule, Validators } from '@angular/forms';
import { AuthService } from './auth.service';

@Component({
  selector: 'app-login',
  standalone: true,
  imports: [CommonModule, ReactiveFormsModule],
  templateUrl: './login.component.html',
  styleUrl: './login.component.scss'
})
export class LoginComponent {
  isSubmitting = false;
  errorMessage = '';
  successMessage = '';

  private fb = inject(FormBuilder);
  private auth = inject(AuthService);

  form = this.fb.group({
    role: ['Employee' as 'Employee' | 'Manager'],
    employeeId: ['', [Validators.required]],
    password: ['', [Validators.required, Validators.minLength(6)]]
  });

  submit(): void {
    this.errorMessage = '';
    this.successMessage = '';

    if (this.form.invalid) {
      this.form.markAllAsTouched();
      return;
    }

    this.isSubmitting = true;
    this.auth.login({
      role: this.form.value.role ?? 'Employee',
      employeeId: this.form.value.employeeId ?? '',
      password: this.form.value.password ?? ''
    }).subscribe({
      next: () => {
        this.isSubmitting = false;
        this.successMessage = 'Login successful.';
      },
      error: (err: Error) => {
        this.isSubmitting = false;
        this.errorMessage = err.message || 'Login failed.';
      }
    });
  }
}
