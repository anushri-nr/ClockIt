import { HttpInterceptorFn } from '@angular/common/http';

export const authInterceptor: HttpInterceptorFn = (req, next) => {
  // Grab the token from local storage
  const token = localStorage.getItem('jwt_token');

  if (token) {
    // Clone the request and staple the Authorization header to it
    const clonedReq = req.clone({
      setHeaders: {
        Authorization: `Bearer ${token}`
      }
    });
    return next(clonedReq);
  }

  // If no token, just send the request as-is (it will likely fail with a 401)
  return next(req);
};