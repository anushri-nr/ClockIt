import { Routes } from '@angular/router';
import { LoginComponent } from './auth/login.component';
import { SupervisorDashboardComponent } from './supervisor-dashboard/supervisor-dashboard.component';

export const routes: Routes = [
  { path: '', component: SupervisorDashboardComponent },
  { path: 'login', component: LoginComponent },
  { path: 'supervisor', component: SupervisorDashboardComponent }
];
