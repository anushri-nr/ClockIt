import { inject } from '@angular/core';
import { CanActivateFn, Router } from '@angular/router';

export const authGuard: CanActivateFn = (route, state) => {
  const router = inject(Router);
  
  // Check if the user has a token
  const token = localStorage.getItem('jwt_token');

  if (token) {
    return true; // Let them in!
  }

  // No token? Kick them to the login page
  router.navigate(['/']);
  return false;
};