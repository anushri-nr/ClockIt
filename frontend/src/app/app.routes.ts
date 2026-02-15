import { Routes } from '@angular/router';
import { LandingComponent } from './landing/landing';
import { WorkerDashboardComponent } from './worker-dashboard/worker-dashboard';
import { LoginComponent } from './login/login';
import { RegisterComponent } from './register/register';
import { SupervisorDashboardComponent } from './supervisor-dashboard/supervisor-dashboard.component';

export const routes: Routes = [
  { path: '', redirectTo: '/home', pathMatch: 'full' },
  { path: 'home', component: LandingComponent },
  { path: 'login', component: LoginComponent },
  { path: 'register', component: RegisterComponent },
  { path: 'worker-dashboard', component: WorkerDashboardComponent },
  { path: 'supervisor', component: SupervisorDashboardComponent }
];
